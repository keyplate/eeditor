package main

import (
	"bytes"
	"image/color"
	"log"
	"strings"

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
    fontSize int
    cursorImg *ebiten.Image
)

func init() {
    fontSize = 14

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

func getCursorMap(text string) []int {
    lines := strings.Split(text, "\n")
    linesLen := make([]int, len(lines)) 

    for i, line := range(lines) {
        linesLen[i] = len(line)
    }
    
    return linesLen
}

type Game struct {
	runes   []rune
	text    string
	counter int
    cursorMap []int
    cursorPosX int
    cursorPosY int
}

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.text += string(g.runes)
    
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.text += "\n"
	}

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.text) >= 1 {
			g.text = g.text[:len(g.text)-1]
		}
	}

    g.cursorMap = getCursorMap(g.text) 
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.text
    
    txtOp := &text.DrawOptions{}
    txtOp.LineSpacing = 14 * 1.2
    txtFace := &text.GoTextFace{
        Source: jetBrainsMonoFaceSource,
        Size: 14,
    } 

    text.Draw(screen, t, txtFace, txtOp)

    curOp := &ebiten.DrawImageOptions{}
    curOp.GeoM.Translate(float64(g.cursorPosX), float64(g.cursorPosY))
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
