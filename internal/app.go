package internal

import (
	"bytes"
	"log"
	"strings"

	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed bard.png
var bard []byte

type App struct {
	ScreenWidth  int
	ScreenHeight int
	runes        []rune
	text         string
	counter      int
	bard         *ebiten.Image
}

func NewApp() *App {
	img, _, err := image.Decode(bytes.NewReader(bard))
	if err != nil {
		log.Fatal(err)
	}

	sw, sh := ebiten.ScreenSizeInFullscreen()

	return &App{
		ScreenWidth:  sw / 2,
		ScreenHeight: sh / 2,
		text:         "Enter a prompt here:\n",
		counter:      0,
		bard:         ebiten.NewImageFromImage(img),
	}
}

func (g *App) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.text += string(g.runes)

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(g.text, "\n")
	if len(ss) > 10 {
		g.text = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.text += "\n"
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.text) >= 1 {
			g.text = g.text[:len(g.text)-1]
		}
	}

	g.counter++
	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	// Blink the cursor.
	t := g.text
	if g.counter%60 < 30 {
		t += "_"
	}
	ebitenutil.DebugPrint(screen, t)

	// Draw image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.05, 0.05)
	op.GeoM.Translate(0, 64)
	screen.DrawImage(g.bard, op)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
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
