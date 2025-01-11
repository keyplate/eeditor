package main

import (
	"bytes"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/keyplate/eeditor/resources/fonts"

	gap "codeberg.org/Release-Candidate/go-gap-buffer"
)

type editorConfig struct {
	screenWidth             int
	screenHeight            int
	jetBrainsMonoFaceSource *text.GoTextFaceSource
	cursorImg               *ebiten.Image
	fontSize                float64
	lineSpacing             float64
}

type Game struct {
	runes     []rune
	gapBuffer gap.GapBuffer
	counter   int
	cursor    Cursor
	cfg       editorConfig
    file      *os.File
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
	txtOp.LineSpacing = g.cfg.lineSpacing
	txtFace := &text.GoTextFace{
		Source: g.cfg.jetBrainsMonoFaceSource,
		Size:   g.cfg.fontSize,
	}

	text.Draw(screen, t, txtFace, txtOp)
	if g.counter%60 < 30 {
		curOp := &ebiten.DrawImageOptions{}
		curOp.GeoM.Translate(
			g.cursor.getGraphicalX(txtFace),
			g.cursor.getGraphicalY(g.cfg.lineSpacing),
		)
		screen.DrawImage(g.cfg.cursorImg, curOp)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.screenWidth, g.cfg.screenHeight
}

func main() {
    file, err := getFile()
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.JetBrainsMonoRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	cursorImg := ebiten.NewImage(2, 14)
	cursorImg.Fill(color.White)

	cfg := editorConfig{
		fontSize:                14,
		lineSpacing:             16.8,
		jetBrainsMonoFaceSource: s,
		cursorImg:               cursorImg,
		screenWidth:             640,
		screenHeight:            480,
	}

	g := &Game{
		gapBuffer: *gap.New(),
		counter:   0,
		cfg:       cfg,
        file:      file,
	}
   
    content, err := loadFile(file)
    if err != nil {
        log.Fatal(err)
    }
    g.gapBuffer.Insert(content)
    
	ebiten.SetWindowSize(g.cfg.screenWidth, g.cfg.screenHeight)
	ebiten.SetWindowTitle("Editor")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
