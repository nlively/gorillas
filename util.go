package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
)

func draw_any_text(screen *ebiten.Image, message string, x int, y int, size float64, color color.Color) {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	text.Draw(screen, message, mplusNormalFont, x, y, color)
}

func draw_text(screen *ebiten.Image, message string, x int, y int) {
	draw_any_text(screen, message, x, y, 24, color.White)
}

func draw_smaller_text(screen *ebiten.Image, message string, x int, y int) {
	draw_any_text(screen, message, x, y, 16, color.White)
}

func draw_error_text(screen *ebiten.Image, message string, x int, y int) {
	draw_any_text(screen, message, x, y, 24, color.RGBA{0xFF, 0x00, 0x00, 0xFF})
}

func random_color() color.Color {
	red := uint8(rand.Intn(255))
	green := uint8(rand.Intn(255))
	blue := uint8(rand.Intn(255))
	return color.RGBA{red, green, blue, 0xFF}
}

func repeating_key_pressed(key ebiten.Key) bool {
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

func handle_backspace(input *string) {
	if repeating_key_pressed(ebiten.KeyBackspace) {
		if len(*input) > 0 {
			*input = (*input)[:len(*input)-1]
		}
	}
}

func valid_number_input(input []rune) string {
	var valid_input string
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			valid_input += string(input[i])
		}
	}
	return valid_input
}
