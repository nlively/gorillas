package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	// "github.com/hajimehoshi/ebiten/v2/vector"
	// "github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/hajimehoshi/ebiten/v2/text"
	// "image"
	// "image/color"
	_ "image/png"
)

type game struct {
	buildings []building
	//holes      []hole
	player1 player
	player2 player
	wind    float64
	//winner     bool
	keys       []ebiten.Key
	firing     bool
	projectile projectile
	angle      float64
	velocity   float64
	setup      bool
}

func (g *game) setup_players() {
	// g.board.players[0] = player{name: capture_player_name(1)}
	// g.board.players[1] = player{name: capture_player_name(2)}
	g.player1 = player{name: "Tori", path: "images/tori.png"}
	g.player2 = player{name: "Evan", path: "images/evan.png"}
}

func (g *game) setup_buildings() {
	g.buildings = make([]building, 0)
	i := 0
	total_windows := 0

	// Randomly generate enough buildings to fill the screen
	for {
		var windows int
		if total_windows >= max_windows-5 {
			windows = max_windows - total_windows
		} else {
			windows = rand.Intn(3) + 3
		}
		total_windows += windows

		b := building{windows: windows, floors: rand.Intn(8) + 2, color: random_color()}

		g.buildings = append(g.buildings, b)

		if total_windows >= max_windows {
			break
		}
		i++
	}
}

func (g *game) setup_game() {
	g.wind = rand.Float64() - 0.5

	g.setup_players()
	g.setup_buildings()
}

func (g *game) setup_grid(screen *ebiten.Image) {
	g.setup = true
	running_width := 0

	for i := 0; i < len(g.buildings); i++ {
		building := &g.buildings[i]

		width := building.building_width()
		height := building.building_height()
		x := running_width
		y := screen.Bounds().Dy() - height

		running_width = running_width + width

		building.set_coordinates(float32(x), float32(y), float32(width), float32(height))

		if i == 1 {
			// player 1 sits on top of 2nd building
			img_x := int(x + width/2)
			g.player1.set_coordinates(float64(img_x), float64(y))
		} else if i == len(g.buildings)-2 {
			// player 2 sits on top of 2nd to last building
			img_x := int(x + width/2)
			g.player2.set_coordinates(float64(img_x), float64(y))
		}
	}
}

func (g *game) fire() {
	g.angle = float64(rand.Intn(90))
	g.velocity = 25 + (rand.Float64() * 175 * .7)
	g.firing = true
	g.projectile.x = g.player1.x
	g.projectile.y = g.player1.y

	fmt.Printf("angle %f, velocity %f, wind %f\n", g.angle, g.velocity, g.wind)

	radian := float64(g.angle) * math.Pi / 180

	// Calculate increment values
	g.projectile.dx = float64(g.velocity)*math.Cos(radian) + float64(g.wind)
	g.projectile.dy = float64(g.velocity) * math.Sin(radian) * -1

	g.projectile.dx *= scale
	g.projectile.dy *= scale

	fmt.Printf("dx,dy = %f, %f\n", g.projectile.dx, g.projectile.dy)
}

func (g *game) move_projectile() {

	// Increment coords by a set increment value
	g.projectile.dx -= float64(g.wind) * scale
	g.projectile.dy += float64(gravity) * scale

	g.projectile.x += g.projectile.dx
	g.projectile.y += g.projectile.dy

	// g.projectile.px += g.projectile.dx * scale;
	// g.projectile.py += g.projectile.dy * scale;

	// adjusted_x := float64(g.projectile.x) + g.projectile.px
	// adjusted_y := float64(g.projectile.y) - g.projectile.py

	// fmt.Printf("projectile coords %f, %f\n", g.projectile.x, g.projectile.y)

	// if adjusted_x > screen_width || adjusted_y > screen_height {
	// 	g.stop_projectile()
	// }
	if g.projectile.x > screen_width || g.projectile.y > screen_height {
		g.stop_projectile()
	}
}

func (g *game) stop_projectile() {
	g.firing = false
}

/*
func (g *game) summarize() {
	// fmt.Println("-----")
	// fmt.Println("Player 1 is", g.player1.name)
	// fmt.Println("Player 2 is", g.player2.name)
	// fmt.Println("Wind direction is", describe_wind_direction(g))
	// fmt.Printf("Wind speed is %dmph\n", g.wind_speed)

	// for i := 0; i < len(g.buildings); i++ {
	//   fmt.Printf("Building %d has %d windows and %d floors\n", i+1, g.buildings[i].windows, g.buildings[i].floors)
	// }
}
*/

func (g *game) draw_projectile(screen *ebiten.Image) {
	path := "images/pizza.png"
	var err error
	var img *ebiten.Image

	img, _, err = ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// new_x := float64(g.projectile.x) + g.projectile.px - float64(img.Bounds().Dx() / 2)
	// new_y := float64(g.projectile.y) - g.projectile.py - float64(img.Bounds().Dy() / 2)

	// Center the image on the current point
	new_x := float64(g.projectile.x) + float64(img.Bounds().Dx()/2)
	new_y := float64(g.projectile.y) - float64(img.Bounds().Dy()/2)

	// fmt.Printf("drawing at %f, %f\n", new_x, new_y)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(new_x, new_y)

	screen.DrawImage(img, op)
}
