package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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

type Game struct {
	runes     []rune
	gapBuffer gap.GapBuffer
	counter   int
	cursor    Cursor
}

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	if len(g.runes) > 0 {
		g.gapBuffer.Insert(string(g.runes))
	}

	g.cursor.updateCursorMap(g.gapBuffer.String())
	for range g.runes {
		g.cursor.moveCursorRight()
	}

	g.enterPressed()
	g.backspacePressed()
	g.arrowUpPressed()
	g.arrowDownPressed()
	g.arrowLeftPressed()
	g.arrowRightPressed()

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
	ebiten.SetWindowTitle("Editor")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
