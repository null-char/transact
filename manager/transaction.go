package manager

import (
	"errors"
	"fmt"

	s "github.com/null-char/transact/store"
)

// Transaction consists of a parent pointer which points to the transaction's parent context. It can be thought as a reversed linked list
// (nil if current Transaction is at the top) and a local store which resembles the data that needs to be updated
// into our global store if the transaction is committed.
type Transaction struct {
	localStore *s.Store
	parent     *Transaction
}

// TransactionManager manages the transaction stack
type TransactionManager struct {
	globalStore       *s.Store
	transactions      *Transaction
	activeTransaction *Transaction // represents top of the stack
	numTransactions   uint
}

// MakeTransactionManager constructs a new transaction manager supplied with a (global) store
func MakeTransactionManager(store *s.Store) *TransactionManager {
	return &TransactionManager{globalStore: store, numTransactions: 0}
}

// empty returns a bool indicating whether or not the transaction stack is empty (no ongoing transactions)
func (tm TransactionManager) empty() bool {
	return tm.numTransactions == 0
}

// getActiveStore returns the global store if there are no active transactions and the local store (of the active transaction) if there is one
func (tm TransactionManager) getActiveStore() *s.Store {
	if tm.empty() {
		return tm.globalStore
	}

	return tm.activeTransaction.localStore
}

// GetGlobalStore returns a copy of the global store for the purpose of reading
func (tm TransactionManager) GetGlobalStore() s.Store {
	return *tm.globalStore
}

// PushTransaction pushes a new transaction onto the stack inheriting the local store of its parent (if it has a parent)
func (tm *TransactionManager) PushTransaction() {
	ls := s.MakeNewStoreWithData(tm.getActiveStore().GetData())

	if tm.empty() {
		t := &Transaction{localStore: ls, parent: nil}
		tm.activeTransaction = t
		tm.transactions = tm.activeTransaction
	} else {
		t := &Transaction{localStore: ls}
		t.parent = tm.activeTransaction
		tm.activeTransaction = t
	}

	tm.numTransactions++
}

// popTransaction pops the active transaction (top of the stack) if there is one otherwise we have a stack underflow so we just log some error
func (tm *TransactionManager) popTransaction() {
	if tm.empty() {
		fmt.Println("ERROR: There are no transactions to pop. Stack is empty.")
	} else {
		tm.activeTransaction = tm.activeTransaction.parent
		tm.numTransactions--
	}
}

// Commit commits any changes made in the active transaction into the global store
func (tm *TransactionManager) Commit() {
	if tm.empty() {
		fmt.Println("ERROR: No transactions to commit")
	} else {
		localStore := tm.activeTransaction.localStore

		for _, pair := range localStore.GetKVPairs() {

			// For the child of a transaction, from their perspective, their parent represents the global store.
			// So, we update the data into the parent's local store instead. This way, the changes will cascade
			// up into the real global store.
			if tm.activeTransaction.parent != nil {
				tm.activeTransaction.parent.localStore.Set(pair.First, pair.Second)
			} else {
				// If the transaction has no parents (which means it's the only one in the stack), then we insert the data into
				// the ACTUAL global store.
				tm.globalStore.Set(pair.First, pair.Second)
			}
		}

		tm.popTransaction()
	}
}

// Rollback omits any changes made in the current transaction and pops it off the stack
func (tm *TransactionManager) Rollback() {
	if tm.empty() {
		fmt.Println("ERROR: No transactions to rollback")
	} else {
		// Make no changes and then simply pop off the transaction
		tm.popTransaction()
	}
}

// Set sets a new key value pair or updates it. An update is indicated by a return bool value of true.
func (tm *TransactionManager) Set(key s.Key, value s.Mappable) bool {
	s := tm.getActiveStore()
	// Checks to see whether we're updating an existing key's value or setting a new one
	updated := false
	if _, ok := s.Get(key); ok {
		updated = true
	}

	s.Set(key, value)
	return updated
}

// Get fetches the value associated with the provided key
func (tm *TransactionManager) Get(key s.Key) (s.Mappable, bool) {
	s := tm.getActiveStore()
	return s.Get(key)
}

// Delete removes the specified key value association
func (tm *TransactionManager) Delete(key s.Key) error {
	s := tm.getActiveStore()
	if ok := s.Delete(key); !ok {
		return errors.New("Not found")
	}

	return nil
}

// Count just counts the number of times the given value has been associated in the underlying store
func (tm *TransactionManager) Count(val s.Mappable) uint {
	s := tm.getActiveStore()
	return s.Count(val)
}
