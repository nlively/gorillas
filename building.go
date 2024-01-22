package main

import (
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/vector"
  "image/color"
  _ "image/png"
)

type building struct {
	color color.Color
	windows int
	floors int
	x int
	y int
	width int
	height int
}

func (b *building) building_width() int {
	return b.windows * (window_width * 2)
}
  
func (b *building) building_height() int {
	return b.floors * (window_height * 2)
}

func (b *building) set_coordinates(x int, y int, width int, height int) {
	b.x = x
	b.y = y
	b.width = width
	b.height = height
}

func (b *building) draw_building(screen *ebiten.Image) {
  vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), b.color, true)

  for j := 0; j < b.windows; j++ {
    for k := 0; k < b.floors; k++ {
      win_x := (b.x + (12*j)) + (window_width / 2);
      win_y := b.y + (16*k) + (window_height / 2)
      vector.DrawFilledRect(screen, float32(win_x), float32(win_y), float32(window_width), float32(window_height), color.RGBA{0xff, 0xff, 0xff, 0x99}, true)
    }
  }
}