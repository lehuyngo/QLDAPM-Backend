package url_parser

import (
	"mvdan.cc/xurls/v2"
)

/*
func ExtractURL(content string) []string {
	xurlsStrict := xurls.Strict()
	result := xurlsStrict.FindAllString("golangbyexample.com is https://golangbyexample.com?abc=123&def=234", -1)
	fmt.Println(result)

	xurlsRelaxed := xurls.Relaxed()
	result = xurlsRelaxed.FindAllString("The website is golangbyexample.com", -1)
	fmt.Println(result)

	result = xurlsRelaxed.FindAllString("golangbyexample.com is https://golangbyexample.com", -1)
	fmt.Println(result)

	return result
}
*/

func ExtractURL(content string) []string {
	return xurls.Strict().FindAllString(content, -1)
}
