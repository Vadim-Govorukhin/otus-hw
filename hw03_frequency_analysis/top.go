package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type FrequencyAnalizer struct {
	freqMap map[string]int
	top     []string
}

func (f FrequencyAnalizer) Len() int {
	return len(f.top)
}

func (f FrequencyAnalizer) Swap(i, j int) {
	f.top[i], f.top[j] = f.top[j], f.top[i]
}

func (f FrequencyAnalizer) Less(i, j int) bool {
	if f.freqMap[f.top[i]] == f.freqMap[f.top[j]] {
		return f.top[i] < f.top[j]
	}
	return f.freqMap[f.top[i]] > f.freqMap[f.top[j]]
}

func Top10(str string) []string {
	freqAn := FrequencyAnalizer{
		freqMap: make(map[string]int),
	}
	splitted := strings.FieldsFunc(strings.ToLower(str), separatorRule)
	for _, word := range splitted {
		if word == "-" {
			continue
		}
		val, ok := freqAn.freqMap[word]
		if ok {
			val++
		} else {
			val = 1
		}
		freqAn.freqMap[word] = val
	}

	// fmt.Println(f.freqMap)

	for key := range freqAn.freqMap {
		freqAn.top = append(freqAn.top, key)
	}
	sort.Sort(freqAn)

	resultLen := 10
	if len(freqAn.top) < resultLen {
		resultLen = len(freqAn.top)
	}
	return freqAn.top[:resultLen]
}

func separatorRule(c rune) bool {
	return unicode.IsSpace(c) || (unicode.IsPunct(c) && c != rune('-'))
}
