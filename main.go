package main

import (
	"bufio"
	"os"
	"fmt"
	// "time"
	"pokedex/internal/repl"
	// "pokedex/internal/pokecache"
	"strings"
)

func main() {
	config := repl.NewConfig()

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
				argument := strings.Join(cleanText[1:], "-")
				// fmt.Println(fmt.Sprintf("%s", argument))
				err := callFunc(config, argument)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("")
				continue
			}
		}
		fmt.Println(fmt.Errorf("Unknown command\n"))
	}
}




/*
if len == 0 or call Func != nil{
	fmt.PrintError
}
=> if len != 0 and Func == nil{

}
*/