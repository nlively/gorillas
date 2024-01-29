package main

type projectile struct {
	x  float64
	y  float64
	dx float64
	dy float64

	//px float64
	//py float64
}

// function to detect collision between projectile and player
func (p *projectile) detect_collision(player *player) bool {
	px := float64(p.x)
	py := float64(p.y)
	if px < player.x+PLAYER_WIDTH && px+PROJECTILE_WIDTH > player.x && py < player.y+PLAYER_HEIGHT && py+PROJECTILE_HEIGHT > player.y {
		return true
	}
	return false
}
