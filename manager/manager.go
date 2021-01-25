package manager

import (
	"fmt"

	s "github.com/null-char/transact/store"
)

// Transaction consists of a next pointer which points to the next transaction in the stack
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

// Empty returns a bool indicating whether or not the transaction stack is empty (no ongoing transactions)
func (tm TransactionManager) Empty() bool {
	return tm.transactions == nil
}

// getActiveStore returns the global store if there are no active transactions and the local store (of the active transaction) if there is one
func (tm TransactionManager) getActiveStore() *s.Store {
	if tm.Empty() {
		return tm.globalStore
	}

	return tm.activeTransaction.localStore
}

// PushTransaction pushes a new transaction onto the stack inheriting the local store of its parent (if it has a parent)
func (tm *TransactionManager) PushTransaction() {
	if tm.Empty() {
		t := &Transaction{localStore: s.MakeNewStore(), parent: nil}
		tm.transactions = t
		tm.activeTransaction = t
	} else {
		store := s.MakeNewStore()

		// We duplicate the data of our parent here
		for _, pair := range tm.activeTransaction.localStore.GetKVPairs() {
			store.Set(pair.First, pair.Second)
		}

		t := &Transaction{localStore: store}
		t.parent = tm.activeTransaction
		tm.activeTransaction = t
	}

	tm.numTransactions++
}

// PopTransaction pops the active transaction (top of the stack) if there is one otherwise we have a stack underflow so we just log some error
func (tm *TransactionManager) PopTransaction() {
	if tm.Empty() {
		fmt.Println("ERROR: There are no transactions to pop. Stack is empty.")
	} else {
		tm.activeTransaction = tm.activeTransaction.parent
		tm.numTransactions--
	}

}

func (tm *TransactionManager) Set(key string, value s.Mappable) {
	s := tm.getActiveStore()
	s.Set(key, value)
}

func (tm *TransactionManager) Get(key string) (s.Mappable, bool) {
	s := tm.getActiveStore()
	return s.Get(key)
}
