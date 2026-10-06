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

//go:embed assets/Feather.ttf
var fontFeather []byte

//go:embed assets/Default.ttf
var fontDefault []byte

var (
	TitleSource   *text.GoTextFaceSource
	FeatherSource *text.GoTextFaceSource
	DefaultSource *text.GoTextFaceSource
)

func init() {
	var err error

	TitleSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontTitle))
	if err != nil {
		log.Fatalf("Error while loading Title font: %v", err)
	}

	FeatherSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontFeather))
	if err != nil {
		log.Fatalf("Error while loading Feather font: %v", err)
	}

	DefaultSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontDefault))
	if err != nil {
		log.Fatalf("Error while loading Default font: %v", err)
	}
}

func DrawText(screen *ebiten.Image, msg string, x float64, y float64, size float64, clr color.Color, fontType string) {
	var source *text.GoTextFaceSource

	switch fontType {
	case "title":
		source = TitleSource
	case "feather":
		source = FeatherSource
	default:
		source = DefaultSource
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
