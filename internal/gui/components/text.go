package components

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/Text.ttf
var fontText []byte

var (
	TextSource *text.GoTextFaceSource
)

func init() {
	var err error

	TextSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontText))
	if err != nil {
		log.Fatalf("Error while loading Default font: %v", err)
	}
}

// Function to print text with the screen to draw, the text, the positions, size and the color
func DrawText(screen *ebiten.Image, msg string, x float64, y float64, size float64, clr color.Color) {
	face := &text.GoTextFace{
		Source: TextSource,
		Size:   size,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)

	text.Draw(screen, msg, face, op)
}

type TextWidget struct {
	X, Y     float32
	Size     float32
	Text     string
	FontType string
	Color    color.Color
	Visible  bool

	Blink      bool
	blinkTicks int
}

// Initialize the text to draw
func NewTextWidget(x, y, size float32, text string, clr color.Color) *TextWidget {
	return &TextWidget{
		X:        x,
		Y:        y,
		Size:     size,
		Text:     text,
		FontType: "default",
		Color:    clr,
		Visible:  true,
	}
}

// Set the text to change the content if needed
func (tw *TextWidget) SetText(newText string) {
	tw.Text = newText
}

func (tw *TextWidget) Update() {
	if !tw.Visible {
		return
	}

	if tw.Blink {
		tw.blinkTicks++
	}
}

func (tw *TextWidget) Draw(screen *ebiten.Image) {
	if !tw.Visible {
		return
	}

	if tw.Blink && (tw.blinkTicks/30)%2 == 0 {
		return
	}

	DrawText(screen, tw.Text, float64(tw.X), float64(tw.Y), float64(tw.Size), tw.Color)
}
