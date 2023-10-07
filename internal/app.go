package internal

import (
	"bytes"
	"log"
	"os"
	"strings"

	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	//go:embed bard.png
	bard []byte
	//go:embed Cubic_11_1.013_R.ttf
	cubicFont []byte
)

type App struct {
	ScreenWidth  int
	ScreenHeight int
	runes        []rune
	promptHint   string
	prompt       string
	content      string
	counter      int
	bard         *ebiten.Image
	bot          *Chatbot
	font         font.Face
	snail        *Snail
}

func NewApp() *App {
	img, _, err := image.Decode(bytes.NewReader(bard))
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(cubicFont)
	if err != nil {
		log.Fatal(err)
	}

	font, err := opentype.NewFace(tt, nil)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := NewChatbot(os.Getenv("BARD_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	sw, sh := ebiten.ScreenSizeInFullscreen()

	return &App{
		ScreenWidth:  sw / 4,
		ScreenHeight: sh / 4,
		promptHint:   "Enter a prompt here:",
		bard:         ebiten.NewImageFromImage(img),
		bot:          bot,
		font:         font,
		snail:        NewSnail(),
	}
}

func (g *App) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.prompt += string(g.runes)

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(g.prompt, "\n")
	if len(ss) > 10 {
		g.prompt = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		// resp, err := g.bot.Ask(g.prompt)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// g.content = resp.Content

		// g.prompt = ""
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.prompt) >= 1 {
			g.prompt = g.prompt[:len(g.prompt)-1]
		}
	}

	g.counter++
	return nil
}

func (g *App) Draw(screen *ebiten.Image) {

	ebitenutil.DebugPrint(screen, g.promptHint)

	// Blink the cursor.
	t := g.prompt
	if g.counter%60 < 30 {
		t += "_"
	}
	ebitenutil.DebugPrintAt(screen, t, 0, 20)

	// Print content
	ebitenutil.DebugPrintAt(screen, g.content, 200, 0)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.05, 0.05)
	op.GeoM.Translate(0, 64)
	screen.DrawImage(g.bard, op)

	text.Draw(screen, "hello 123 中文測試", g.font, 0, 30, color.White)

	g.snail.Draw(screen, g.counter)
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
