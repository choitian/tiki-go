package syntax

import (
	"tools/util"
)

type GrammerDescription struct {
	content []string
}

func NewGrammerDescription(path string) *GrammerDescription {
	var gd GrammerDescription
	var err error
	gd.content, err = util.ReadTextLines(path)
	if err != nil {
		panic(err)
	}
	return &gd
}
