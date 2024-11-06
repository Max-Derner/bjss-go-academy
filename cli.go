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
	for err = nil ; err == nil ; {
		_, err = fmt.Fscanf(r, "%c", &char)
		if char == '\n' {break}
		input += string(char)
	}
	fmt.Println()
	return input
}

func CliPromptForToDoItem() ToDoItem {
	title := Title(
		Input("Enter a title for this to do item: "),
	)
	priority := Priority(
		Input("Enter a priority for this to do item: "),
	)
	return ConstructToDoItem(title, priority, false)
}


