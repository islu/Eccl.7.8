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
	"github.com/islu/bard-sdk-go/bard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	//go:embed bard.png
	bardImage []byte
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
	bot          *bard.Chatbot
	font         font.Face
	snail        *Snail
	canPrompt    bool
}

func NewApp() *App {
	img, _, err := image.Decode(bytes.NewReader(bardImage))
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

	bot, err := bard.NewChatbot(os.Getenv("BARD_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	sw, sh := ebiten.ScreenSizeInFullscreen()

	return &App{
		ScreenWidth:  sw / 2,
		ScreenHeight: sh / 2,
		promptHint:   "Enter a prompt here:",
		bard:         ebiten.NewImageFromImage(img),
		bot:          bot,
		font:         font,
		snail:        NewSnail(),
		canPrompt:    false,
	}
}

func (g *App) Update() error {

	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if g.snail.IsOnClick(mx, my) {
			g.canPrompt = !g.canPrompt
		}
	}

	if g.canPrompt {
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

		if inpututil.IsKeyJustPressed(ebiten.KeyControl) {

			go func(prompt string) {
				resp, err := g.bot.Ask(prompt)
				if err != nil {
					log.Fatal(err)
				}

				c := strings.Split(resp.Content, "\n")

				for i := 0; i < len(c); i++ {
					c[i] = addNewlines(c[i])
				}

				g.content = strings.Join(c, "\n")
				g.prompt = ""
			}(g.prompt)

			g.prompt = "等待中..."
		}

		// If the backspace key is pressed, remove one character.
		if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(g.prompt) >= 1 {
				g.prompt = g.prompt[:len(g.prompt)-1]
			}
		}
	}

	g.counter++
	return nil
}

func (g *App) Draw(screen *ebiten.Image) {

	if g.canPrompt {
		ebitenutil.DebugPrintAt(screen, g.promptHint, 0, 300)

		// Blink the cursor.
		t := g.prompt
		if g.counter%60 < 30 {
			t += "_"
		}
		text.Draw(screen, t, g.font, 0, 325, color.White)

		// Print content
		text.Draw(screen, g.content, g.font, 200, 25, color.White)
	}

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(0.05, 0.05)
	// op.GeoM.Translate(0, 64)
	// screen.DrawImage(g.bard, op)

	// text.Draw(screen, "hello 123 中文測試", g.font, 0, 30, color.White)

	g.snail.Draw(screen, g.counter)

	// mx, my := ebiten.CursorPosition()
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mx: %d, my: %d", mx, my), 0, 25)
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("hover: %v", g.snail.IsOnClick(mx, my)), 0, 50)
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

func addNewlines(s string) string {
	b := strings.Builder{}
	count := 0
	for _, ss := range s {
		count++
		if count%30 == 0 {
			b.WriteString("\n")
		}
		b.WriteRune(ss)
	}
	return b.String()
}
