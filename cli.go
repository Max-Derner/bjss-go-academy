package main

import (
	"fmt"
	"io"
	"os"
)

func Input(prompt string) string {
	fmt.Print(prompt)
	return FInput(os.Stdin)
}

func FInput(r io.Reader) string {
	var input string
	var char rune
	var err error
	for err = nil; err == nil; {
		_, err = fmt.Fscanf(r, "%c", &char)
		if char == '\n' {
			break
		}
		input += string(char)
	}
	return input
}

func getTitle() Title {
	return Title(
		Input("Enter a title for this to do item: "),
	)
}

func getPriority() Priority {
	return Priority(
		Input("Enter a priority for this to do item: "),
	)
}

func cliPromptForToDoItem() ToDoItem {
	title := getTitle()
	priority := getPriority()
	return ConstructToDoItem(title, priority, false)
}

func printToDoItem(item ToDoItem) {
	var status string
	if item.Complete {
		status = "complete"
	} else {
		status = "incomplete"
	}
	fmt.Printf("| %s | %s | %s |\n", item.Title, item.Priority, status)
}

func choseFromList(commandList []string) int {
	if len(commandList) == 0 {
		return -1
	}
	for {
		fmt.Println("Choose a command!")
		for i, command := range commandList {
			fmt.Printf("| %d: %s ", i, command)
		}
		fmt.Print("|\n-> ")
		var selection int
		fmt.Scanf("%d", &selection)
		if selection < 0 || selection >= len(commandList) {
			fmt.Printf("Choose a number between 0 and %d", len(commandList)-1)
		} else {
			return selection
		}
	}
}

func cliRead(db *dataStore) {
	items := db.Read()
			if len(items) == 0 {
				fmt.Println("No items to show")
			} else {
				for _, item := range items {
					printToDoItem(item)
				}
			}
}

func cliAdd(db *dataStore) {
	err := db.Add(
		cliPromptForToDoItem(),
	)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func cliDelete(db *dataStore) {
	items := db.Read()
			if len(items) == 0 {
				fmt.Println("No items in database")
			} else {
				for i, item := range items {
					fmt.Printf("%d : ", i)
					printToDoItem(item)
				}
				fmt.Print("Choose item to delete: ")
				var choice int
				fmt.Scanf("%d", &choice)
				itemToDelete := items[choice]
				db.Delete(itemToDelete)
			}
}

func cliUpdate(db *dataStore) {
	items := db.Read()
	if len(items) == 0 {
		fmt.Println("No items in database")
	} else {
		for i, item := range items {
			fmt.Printf("%d : ", i)
			printToDoItem(item)
		}
		fmt.Print("Choose item to update: ")
		var choice int
		fmt.Scanf("%d", &choice)
		itemToUpdate := items[choice]
		actions := []string{
			"update title",
			"update priority",
			"mark complete",
			"mark incomplete",
		}
		selection := choseFromList(actions)
		switch actions[selection] {
		case "update title":
			itemToUpdate.Title = getTitle()
		case "update priority":
			itemToUpdate.Priority = getPriority()
		case "mark complete":
			if itemToUpdate.Complete {
				fmt.Println("To Do item is already complete!")
			} else {
				itemToUpdate.Complete = true
			}
		case "mark incomplete":
			if !itemToUpdate.Complete {
				fmt.Println("To Do item is already incomplete!")
			} else {
				itemToUpdate.Complete = false
			}
		}
		db.update(itemToUpdate)
	}
}

func RunCli() {
	fmt.Println("It's a todo app!")
	commandList := []string{
		"exit",
		"read",
		"add",
		"update",
		"delete",
	}
	db := NewEmptyDataStore()
	for {
		fmt.Println("\n=================================================")
		selection := choseFromList(commandList)
		switch commandList[selection] {
		case "exit":
			os.Exit(0)
		case "read":
			cliRead(&db)
		case "add":
			cliAdd(&db)
		case "delete":
			cliDelete(&db)
		case "update":
			cliUpdate(&db)
		}
	}
}
