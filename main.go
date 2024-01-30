package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
)

func (g *Game) Update() error {
	switch g.state {
	case INTRO:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.switch_state(MAIN_GAME)
		} else {
			g.intro.blink()
		}
	case MAIN_GAME:
		// Capture turn input
		switch g.turn_state {
		case COLLECT_ANGLE:
			g.angle_runes = ebiten.AppendInputChars(g.angle_runes[:0])
			g.angle_input += valid_number_input(g.angle_runes)
			handle_backspace(&g.angle_input)
		case COLLECT_VELOCITY:
			g.velocity_runes = ebiten.AppendInputChars(g.velocity_runes[:0])
			g.velocity_input += valid_number_input(g.velocity_runes)
			handle_backspace(&g.velocity_input)
		}

		// Handle completed turn input
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.firing {
			if g.turn_state == COLLECT_ANGLE {
				if g.validate_angle() {
					g.turn_state = COLLECT_VELOCITY
				}
			} else if g.turn_state == COLLECT_VELOCITY {
				if g.validate_velocity() {
					g.fire_projectile()
				}
			}
		}

		if g.firing {
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

func (g *Game) draw_main_game(screen *ebiten.Image) {
	if !g.setup {
		g.setup_grid(screen)
	}

	for i := 0; i < len(g.buildings); i++ {
		g.buildings[i].draw_building(screen)
	}

	g.player1.draw_image(screen)
	g.player2.draw_image(screen)

	for i := 0; i < len(g.holes); i++ {
		g.holes[i].draw_hole(screen)
	}

	g.draw_wind(screen)

	if g.firing {
		g.draw_projectile(screen)
	} else {
		g.draw_turn_inputs(screen)
	}

	if g.error_text != "" {
		draw_error_text(screen, g.error_text, 240, 200)
	}
}

func (g *Game) draw_victory(screen *ebiten.Image) {
	var message string
	message = fmt.Sprintf("Winner: %s", g.winner.name)
	draw_text(screen, message, 240, 200)

	draw_text(screen, "Press Enter to Start Over", 190, 360)
}

func (g *Game) draw_background(screen *ebiten.Image) {
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
