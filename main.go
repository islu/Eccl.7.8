package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	itn "github.com/islu/Eccl.7.8/internal"
)

func main() {

	app := itn.NewApp()

	ebiten.SetWindowSize(app.ScreenWidth*2, app.ScreenHeight*2)
	ebiten.SetWindowTitle("Eccl.7.8")
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
