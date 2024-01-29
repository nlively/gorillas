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

type Game struct {
	buildings  []building
	holes      []hole
	player1    player
	player2    player
	wind       float64
	winner     *player
	keys       []ebiten.Key
	firing     bool
	projectile projectile
	angle      float64
	velocity   float64
	setup      bool
	state      int
	intro      Intro
}

type GameState struct {
	name string
}

const INTRO = 1
const MAIN_GAME = 2
const VICTORY = 3

func (g *Game) switch_state(new_state int) {
	g.state = new_state
}

func (g *Game) trigger_win(player *player) {
	g.winner = player
	g.switch_state(VICTORY)
}

func (g *Game) setup_players() {
	g.player1 = player{name: "Tori", path: "images/tori.png"}
	g.player2 = player{name: "Evan", path: "images/evan.png"}
}

func (g *Game) setup_buildings() {
	g.buildings = make([]building, 0)
	i := 0
	total_windows := 0

	// Randomly generate enough buildings to fill the screen
	for {
		var windows int
		if total_windows >= MAX_WINDOWS-5 {
			windows = MAX_WINDOWS - total_windows
		} else {
			windows = rand.Intn(3) + 3
		}
		total_windows += windows

		b := building{windows: windows, floors: rand.Intn(8) + 2, color: random_color()}

		g.buildings = append(g.buildings, b)

		if total_windows >= MAX_WINDOWS {
			break
		}
		i++
	}
}

func (g *Game) setup_game() {
	g.wind = rand.Float64() - 0.5

	g.setup_players()
	g.setup_buildings()
}

// function to start over
func (g *Game) start_over() {
	g.setup = false
	g.setup_game()
	g.holes = make([]hole, 0)
	g.firing = false
}

func (g *Game) setup_grid(screen *ebiten.Image) {
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

func (g *Game) fire() {
	g.angle = float64(rand.Intn(90))
	g.velocity = 25 + (rand.Float64() * 175 * .7)
	g.firing = true
	g.projectile.x = g.player1.x
	g.projectile.y = g.player1.y - 32

	fmt.Printf("angle %f, velocity %f, wind %f\n", g.angle, g.velocity, g.wind)

	radian := float64(g.angle) * math.Pi / 180

	// Calculate increment values
	g.projectile.dx = float64(g.velocity)*math.Cos(radian) + float64(g.wind)
	g.projectile.dy = float64(g.velocity) * math.Sin(radian) * -1

	g.projectile.dx *= SCALE
	g.projectile.dy *= SCALE

	fmt.Printf("dx,dy = %f, %f\n", g.projectile.dx, g.projectile.dy)
}

func (g *Game) move_projectile() {

	// Increment coords by a set increment value
	g.projectile.dx -= float64(g.wind) * SCALE
	g.projectile.dy += float64(GRAVITY) * SCALE

	g.projectile.x += g.projectile.dx
	g.projectile.y += g.projectile.dy

	// detect a collision between the projectile and a player
	if g.projectile.detect_collision(&g.player1) {
		g.stop_projectile()
		g.trigger_win(&g.player2)
		return
	} else if g.projectile.detect_collision(&g.player2) {
		g.stop_projectile()
		g.trigger_win(&g.player1)
		return
	}

	// detect a collision between the projectile and a building
	for i := 0; i < len(g.buildings); i++ {
		building := &g.buildings[i]
		if building.detect_collision(&g.projectile) {
			g.stop_projectile()
			// add hole to game
			g.add_hole(&g.projectile)
			break
		}
	}

	if g.projectile.x > SCREEN_WIDTH || g.projectile.y > SCREEN_HEIGHT {
		g.stop_projectile()
	}
}

// function to add a hole to the game
func (g *Game) add_hole(projectile *projectile) {
	px := projectile.x + (PROJECTILE_WIDTH / 2)
	py := projectile.y + (PROJECTILE_HEIGHT / 2)
	h := hole{int(px), int(py)}
	g.holes = append(g.holes, h)
}

func (g *Game) stop_projectile() {
	g.firing = false
}

func (g *Game) draw_projectile(screen *ebiten.Image) {
	path := "images/pizza.png"
	var err error
	var img *ebiten.Image

	img, _, err = ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Center the image on the current point
	new_x := float64(g.projectile.x) + float64(img.Bounds().Dx()/2)
	new_y := float64(g.projectile.y) - float64(img.Bounds().Dy()/2)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(new_x, new_y)

	screen.DrawImage(img, op)
}
