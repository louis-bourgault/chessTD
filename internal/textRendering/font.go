package textRendering

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	ebitenText "github.com/hajimehoshi/ebiten/v2/text/v2"
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

func DrawText(screen *ebiten.Image, text string, x, y float64, size float64, colorScale color.Color) {
	face := &ebitenText.GoTextFace{
		Source: FontSource,
		Size:   size,
	}
	op := &ebitenText.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(colorScale)

	ebitenText.Draw(screen, text, face, op)
}
