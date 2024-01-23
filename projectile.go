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
	if px < player.x+player_width && px+projectile_width > player.x && py < player.y+player_height && py+projectile_height > player.y {
		return true
	}
	return false
} 	
