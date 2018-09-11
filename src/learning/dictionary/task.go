package dictionary

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	words := strings.Split(s, " ")
	wc := make(map[string]int)
	for _, word := range words {
		count, ok := wc[word]
		if ok {
			wc[word] = count + 1
		} else {
			wc[word] = 1
		}
	}
	return wc
}

func Task() {
	wc.Test(WordCount)
}
