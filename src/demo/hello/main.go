package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"tools/lex"
)

func main() {
	// Open our xmlFile
	xmlFile, err := os.Open("lex.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var tokens lex.Tokens
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &tokens)

	fmt.Printf("XMLName %v\n", tokens.XMLName)
	for _, token := range tokens.Tokens {
		fmt.Printf("XMLName %v %v\n", token.XMLName, token)
	}
}
