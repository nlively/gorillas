package main

// import dependencies
import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type hole struct {
	x int
	y int
}

// function to draw a hole
func (h *hole) draw_hole(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(h.x), float32(h.y), 16, color.Black, true)
}