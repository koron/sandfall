package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type buttonState struct {
	curr bool
	last bool
}

func (bs *buttonState) update(v bool) {
	bs.last = bs.curr
	bs.curr = v
}

func (bs *buttonState) isOn() bool {
	return bs.curr
}

func (bs *buttonState) triggerOn() bool {
	return bs.curr && !bs.last
}

func (bs *buttonState) triggerOff() bool {
	return !bs.curr && bs.last
}

type Game struct {
	sandbox *Sandbox

	currColor int
	addMode   int

	keys             []ebiten.Key
	mouseButtonLeft  buttonState
	mouseButtonRight buttonState
}

func (g *Game) Update() error {
	// fetch inputs
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, k := range g.keys {
		switch k {
		case ebiten.Key1:
			g.currColor = 1
		case ebiten.Key2:
			g.currColor = 2
		case ebiten.Key3:
			g.currColor = 3
		case ebiten.Key4:
			g.currColor = 4
		}
	}
	mx, my := ebiten.CursorPosition()
	g.mouseButtonLeft.update(ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	g.mouseButtonRight.update(ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))

	// update Game status

	g.sandbox.Update()

	switch g.addMode {
	case 0:
		if g.mouseButtonLeft.isOn() && g.sandbox.Get(mx, 119-my) == empty {
			g.sandbox.Set(mx, 119-my, Cell(g.currColor))
		}
	case 1:
		if g.mouseButtonLeft.triggerOn() {
			for dy := -4; dy < 3; dy++ {
				y := 119 - (my + dy)
				if y < 0 || y >= 120 {
					continue
				}
				for dx := -4; dx < 3; dx++ {
					x := mx + dx
					if x < 0 || x >= 80 {
						continue
					}
					g.sandbox.Set(x, y, Cell(g.currColor))
				}
			}
		}
	}

	if g.mouseButtonRight.triggerOff() {
		g.sandbox.Clear()
	}

	return nil
}

var colorPallet = []color.Color{
	color.Black,
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xff, 0xff},
	color.RGBA{0xff, 0xff, 0x00, 0xff},
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray16{0x6666})
	w, h := g.sandbox.Width(), g.sandbox.Height()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := int(g.sandbox.Get(x, y))
			if c >= len(colorPallet) {
				c = 0
			}
			screen.Set(x, 119-y, colorPallet[c])
		}
	}

	screen.Set(81, 1, colorPallet[g.currColor])
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth / 4, screenHeight / 4
}

func main() {
	game := &Game{
		sandbox:   NewSandbox(80, 120),
		currColor: 1,
		addMode:   1,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sandbox demo")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
