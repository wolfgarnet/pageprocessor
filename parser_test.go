package html2

import (
	"testing"
	"fmt"
)

func TestParser(t *testing.T) {
	src := `<a href="url">text<img src="imgurl"></a>`
	parser, err := ParseHtml(src)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	println("IMAGE", parser.Images[0].Url)
}
