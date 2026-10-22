package components

import (
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const MaxCharacterPerLine = 50
const MaxLines = 15

type ChatWidget struct {
	X, Y, W, H float32
	Lines      []string
	TextSize   float32
	Color      color.Color
	Spacing    float64
	InputTitle string
	Text       string
	Runes      []rune
	HasFocus   bool
	OnSubmit   func(text string)
	HasBtn     bool
	ModeButton *Button
	Mode       string
}

func NewChatWidget(x, y, w, h, textSize float32, col color.Color, inputTitle string, hasBtn bool) *ChatWidget {
	chat := &ChatWidget{
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		Lines:      []string{},
		TextSize:   textSize,
		Color:      color.RGBA{R: 0, G: 255, B: 180, A: 255},
		Spacing:    20.0,
		InputTitle: inputTitle,
		HasFocus:   false,
		HasBtn:     hasBtn,
		Mode:       "GROUP",
	}

	if hasBtn {
		inputHeight := float32(35.0)
		boxY := y + h - inputHeight - 5.0
		btnW := float32(80.0)

		chat.ModeButton = NewButton("GROUP", x, boxY, btnW, inputHeight)
		chat.ModeButton.OnClick = func() {
			if chat.Mode == "GROUP" {
				chat.Mode = "GLOBAL"
				chat.ModeButton.Label = "GLOBAL"
			} else {
				chat.Mode = "GROUP"
				chat.ModeButton.Label = "GROUP"
			}
		}
	}

	return chat
}

func (c *ChatWidget) AddLine(newLine string) {
	if len(newLine) > MaxCharacterPerLine {
		c.AddLine(newLine[:MaxCharacterPerLine])
		c.AddLine(newLine[MaxCharacterPerLine:])
		return
	}

	c.Lines = append(c.Lines, newLine)
	if len(c.Lines) > MaxLines {
		c.Lines = c.Lines[len(c.Lines)-MaxLines:]
	}
}

func (c *ChatWidget) Update() {
	inputHeight := float32(35.0)
	boxY := c.Y + c.H - inputHeight - 5.0

	inputX := c.X
	inputW := c.W

	if c.HasBtn && c.ModeButton != nil {
		c.ModeButton.Update()
		inputX = c.X + c.ModeButton.W + 5.0
		inputW = c.W - c.ModeButton.W - 5.0
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		cursorX, cursorY := float32(mx), float32(my)

		if cursorX >= inputX && cursorX <= inputX+inputW && cursorY >= boxY && cursorY <= boxY+inputHeight {
			c.HasFocus = true
		} else {
			c.HasFocus = false
		}
	}

	if !c.HasFocus {
		return
	}

	c.Runes = ebiten.AppendInputChars(c.Runes[:0])

	if len(c.Text) < 30 {
		c.Text += string(c.Runes)
	}

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		if len(c.Text) > 0 {
			if !c.HasBtn {
				c.AddLine(c.InputTitle + c.Text)
			}

			if c.HasBtn && c.OnSubmit != nil {
				c.OnSubmit("CHAT " + c.Mode + " " + c.Text)
			} else if c.OnSubmit != nil {
				c.OnSubmit(c.Text)
			}

			c.Text = ""
		}
	}

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(c.Text) > 0 {
			_, size := utf8.DecodeLastRuneInString(c.Text)
			c.Text = c.Text[:len(c.Text)-size]
		}
	}
}

func (c *ChatWidget) Draw(screen *ebiten.Image) {
	inputHeight := float32(35.0)

	for i, line := range c.Lines {
		lineY := float64(c.Y) + (float64(i) * c.Spacing)
		DrawText(screen, line, float64(c.X), lineY, float64(c.TextSize), c.Color, "Title")
	}

	boxY := c.Y + c.H - inputHeight - 5.0

	inputX := c.X
	inputW := c.W

	if c.HasBtn && c.ModeButton != nil {
		c.ModeButton.Draw(screen)
		inputX = c.X + c.ModeButton.W + 5.0
		inputW = c.W - c.ModeButton.W - 5.0
	}

	borderColor := color.RGBA{R: 0, G: 120, B: 150, A: 255}
	if c.HasFocus {
		borderColor = color.RGBA{R: 0, G: 255, B: 220, A: 255}
	}

	DrawFilledRect(screen, inputX, boxY, inputW, inputHeight, color.RGBA{R: 5, G: 15, B: 25, A: 200})
	DrawStrokeRect(screen, inputX, boxY, inputW, inputHeight, 1.0, borderColor, []Hole{})

	cursor := ""
	if c.HasFocus {
		cursor = "|"
	}

	displayText := c.Text + cursor

	DrawText(screen, displayText, float64(inputX+8), float64(boxY+8), float64(c.TextSize), color.RGBA{R: 0, G: 255, B: 180, A: 255}, "Title")
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
