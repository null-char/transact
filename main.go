package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	m "github.com/null-char/transact/manager"
	p "github.com/null-char/transact/parser"
	"github.com/null-char/transact/store"
	"github.com/null-char/transact/utils"
)

func main() {
	initial := utils.LoadData()

	// The global store
	s := store.MakeNewStoreWithData(initial)
	tm := m.MakeTransactionManager(s)
	om := m.MakeOperationsManager(tm)

	rd := os.Stdin
	parser := p.MakeParser(rd)

	go onShutdown(s)

	// Parse input from the user until sigint
	for {
		rd.WriteString("> ")
		op, args, err := parser.Run()
		if err != nil {
			continue
		}
		errorMsg := "ERROR: Insufficient arguments"

		switch op {
		case "SET":
			if len(args) < 2 {
				fmt.Println(errorMsg)
				continue
			}

			value := utils.ParseValue(args[1])
			om.Set(args[0], value)
			break

		case "GET":
			if len(args) < 1 {
				fmt.Println(errorMsg)
				continue
			}
			om.Get(args[0])
			break

		case "DELETE":
			if len(args) < 1 {
				fmt.Println(errorMsg)
				continue
			}
			om.Delete(args[0])
			break

		case "COUNT":
			if len(args) < 1 {
				fmt.Println(errorMsg)
				continue
			}

			val := utils.ParseValue(args[0])
			om.Count(val)
			break

		case "BEGIN":
			tm.PushTransaction()
			break

		case "COMMIT":
			tm.Commit()
			break

		case "ROLLBACK":
			tm.Rollback()
			break

		default:
			fmt.Println("ERROR: Unknown operation")
		}
	}
}

func onShutdown(s *store.Store) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		<-done
		fmt.Println("\n Saving data and exiting...")
		utils.SaveStore(*s)
		os.Exit(0)
	}()

	// Notify us (through the given channel) whenever we receive the following signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	done <- true
}
