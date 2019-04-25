package util

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
)

func ReadTextLines(path string) ([]string, error) {
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
func RegxNamedMap(r *regexp.Regexp, str string) map[string]string {
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
func MatchNamedMap(r *regexp.Regexp, match []string) map[string]string {
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
func StringBoolMapKeys(m map[string]bool) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool { return strings.Compare(keys[i], keys[j]) < 0 })
	return keys
}
func ToArrayString(values []interface{}) []string {
	var keys []string
	for _, v := range values {
		keys = append(keys, v.(string))
	}
	return keys
}
