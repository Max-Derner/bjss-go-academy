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

func FInput(w io.Reader) string {
	var input string
	fmt.Fscanln(w, &input)
	return input
}
