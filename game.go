package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"strconv"
)

type Game struct {
	buildings      []building
	holes          []hole
	player1        player
	player2        player
	wind           float64
	winner         *player
	keys           []ebiten.Key
	firing         bool
	projectile     projectile
	angle          float64
	velocity       float64
	setup          bool
	state          int
	turn_player    *player
	turn_state     int
	intro          Intro
	angle_runes    []rune
	velocity_runes []rune
	angle_input    string
	velocity_input string
	error_text     string
}

type GameState struct {
	name string
}

const INTRO = 1
const MAIN_GAME = 2
const VICTORY = 3

const COLLECT_ANGLE = 10
const COLLECT_VELOCITY = 11

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

func (g *Game) next_turn() {
	if g.turn_player == &g.player1 {
		g.turn_player = &g.player2
	} else {
		g.turn_player = &g.player1
	}
	g.turn_state = COLLECT_ANGLE
	g.angle_input = ""
	g.velocity_input = ""
}

func (g *Game) setup_game() {
	g.wind = rand.Float64() - 0.5

	g.setup_players()
	g.setup_buildings()

	g.turn_player = &g.player1
	g.turn_state = COLLECT_ANGLE
}

func (g *Game) start_over() {
	g.setup = false
	g.setup_game()
	g.holes = make([]hole, 0)
	g.firing = false
	g.turn_state = COLLECT_ANGLE
	g.angle_input = ""
	g.velocity_input = ""
	g.error_text = ""
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

func (g *Game) fire_projectile() {
	g.firing = true
	g.projectile.x = g.turn_player.x
	g.projectile.y = g.turn_player.y - 32

	fmt.Printf("angle %f, velocity %f, wind %f\n", g.angle, g.velocity, g.wind)

	radian := float64(g.angle) * math.Pi / 180

	// Calculate increment values
	g.projectile.dx = float64(g.velocity)*math.Cos(radian) + float64(g.wind)
	g.projectile.dy = float64(g.velocity) * math.Sin(radian) * -1

	if g.turn_player == &g.player2 {
		g.projectile.dx *= -1
	}

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

func (g *Game) add_hole(projectile *projectile) {
	px := projectile.x + (PROJECTILE_WIDTH / 2)
	py := projectile.y + (PROJECTILE_HEIGHT / 2)
	h := hole{int(px), int(py)}
	g.holes = append(g.holes, h)
}

func (g *Game) stop_projectile() {
	g.firing = false
	g.next_turn()
}

func (g *Game) draw_wind(screen *ebiten.Image) {
	var wind_text string
	wind_human := int(math.Abs(g.wind) * 300)
	if g.wind > 0 {
		wind_text = fmt.Sprintf("Wind: %dmph east", wind_human)
	} else {
		wind_text = fmt.Sprintf("Wind: %dmph west", wind_human)
	}

	draw_any_text(screen, wind_text, 475, 30, 16, color.RGBA{0xFF, 0xFF, 0x00, 0xFF})
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

func (g *Game) draw_turn_inputs(screen *ebiten.Image) {
	var angle_x, velocity_x int
	var angle_input_text, velocity_input_text string
	angle_y := 100
	velocity_y := 125

	if g.turn_player == &g.player2 {
		angle_x = 450
		velocity_x = 450
	} else {
		angle_x = 10
		velocity_x = 10
	}

	draw_text(screen, fmt.Sprintf("%s's Turn", g.turn_player.name), 10, 40)

	angle_input_text = fmt.Sprintf("Enter Angle: %s", g.angle_input)

	if g.turn_state == COLLECT_VELOCITY {
		velocity_input_text = fmt.Sprintf("Enter Velocity: %s", g.velocity_input)
	}

	draw_smaller_text(screen, angle_input_text, angle_x, angle_y)
	draw_smaller_text(screen, velocity_input_text, velocity_x, velocity_y)
}

func (g *Game) validate_angle() bool {
	angle, err := strconv.ParseFloat(g.angle_input, 64)
	if err != nil {
		return false
	}

	if angle < 0 || angle > 90 {
		g.error_text = "Angle must be between 0 and 90"
		return false
	}

	g.angle = angle
	g.error_text = ""
	return true
}

func (g *Game) validate_velocity() bool {
	velocity, err := strconv.ParseFloat(g.velocity_input, 64)
	if err != nil {
		return false
	}

	if velocity < 25 || velocity > 200 {
		g.error_text = "Velocity must be between 25 and 200"
		return false
	}

	g.velocity = velocity
	g.error_text = ""
	return true
}
