package main

import (
	"math/rand/v2"
)

type Cell uint8

const (
	empty Cell = 0x00
	wall  Cell = 0xff
)

type Sandbox struct {
	w   int
	h   int
	buf []Cell
}

func NewSandbox(w, h int) *Sandbox {
	if w <= 0 || h <= 0 {
		panic("unsupport zero width nor height")
	}
	// Add two for wall on all four sides.
	w = w + 2
	h = h + 2
	sandbox := &Sandbox{w: w, h: h, buf: make([]Cell, w*h)}
	sandbox.initWall()
	return sandbox
}

func (box *Sandbox) initWall() {
	// Surround a rectangle buffer with wall (0xff) on all four sides.
	for x := 0; x < box.w; x++ {
		box.setCell(x, 0, wall)
		box.setCell(x, box.h-1, wall)
	}
	for y := 0; y < box.h; y++ {
		box.setCell(0, y, wall)
		box.setCell(box.w-1, y, wall)
	}
}

func (box *Sandbox) Width() int { return box.w - 2 }

func (box *Sandbox) Height() int { return box.h - 2 }

func (box *Sandbox) Get(x, y int) Cell {
	if x < 0 || x > box.w-2 || y < 0 || y > box.h-2 {
		return wall
	}
	return box.getCell(x+1, y+1)
}

func (box *Sandbox) Set(x, y int, c Cell) {
	if x < 0 || x > box.w-2 || y < 0 || y > box.h-2 {
		return
	}
	box.setCell(x+1, y+1, c)
}

func (box *Sandbox) index(x, y int) int {
	return x + y*box.w
}

func (box *Sandbox) getCell(x, y int) Cell {
	return box.buf[box.index(x, y)]
}

func (box *Sandbox) setCell(x, y int, c Cell) {
	box.buf[box.index(x, y)] = c
}

// fallDir determines whether the cell falls to the left or right.
// 0:not fall, 1:left, 2:right
func (box *Sandbox) fallDir(wallLeft, wallRight bool) int {
	// Can't fall
	if wallLeft && wallRight {
		return 0
	}
	// Can fall to only bottom-left
	if !wallLeft && wallRight {
		return 1
	}
	// Can fall to only bottom-right
	if wallLeft && !wallRight {
		return 2
	}
	// Can fall on either side, so it is randomly selected.
	return int(rand.Int64())%2 + 1
}

func (box *Sandbox) Update() int {
	updated := 0
	// Fall all cells
	for y := 1; y < box.h-1; y++ {
		for x := 1; x < box.w-1; x++ {
			c := box.getCell(x, y)
			if c == empty || c == wall {
				continue
			}

			// Fall to bottom
			if b := box.getCell(x, y-1); b == empty {
				box.setCell(x, y-1, c)
				box.setCell(x, y, empty)
				updated++
				continue
			}

			// Check if it can be dropped to the bottom-left or the
			// bottom-right
			//
			// Make asymmetric decisions for left and right. The left cell is
			// dropped in the previous loop, so only the bottom-left cell needs
			// to be checked. The right cell is not yet dropped, so the
			// bottom-right cell needs to be checked along with it.
			wallLeft := box.getCell(x-1, y-1) != empty
			wallRight := box.getCell(x+1, y-1) != empty || box.getCell(x+1, y) != empty
			switch box.fallDir(wallLeft, wallRight) {
			case 0:
				continue
			case 1:
				box.setCell(x-1, y-1, c)
			case 2:
				box.setCell(x+1, y-1, c)
			}
			box.setCell(x, y, empty)
			updated++
		}
	}
	return updated
}

func (box *Sandbox) Clear() {
	for i := range box.buf {
		box.buf[i] = empty
	}
	box.initWall()
}
