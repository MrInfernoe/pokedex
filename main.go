package main

import (
	langShell "pokedex/internal/langShell"
)

func main() {
	print("Hello, World!\n")
	stringSlice := langShell.CleanInput("hello world")
	for _, s := range stringSlice {
		print(s + "\n")
	}
}