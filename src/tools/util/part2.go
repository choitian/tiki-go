package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"tools/lexical"
)

func test() {
	byteValue, err := ioutil.ReadFile("lexical.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var tokens lexical.Tokens
	xml.Unmarshal(byteValue, &tokens)

	fmt.Printf("XMLName %v\n", tokens.XMLName)
	for _, token := range tokens.Tokens {
		fmt.Printf("%v %v\n", token.Token, token.Regx)
	}
}
