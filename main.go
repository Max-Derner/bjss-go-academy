package main

// import (
// 	"fmt"
// )

func main() {
	db := inMemoryDataStore{make(map[Id]ToDoItem)}
	dal := NewDataAccessLayer(&db)
	go StartAPI(dal)
	RunCli(dal)
}
