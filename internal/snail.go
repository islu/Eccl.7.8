package internal

import (
	"bytes"
	"log"

	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed RedSnail.png
var RedSnail []byte

type Anim struct {
	Sprites []Sprite
	Period  int
}

type Sprite struct {
	Image  *ebiten.Image
	OX     int
	OY     int
	Width  int
	Height int
}

type Snail struct {
	State string
	Anim  []Anim
	PosX  int
	PosY  int
}

func NewSnail() *Snail {
	img, _, err := image.Decode(bytes.NewReader(RedSnail))
	if err != nil {
		log.Fatal(err)
	}
	image := ebiten.NewImageFromImage(img)

	walk1 := Sprite{image, 0, 40, 40, 40}
	walk2 := Sprite{image, 44, 40, 40, 40}
	walk3 := Sprite{image, 96, 40, 50, 40}
	walk4 := Sprite{image, 156, 40, 54, 40}

	walk := Anim{
		Sprites: []Sprite{walk1, walk2, walk3, walk4},
		Period:  20,
	}

	return &Snail{
		State: "walk",
		Anim:  []Anim{walk},
		PosX:  50,
		PosY:  400,
	}
}

func (s *Snail) IsOnClick(mouseX, mouseY int) bool {
	return s.PosX <= mouseX && mouseX <= s.PosX+40 && s.PosY <= mouseY && mouseY <= s.PosY+40
}

func (s *Snail) Draw(screen *ebiten.Image, count int) {

	anim := s.Anim[0]

	i := (count / anim.Period) % len(anim.Sprites)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.PosX), float64(s.PosY))

	sprite := anim.Sprites[i]

	screen.DrawImage(sprite.Image.SubImage(image.Rect(sprite.OX, sprite.OY, sprite.OX+sprite.Width, sprite.OY+sprite.OY)).(*ebiten.Image), op)
}
