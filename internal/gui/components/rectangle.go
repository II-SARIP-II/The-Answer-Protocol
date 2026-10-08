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

const doorSize float32 = 20.0

func DrawStrokeRect(dst *ebiten.Image, x, y, w, h, strokeWidth float32, clr color.Color, holes []Hole) {
	drawWallHorizontal(dst, x, y, w, strokeWidth, clr, 0, holes)               // north
	drawWallHorizontal(dst, x, y+h-strokeWidth, w, strokeWidth, clr, 2, holes) // south

	drawWallVertical(dst, x, y, strokeWidth, h, clr, 3, holes)               // west
	drawWallVertical(dst, x+w-strokeWidth, y, strokeWidth, h, clr, 1, holes) //est
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

func drawWallHorizontal(dst *ebiten.Image, wallX, wallY, wallWidth, strokeWidth float32, clr color.Color, wallID int, holes []Hole) {
	var holePos *float32
	for _, h := range holes {
		if h.Wall == wallID {
			pos := h.Pos
			holePos = &pos
			break
		}
	}

	if holePos == nil {
		DrawFilledRect(dst, wallX, wallY, wallWidth, strokeWidth, clr)
		return
	}

	seg1Width := *holePos
	if seg1Width > 0 {
		DrawFilledRect(dst, wallX, wallY, seg1Width, strokeWidth, clr)
	}

	seg2X := wallX + *holePos + doorSize
	seg2Width := (wallX + wallWidth) - seg2X
	if seg2Width > 0 {
		DrawFilledRect(dst, seg2X, wallY, seg2Width, strokeWidth, clr)
	}
}

func drawWallVertical(dst *ebiten.Image, wallX, wallY, strokeWidth, wallHeight float32, clr color.Color, wallID int, holes []Hole) {
	var holePos *float32
	for _, h := range holes {
		if h.Wall == wallID {
			pos := h.Pos
			holePos = &pos
			break
		}
	}

	if holePos == nil {
		DrawFilledRect(dst, wallX, wallY, strokeWidth, wallHeight, clr)
		return
	}

	seg1Height := *holePos
	if seg1Height > 0 {
		DrawFilledRect(dst, wallX, wallY, strokeWidth, seg1Height, clr)
	}

	seg2Y := wallY + *holePos + doorSize
	seg2Height := (wallY + wallHeight) - seg2Y
	if seg2Height > 0 {
		DrawFilledRect(dst, wallX, seg2Y, strokeWidth, seg2Height, clr)
	}
}
