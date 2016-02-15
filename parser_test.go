package html2

import (
	"testing"
	"fmt"
)

func TestParser(t *testing.T) {
	src := `<a href="url">text<img src="imgurl"></a><title>JJAJA</title>`
	parser, err := ParseHtml(src)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	fmt.Printf("TITLE: %v\n", parser.Title)

	println("IMAGE", parser.Images[0].URL.String())
}
