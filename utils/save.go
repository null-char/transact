package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/null-char/transact/store"
)

func SaveStore(s store.Store) {
	if j, err := json.Marshal(s.GetData()); err != nil {
		fmt.Println("ERROR: Saving data failed")
	} else {
		// Truncates if the file already exists or creates it if it doesn't
		e := ioutil.WriteFile("data.json", j, 0644)
		if e != nil {
			panic(err)
		}
	}
}

func LoadData() map[store.Key]store.Value {
	fmt.Println("Attempting to load saved data from disk...")
	// The map to dump all the decoded JSON into
	dmp := make(map[store.Key]interface{})
	// The actual map that'll we use
	m := make(map[store.Key]store.Value)

	// Dump all file contents into memory so that we can unmarshal it
	dat, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Printf("%s \n", err.Error())
		return m
	}

	if err := json.Unmarshal(dat, &dmp); err != nil {
		fmt.Printf("%v \n", dmp)
		fmt.Println("ERROR: Unable to decode saved data from JSON. Defaulting to empty store.")
		return m
	}

	for k, v := range dmp {
		m[k] = ParseValue(fmt.Sprintf("%v", v))
	}

	fmt.Println("Loaded saved data from disk")
	return m
}
