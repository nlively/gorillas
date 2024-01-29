package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"log"
	"time"
)

type Intro struct {
	last_millisecond int64
	blink_state      bool
}

func (i *Intro) blink() {
	millisecond := time.Now().UnixMilli()
	if i.blink_state && millisecond > i.last_millisecond+1000 || !i.blink_state && millisecond > i.last_millisecond+500 {
		i.blink_state = !i.blink_state
		i.last_millisecond = millisecond
	}
}

func (i *Intro) draw_intro(screen *ebiten.Image) {
	var err error
	var img *ebiten.Image

	img, _, err = ebitenutil.NewImageFromFile("images/intro.png")
	if err != nil {
		log.Fatal(err)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(200, 80)

	screen.DrawImage(img, op)

	if i.blink_state {
		draw_text(screen, "Press Enter to Start", 200, 360)
	}
}
