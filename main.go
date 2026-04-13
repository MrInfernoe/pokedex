package main

import (
	"bufio"
	"os"
	"fmt"
	"pokedex/internal/repl"
)

func main() {
	var config repl.Config
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Println(scanner.Err())
			break
		}
		userText := scanner.Text()
		cleanText := repl.CleanInput(userText)
		if len(cleanText) != 0 {
			firstWord := cleanText[0]
			callFunc := repl.GetRegistry()[firstWord].Callback
			if callFunc != nil {
				callFunc(&config)
				fmt.Println("")
				continue
			}
		}
		fmt.Println(fmt.Errorf("Unknown command"))
	}
}




/*
if len == 0 or call Func != nil{
	fmt.PrintError
}
=> if len != 0 and Func == nil{

}
*/