package lexical

import (
	"encoding/xml"
	"errors"
	"unicode"
)

type Tokens struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token  `xml:"token"`
}

type Token struct {
	XMLName xml.Name `xml:"token"`
	Token   string   `xml:"token,attr"`
	Regx    string   `xml:"regx"`
}
type RegxElementScaner struct {
	regx  string
	cur   int
	len   int
	atEOF bool
}

func NewRegxElementScaner(regx string) *RegxElementScaner {
	var newed RegxElementScaner
	newed.regx = regx
	newed.cur = 0
	newed.len = len(regx)
	return &newed
}
func (s *RegxElementScaner) Cur() (rune, error) {
	if s.cur < s.len {
		return rune(s.regx[s.cur]), nil
	}
	return -1, errors.New("out of character")
}
func (s *RegxElementScaner) Next() (rune, error) {
	r, err := s.Cur()
	if err == nil {
		s.cur++
	}
	return r, err
}
func (s *RegxElementScaner) Util(end rune) (string, error) {
	start := s.cur
	c, err := s.Next()
	for err == nil {
		c, err = s.Next()
		if c == end {
			return string(s.regx[start:s.cur]), err
		}
	}
	return "", err
}
func (s *RegxElementScaner) NextNonSpace() (rune, error) {
	c, err := s.Next()
	for err == nil && unicode.IsSpace(c) {
		c, err = s.Next()
	}
	return c, err
}
func (s *RegxElementScaner) NextItem() (kind string, data string, err error) {
	c, err := s.NextNonSpace()
	if err != nil {
		return "END", "", nil
	}
	switch c {
	case '|':
		return "OR", "", nil
	case '*':
		return "STAR", "", nil
	case '+':
		return "PLUS", "", nil
	case '?':
		return "QUESTION", "", nil
	case '.':
		return "DOT", "", nil
	case '(':
		return "LPAREN", "", nil
	case ')':
		return "RPAREN", "", nil
	case '[':
		data0, err0 := s.Util(']')
		if err0 != nil {
			return "ERR", "", nil
		}
		kind = "ARRAY"
		data = data0
		err = nil
		return
	case '\\':
		c, err := s.Next()
		if err != nil {
			return "ERR", "", nil
		}
		switch c {
		case 'd':
			return "DIGIT", "", nil
		case 'a':
			return "Alpha", "", nil
		case 'w':
			return "WORD", "", nil
		case 's':
			return "BLANK", "", nil
		default:
			return "CHAR", string(c), nil
		}
	default:
		return "CHAR", string(c), nil
	}
}
func (s *RegxElementScaner) getNextRegxElement() (string, string) {
	return "", ""
}
func (t *Token) getRegxSyntaxTree() SyntaxTree {

	return SyntaxTree{}
}
