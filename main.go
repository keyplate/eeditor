package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/keyplate/eeditor/resources/fonts"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	jetBrainsMonoFaceSource *text.GoTextFaceSource
	cursorImg               *ebiten.Image
	fontSize                float64
    lineSpacing             float64
)

func init() {
	fontSize = 14
    lineSpacing = float64(fontSize) * 1.2

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.JetBrainsMonoRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	jetBrainsMonoFaceSource = s

	cursorImg = ebiten.NewImage(2, 14)
	cursorImg.Fill(color.White)
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

type Game struct {
	runes   []rune
	text    string
	counter int
	cursor  Cursor
}

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.text += string(g.runes)
	g.cursor.updateCursorMap(g.text)

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.text += "\n"
	}

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.text) >= 1 {
			g.text = g.text[:len(g.text)-1]
		}
	}

	if repeatingKeyPressed(ebiten.KeyArrowUp) {
		g.cursor.moveCursorUp()
	}

	if repeatingKeyPressed(ebiten.KeyArrowDown) {
		g.cursor.moveCursorDown()
	}
    if repeatingKeyPressed(ebiten.KeyArrowLeft) {
        g.cursor.moveCursorLeft()
    }
    if repeatingKeyPressed(ebiten.KeyArrowRight) {
        g.cursor.moveCursorRight()
    }

	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.text

	txtOp := &text.DrawOptions{}
	txtOp.LineSpacing = lineSpacing
	txtFace := &text.GoTextFace{
		Source: jetBrainsMonoFaceSource,
		Size:   fontSize,
	}

	text.Draw(screen, t, txtFace, txtOp)

	curOp := &ebiten.DrawImageOptions{}
	curOp.GeoM.Translate(
        g.cursor.getGraphicalX(txtFace),
        g.cursor.getGraphicalY(lineSpacing),
    )
	screen.DrawImage(cursorImg, curOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		text:    "",
		counter: 0,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EEditero")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
