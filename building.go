package main

import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type building struct {
	color   color.Color
	windows int
	floors  int
	x       float32
	y       float32
	width   float32
	height  float32
}

func (b *building) building_width() int {
	return b.windows * (BLDG_WINDOW_WIDTH * 2)
}

func (b *building) building_height() int {
	return b.floors * (BLDG_WINDOW_HEIGHT * 2)
}

func (b *building) set_coordinates(x float32, y float32, width float32, height float32) {
	b.x = x
	b.y = y
	b.width = width
	b.height = height
}

func (b *building) draw_building(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, b.x, b.y, b.width, b.height, b.color, true)

	for j := 0; j < b.windows; j++ {
		for k := 0; k < b.floors; k++ {
			win_x := (b.x + float32(12*j)) + (BLDG_WINDOW_WIDTH / 2)
			win_y := b.y + float32(16*k) + (BLDG_WINDOW_HEIGHT / 2)
			vector.DrawFilledRect(screen, win_x, win_y, BLDG_WINDOW_WIDTH, BLDG_WINDOW_HEIGHT, color.RGBA{0xff, 0xff, 0xff, 0x99}, true)
		}
	}
}

// detect collision between projectile and building
func (b *building) detect_collision(projectile *projectile) bool {
	px := float32(projectile.x) + (PROJECTILE_WIDTH / 2)
	py := float32(projectile.y) + (PROJECTILE_HEIGHT / 2)
	if px > b.x && px < b.x+b.width && py > b.y && py < b.y+b.height {
		return true
	}
	return false
}
