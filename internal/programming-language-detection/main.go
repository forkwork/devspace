package main

import (
	"os"

	"dev.khulnasoft.com/programming-language-detection/pkg/detector"
)

func main() {
	lang := detector.GetLanguage(os.Args[1], 5000)
	println(lang)
}
