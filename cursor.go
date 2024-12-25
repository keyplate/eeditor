package main

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Cursor struct {
	x          int
	y          int
	preferredX int
	cursorMap  []int
}

func (c *Cursor) moveCursorUp() {
    //If cursor is in the top do nothing
    if c.y == 0 {
		return
	}
	//If above is too short to jump directly above, we jump to the end of the line
	if c.cursorMap[c.y-1] < c.x {
		c.x = c.cursorMap[c.y-1]
	} else {
    	c.x = c.preferredX
    } 
	c.y--
}

func (c *Cursor) moveCursorDown() {
    //If cursor is at the bottom do nothing
	if c.y == len(c.cursorMap) - 1 {
		return
	}
	//If a line is too short to jump directly below, we jump to the end of the line
	if c.cursorMap[c.y+1] < c.x {
		c.x = c.cursorMap[c.y+1]
	} else {
        c.x = c.preferredX
    }
    c.y++
}

func (c *Cursor) moveCursorLeft() {
    //If cursor is in the top left position do nothing
    if c.y == 0 && c.x == 0 {
        return
    }
    //If cursor is in the edge left position move it to the end of the line above
    if c.x == 0 {
        c.y--
        c.x = c.cursorMap[c.y]
    } else {
        c.x--
    }
    c.preferredX = c.x
}

func (c *Cursor) moveCursorRight() {
    //If cursor is in the bottom right position do nothing
    fmt.Print(c)
    if c.y == len(c.cursorMap) - 1 && c.x == c.cursorMap[len(c.cursorMap) - 1] {
        return
    }
    //If cursor is in the edge right position move it to the below of the line above
    if c.x == c.cursorMap[c.y] {
        c.y++
        c.x = 0
    } else {
        c.x++
    }
    c.preferredX = c.x
}

func (c *Cursor) updateCursorMap(text string) {
    lines := strings.Split(text, "\n")
	linesLen := make([]int, len(lines))

	for i, line := range lines {
		linesLen[i] = len(line)
	}
	c.cursorMap = linesLen
}

//Works only for monospaced fonts
func (c *Cursor) getGraphicalX(face text.Face) float64 { 
    return text.Advance(" ", face) * float64(c.x) 
}

func (c *Cursor) getGraphicalY(lineHeigth float64) float64 {
    return float64(c.y) * lineHeigth
}
