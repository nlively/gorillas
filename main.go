package main

import (
  "fmt"
  "math/rand"
  "log"
  "github.com/hajimehoshi/ebiten/v2"
  // "github.com/hajimehoshi/ebiten/v2/ebitenutil"
  // "github.com/hajimehoshi/ebiten/v2/vector"
  "github.com/hajimehoshi/ebiten/v2/inpututil"
  // "github.com/hajimehoshi/ebiten/v2/text"
  // "image"
  "image/color"
  _ "image/png"
)



func (g *game) Update() error {
  g.keys = inpututil.AppendPressedKeys(g.keys[:0])
  if len(g.keys) > 0 && g.keys[0].String() == "F" && !g.firing {
    g.fire()
    fmt.Println("Fire!")
  } else if g.firing {
    g.projectile.x += 5
    if g.projectile.x >= screen_width {
      g.stop_projectile()
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

func (g *game) Draw(screen *ebiten.Image) {
  if !g.setup {
    g.setup_grid(screen)
  }

  for i:= 0; i < len(g.buildings); i++ {
    g.buildings[i].draw_building(screen)
  }

  g.player1.draw_image(screen)
  g.player2.draw_image(screen)

  if g.firing {
    g.draw_projectile(screen)
  }
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return screen_width, screen_height
}

func main() {
  // new_game()
  // summarize()
  
  g := game{}
  g.setup_game()

  // var counter = 0
  // var current_player = 0
  // for {
  //   counter++
  //   p := board.players[current_player]
  //   turn(&p)
  //   if winner || counter > 3 { 
  //     break
  //   }
  //   current_player = current_player ^ 1
  // }

  ebiten.SetWindowSize(screen_width, screen_height)
  ebiten.SetWindowTitle("Tori vs Evan")
  if err := ebiten.RunGame(&g); err != nil {
    log.Fatal(err)
  }
  
}
