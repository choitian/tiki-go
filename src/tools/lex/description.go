package lex

import (
	"encoding/xml"
)

type Tokens struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token  `xml:"token"`
}

type Token struct {
	XMLName xml.Name `xml:"token"`
	Tok     string   `xml:"tok,attr"`
	Regx    string   `xml:"regx"`
}
