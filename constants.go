package main

import (
	"image/color"
)

const SCREEN_WIDTH = 640
const SCREEN_HEIGHT = 480

const BLDG_WINDOW_WIDTH = 6
const BLDG_WINDOW_HEIGHT = 8

const MAX_WINDOWS = int(SCREEN_WIDTH / (BLDG_WINDOW_WIDTH * 2))

const GRAVITY = 2.2
const SCALE = 0.2

const PROJECTILE_WIDTH = 16
const PROJECTILE_HEIGHT = 16

const PLAYER_WIDTH = 32
const PLAYER_HEIGHT = 32

var BACKGROUND_COLOR = color.RGBA{0x1B, 0x00, 0x2F, 0xFF}
