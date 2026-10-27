package components

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/Title.ttf
var fontTitle []byte

//go:embed assets/Text.ttf
var fontText []byte

var (
	TitleSource *text.GoTextFaceSource
	TextSource  *text.GoTextFaceSource
)

func init() {
	var err error

	TitleSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontTitle))
	if err != nil {
		log.Fatalf("Error while loading Title font: %v", err)
	}

	TextSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontText))
	if err != nil {
		log.Fatalf("Error while loading Default font: %v", err)
	}
}

func drawText(screen *ebiten.Image, msg string, x float64, y float64, size float64, clr color.Color, fontType string) {
	var source *text.GoTextFaceSource

	switch fontType {
	case "title":
		source = TitleSource
	default:
		source = TextSource
	}

	face := &text.GoTextFace{
		Source: source,
		Size:   size,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)

	text.Draw(screen, msg, face, op)
}
