package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	// "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/hajimehoshi/ebiten/v2/text"
	// "image"
	"image/color"
	_ "image/png"
)

func (g *Game) Update() error {
	// g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	switch g.state {
	case INTRO:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.switch_state(MAIN_GAME)
		} else {
			g.intro.blink()
		}
	case MAIN_GAME:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !g.firing {
			g.fire()
			fmt.Println("Fire!")
		} else if g.firing {
			g.move_projectile()
		}
	case VICTORY:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.start_over()
			g.switch_state(MAIN_GAME)
		}
	}

	return nil
}

func random_color() color.Color {
	red := uint8(rand.Intn(255))
	green := uint8(rand.Intn(255))
	blue := uint8(rand.Intn(255))
	return color.RGBA{red, green, blue, 0xFF}
}

func (g *Game) draw_main_game(screen *ebiten.Image) {
	if !g.setup {
		g.setup_grid(screen)
	}

	for i := 0; i < len(g.buildings); i++ {
		g.buildings[i].draw_building(screen)
	}

	g.player1.draw_image(screen)
	g.player2.draw_image(screen)

	// draw holes as black circles
	for i := 0; i < len(g.holes); i++ {
		g.holes[i].draw_hole(screen)
	}

	if g.firing {
		g.draw_projectile(screen)
	}
}

func (g *Game) draw_victory(screen *ebiten.Image) {
	var message string
	message = fmt.Sprintf("Winner: %s", g.winner.name)
	draw_text(screen, message, 240, 200)

	draw_text(screen, "Press Enter to Start Over", 190, 360)
}

// function to fill the background with a navy blue color
func (g *Game) draw_background(screen *ebiten.Image) {
	// screen.Fill(color.Navy)
	screen.Fill(BACKGROUND_COLOR)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.draw_background(screen)
	switch g.state {
	case INTRO:
		g.intro.draw_intro(screen)
	case MAIN_GAME:
		g.draw_main_game(screen)
	case VICTORY:
		g.draw_victory(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	game := &Game{state: INTRO}
	game.setup_game()

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Tori vs Evan")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
