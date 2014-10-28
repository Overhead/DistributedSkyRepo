package dht

import (
	"fmt"
	"time"
	"sync"
)

type DataTable struct { 
	Data map[string] string
	Timestamp time.Time
	lock sync.RWMutex
}

func CreateDataTable() *DataTable {
	table := DataTable{}
	table.Data = make(map[string]string)
	table.Timestamp = time.Now()
	return &table
}

func (table *DataTable) Get(value string) (string, bool) {
	key := sha1hash(value)
	_, val := table.Data[key]
	if val {
		return table.Data[key], true
	} else {
		fmt.Printf("Item doesn't exists\n")
		return "", false
	}
}

func (table *DataTable) Add(value string) bool {
	key := sha1hash(value)
	table.lock.Lock()
	var result bool = false
	defer table.lock.Unlock()
	_, val := table.Data[key]
	if !val {
		table.Data[key] = value
		table.Timestamp = time.Now()
		result = true
	} else {
		fmt.Printf("Item allready exists\n")
		result = false
	}
	table.lock.Unlock()
	return result
}

func (table *DataTable) Remove(value string) bool {
	key := sha1hash(value)
	table.lock.Lock()
	var result bool = false
	defer table.lock.Unlock()
	_, val := table.Data[key]
	if val {
		delete(table.Data, key)
		table.Timestamp = time.Now()
		result = true
	} else {
		fmt.Printf("Item doesn't exists\n")
		result = false
	}
	table.lock.Unlock()
	return result
}

func (table *DataTable) Update(value string) bool {
	key := sha1hash(value)
	table.lock.Lock()
	var result bool = false
	defer table.lock.Unlock()
	_, val := table.Data[key]
	if val {
		table.Data[key] = value
		table.Timestamp = time.Now()
		result = true
	} else {
		fmt.Printf("Item doesn't exists\n")
		result = false
	}
	table.lock.Unlock()
	return result	
}
