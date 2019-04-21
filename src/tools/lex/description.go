package lex

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Line struct {
	Pos    string
	Text   string
	Script string
}

type Description struct {
	Lines []Line
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func reNamedMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
func NewDescription(path string) (*Description, error) {
	des := new(Description)

	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}

	var myExp = regexp.MustCompile(`^(?P<text>.+)(?P<script>\{.+\})`)
	for i, v := range lines {
		ret := reNamedMap(myExp, v)
		if ret != nil {
			var line Line
			line.Pos = strconv.Itoa(i)
			line.Text = strings.TrimSpace(ret["text"])
			line.Script = strings.TrimSpace(ret["script"])

			des.Lines = append(des.Lines, line)
		}
	}
	return des, nil
}
