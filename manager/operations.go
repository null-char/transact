package manager

import (
	"fmt"

	s "github.com/null-char/transact/store"
	"github.com/null-char/transact/utils"
)

// OperationsManager manages all operations performed on the store by the user and prints any required information
type OperationsManager struct {
	tm *TransactionManager
}

// MakeOperationsManager constructs a new operatiton manager from a transaction manager
func MakeOperationsManager(tm *TransactionManager) *OperationsManager {
	return &OperationsManager{tm}
}

// Get fetches the value associated with the key and prints and formats it appropriately
func (om OperationsManager) Get(key s.Key) {
	if val, ok := om.tm.Get(key); ok {
		utils.PrintValue(val)
	} else {
		fmt.Println("Not found")
	}
}

// Set sets the given key's associated value. Updates the value if the key already exists.
func (om OperationsManager) Set(key s.Key, val s.Value) {
	if updated := om.tm.Set(key, val); updated {
		fmt.Println("Updated")
	} else {
		fmt.Println("Added")
	}
}

// Delete removes the key from the store entirely. Prints out an error if the key doesn't exist in the first place.
func (om OperationsManager) Delete(key s.Key) {
	if err := om.tm.Delete(key); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Deleted")
	}
}

// Count counts all instances of the provided value and prints the output. Prints 0 if there are no associations with this value.
func (om OperationsManager) Count(val s.Value) {
	n := om.tm.Count(val)
	fmt.Println(n)
}
