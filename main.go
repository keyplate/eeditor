package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/keyplate/eeditor/resources/fonts"

	gap "codeberg.org/Release-Candidate/go-gap-buffer"
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
	runes     []rune
	gapBuffer gap.GapBuffer
	counter   int
	cursor    Cursor
}

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.gapBuffer.Insert(string(g.runes))

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.gapBuffer.Insert("\n")
		g.cursor.updateCursorMap(g.gapBuffer.String())
		g.cursor.moveCursorRight()
	}

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		g.cursor.moveCursorLeft()
		g.gapBuffer.LeftDel()
	}

	g.cursor.updateCursorMap(g.gapBuffer.String())
	for _, _ = range g.runes {
		g.cursor.moveCursorRight()
	}

	if repeatingKeyPressed(ebiten.KeyArrowUp) {
		g.cursor.moveCursorUp()
		g.gapBuffer.UpMv()
	}
	if repeatingKeyPressed(ebiten.KeyArrowDown) {
		g.cursor.moveCursorDown()
		g.gapBuffer.DownMv()
	}
	if repeatingKeyPressed(ebiten.KeyArrowLeft) {
		g.cursor.moveCursorLeft()
		g.gapBuffer.LeftMv()
	}
	if repeatingKeyPressed(ebiten.KeyArrowRight) {
		g.cursor.moveCursorRight()
		g.gapBuffer.RightMv()
	}

	g.counter++
	if g.counter > 61 {
		g.counter = 0
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.gapBuffer.String()

	txtOp := &text.DrawOptions{}
	txtOp.LineSpacing = lineSpacing
	txtFace := &text.GoTextFace{
		Source: jetBrainsMonoFaceSource,
		Size:   fontSize,
	}

	text.Draw(screen, t, txtFace, txtOp)
	if g.counter%60 < 30 {
		curOp := &ebiten.DrawImageOptions{}
		curOp.GeoM.Translate(
			g.cursor.getGraphicalX(txtFace),
			g.cursor.getGraphicalY(lineSpacing),
		)
		screen.DrawImage(cursorImg, curOp)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		gapBuffer: *gap.New(),
		counter:   0,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EEditero")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
