package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	itn "github.com/islu/Eccl.7.8/internal"
)

func main() {

	app := itn.NewApp()
	ebiten.SetWindowTitle("Eccl.7.8")
	ebiten.SetWindowSize(app.ScreenWidth*2, app.ScreenHeight*2)
	// ebiten.SetWindowPosition(0, 0)
	// ebiten.SetFullscreen(true)
	ebiten.SetWindowFloating(true)
	// ebiten.SetWindowDecorated(false)

	err := ebiten.RunGameWithOptions(app, &ebiten.RunGameOptions{
		// ScreenTransparent: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}
