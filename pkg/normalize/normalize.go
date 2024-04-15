package normalize

import (
	"net/url"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NormalizeText(rawText string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, rawText)
	if err != nil {
		return "", err
	}

	return result, nil
}

func URLDecodeV1(str string) (string, error) {
	return url.QueryUnescape(str)
}

func URLEncode(str string) string {
	return url.QueryEscape(str)
}