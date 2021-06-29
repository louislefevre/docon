package internal

import (
	"regexp"
	"strings"
)

type keywordType string

type keyword struct {
	word   string
	kwType keywordType
}

func newKeyword(kwType keywordType) keyword {
	return keyword{string(kwType), kwType}
}

func (kw keyword) transform(text string, newText string) string {
	text = strings.ReplaceAll(text, strings.ToLower(kw.word), strings.ToLower(newText))
	text = strings.ReplaceAll(text, strings.ToUpper(kw.word), strings.ToUpper(newText))
	text = strings.ReplaceAll(text, strings.Title(kw.word), strings.Title(newText))
	return text
}

type keywordSet []keyword

func (kws keywordSet) contains(text string) bool {
	for _, kw := range kws {
		if text == strings.ToLower(kw.word) {
			return true
		} else if text == strings.ToUpper(kw.word) {
			return true
		} else if text == strings.Title(kw.word) {
			return true
		}
	}
	return false
}

func containsValidKeywords(text string, set keywordSet) (string, bool) {
	re := regexp.MustCompile(`\{(.*?)\}`)
	submatchall := re.FindAllString(text, -1)
	for _, element := range submatchall {
		if !set.contains(element) {
			return element, false
		}
	}
	return "", true
}
