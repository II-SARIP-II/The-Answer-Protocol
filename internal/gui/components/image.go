package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image struct {
	Path       string
	X, Y, W, H float32
	img        *ebiten.Image
}

func NewImage(path string, x, y, w, h float32) (*Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from %s: %w", path, err)
	}

	return &Image{
		Path: path,
		X:    x,
		Y:    y,
		W:    w,
		H:    h,
		img:  img,
	}, nil
}

func (i *Image) Draw(screen *ebiten.Image) {
	if i.img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	bounds := i.img.Bounds()
	srcW := float32(bounds.Dx())
	srcH := float32(bounds.Dy())

	if srcW > 0 && srcH > 0 {
		scaleX := float64(i.W / srcW)
		scaleY := float64(i.H / srcH)
		op.GeoM.Scale(scaleX, scaleY)
	}

	op.GeoM.Translate(float64(i.X), float64(i.Y))

	screen.DrawImage(i.img, op)
}

func (i *Image) Update() {
}
