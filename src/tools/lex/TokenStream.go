package lex

type Token struct {
	tok  int
	data string
	info string
}

type TokenStream interface {
	NextToken() Token
}
