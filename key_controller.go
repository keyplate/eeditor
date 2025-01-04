package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func repeatingKeyPressed(key ebiten.Key) bool {
	const delay = 30
	const interval = 3

	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}

	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (g *Game) enterPressed() {
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.gapBuffer.Insert("\n")
		g.cursor.updateCursorMap(g.gapBuffer.String())
		g.cursor.moveCursorRight()
	}
}

func (g *Game) backspacePressed() {
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		g.cursor.moveCursorLeft()
		g.gapBuffer.LeftDel()
	}
}

func (g *Game) arrowUpPressed() {
	if repeatingKeyPressed(ebiten.KeyArrowUp) {
		g.cursor.moveCursorUp()
		g.gapBuffer.UpMv()
	}
}

func (g *Game) arrowDownPressed() {
	if repeatingKeyPressed(ebiten.KeyArrowDown) {
		g.cursor.moveCursorDown()
		g.gapBuffer.DownMv()
	}
}

func (g *Game) arrowLeftPressed() {
	if repeatingKeyPressed(ebiten.KeyArrowLeft) {
		g.cursor.moveCursorLeft()
		g.gapBuffer.LeftMv()
	}
}

func (g *Game) arrowRightPressed() {
	if repeatingKeyPressed(ebiten.KeyArrowRight) {
		g.cursor.moveCursorRight()
		g.gapBuffer.RightMv()
	}
}
