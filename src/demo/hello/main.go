package main

import (
	"fmt"
	"tools/lex"
)

func main() {
	des, err := lex.NewDescription("re.txt")
	if err != nil {
		fmt.Printf("error,%v\n", err)
		return
	}
	for _, v := range des.Lines {
		fmt.Printf("%v   =   %v   =   %v\n", v.Pos, v.Text, v.Script)
	}
	fmt.Println("done main.")
}
