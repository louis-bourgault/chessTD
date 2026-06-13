package font

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// source - google fonts

//go:embed RobotoMono-Regular.ttf
var DefaultFont []byte

var FontSource *text.GoTextFaceSource
var Font text.Face

func init() {
	// Parse the raw bytes into an Ebitengine font source
	s, err := text.NewGoTextFaceSource(bytes.NewReader(DefaultFont))
	if err != nil {
		log.Fatal(err)
	}
	FontSource = s
	Font = &text.GoTextFace{
		Source: FontSource,
		Size:   16,
	}
}
