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

type Game struct {
	sandbox *Sandbox

	currCell int

	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, k := range g.keys {
		switch k {
		case ebiten.Key1:
			g.currCell = 1
		case ebiten.Key2:
			g.currCell = 2
		case ebiten.Key3:
			g.currCell = 3
		case ebiten.Key4:
			g.currCell = 4
		}
	}

	g.sandbox.Update()

	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.sandbox.Set(mx, 119-my, Cell(g.currCell))
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.sandbox.Clear()
	}

	return nil
}

var cellPallet = []color.Color{
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
			if c >= len(cellPallet) {
				c = 0
			}
			screen.Set(x, 119-y, cellPallet[c])
		}
	}

	screen.Set(81, 1, cellPallet[g.currCell])
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth / 4, screenHeight / 4
}

func main() {
	game := &Game{
		sandbox:  NewSandbox(80, 120),
		currCell: 1,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sandbox demo")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
