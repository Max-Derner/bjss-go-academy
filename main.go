package main

// import (
// 	"fmt"
// )

func main() {
	// db := inMemoryDataStore{make(map[Id]ToDoItem)}
	db := newJSONDataStore()
	dal := NewDataAccessLayer(&db)
	go StartAPI(dal)
	RunCli(dal)
}
