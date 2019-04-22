package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"tiki-go/src/tools/lex"
)

func main() {
	byteValue, err := ioutil.ReadFile("lex.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var tokens lex.Tokens
	xml.Unmarshal(byteValue, &tokens)

	fmt.Printf("XMLName %v\n", tokens.XMLName)
	for _, token := range tokens.Tokens {
		fmt.Printf("%v %v\n", token.Token, token.Regx)
	}
}
