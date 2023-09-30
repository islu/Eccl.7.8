package main

import (
	"log"
	"os"

	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	bard *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.05, 0.05)
	op.GeoM.Translate(0, 64)
	screen.DrawImage(g.bard, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	f, err := os.Open("bard.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bard := ebiten.NewImageFromImage(img)

	g := &Game{
		bard: bard,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Eccl.7.8")
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
