package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var whitePixel *ebiten.Image

func getWhitePixel() *ebiten.Image {
	if whitePixel == nil {
		whitePixel = ebiten.NewImage(1, 1)
		whitePixel.Fill(color.White)
	}
	return whitePixel
}

func DrawFilledRect(dst *ebiten.Image, x, y, w, h float32, clr color.Color) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(h))
	op.GeoM.Translate(float64(x), float64(y))

	r, g, b, a := clr.RGBA()
	if a > 0 {
		op.ColorScale.Scale(
			float32(r)/65535.0,
			float32(g)/65535.0,
			float32(b)/65535.0,
			float32(a)/65535.0,
		)
		dst.DrawImage(getWhitePixel(), op)
	}
}

func DrawStrokeRect(dst *ebiten.Image, x, y, w, h, strokeWidth float32, clr color.Color) {
	DrawFilledRect(dst, x, y, w, strokeWidth, clr)               // up
	DrawFilledRect(dst, x, y+h-strokeWidth, w, strokeWidth, clr) // down
	DrawFilledRect(dst, x, y, strokeWidth, h, clr)               // left
	DrawFilledRect(dst, x+w-strokeWidth, y, strokeWidth, h, clr) // right
}
