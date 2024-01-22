package main

// NOTE: this file can probably be deleted.  input/output in text mode

import (
	"fmt"
)

func describe_wind_direction(g *game) string {
  switch g.wind_direction {
    case true:
      return "east"
    default:
      return "west"
  }
}

func capture_angle(angle *int) bool {
  fmt.Scanln(angle)

  if *angle < 0 || *angle > 90 {
    fmt.Println("Angle must be between 0 and 90. Please try again.")
    return false
  } else {
    return true
  }
}

func capture_velocity(velocity *int) bool {
  fmt.Scanln(velocity)

  if *velocity < 1 || *velocity > 200 {
    fmt.Println("Velocity must be between 1 and 200. Please try again.")
    return false
  } else {
    return true
  }
}

func turn(p *player) {
  var angle, velocity int
  fmt.Printf("%s's turn!\n", p.name)

  fmt.Println("Enter an angle between 0 and 90: ")
  for {
    if capture_angle(&angle) {
      break
    }
  }

  fmt.Println("Enter a velocity between 1 and 200: ")
  for {
    if capture_velocity(&velocity) {
      break
    }
  }


  fmt.Printf("Firing at an angle of %d and a velocity of %d...\n",angle, velocity)
}

func capture_player_name(number int) string {
  var name_input string
  fmt.Printf("Enter player %d name: ", number)
  fmt.Scanln(&name_input)
  return name_input
}