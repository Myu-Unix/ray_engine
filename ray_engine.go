/*
ray_engine by Myu-Unix - Sept 2020
Inspired by 3DSage's C/OpenGL raycasting engine : https://www.youtube.com/watch?v=gYRrGTC7GtA
Just a prototype from a non-pro :)
*/

package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"math"
	"os"
	"strings"
	//"strconv"
)

var (
	backgroundImage   *ebiten.Image
	SplashImage       *ebiten.Image
	wallImage         *ebiten.Image
	keyStates                 = map[ebiten.Key]int{}
	STATE_SHOW_DEBUG          = 1
	STATE_SHOW_2D_MAP         = 1
	STATE_FULLSCREEN          = 0
	bg_posx           float64 = 0
	bg_posy           float64 = 0
	wall_posx         float64 = 0
	wall_posy         float64 = 0
	CONST_PI          float64 = 3.1415926535
	CONST_PI2         float64 = CONST_PI / 2
	CONST_PI3         float64 = (3 * CONST_PI) / 2
	CONST_DR          float64 = 0.0174533 // one radian in degrees
	player_pos_x      float64 = 220
	player_pos_y      float64 = 320
	player_delta_x    float64 = 0
	player_delta_y    float64 = 0
	player_angle      float64 = 0 // Initial angle
	mapX              int     = 8
	mapY              int     = 8
	boot                      = 56
	engine_version            = "ray_engine 0.6.0 (Summer 2021)"
	debug_str                 = "'z/s/q/d' (Azerty) to move, 'k' to exit, 'm' to toogle the 2D map"
	debug_str2                = "'f' for sullscreen, 'i' for debug info"
	str               string
	map_array         = [64]int{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 1, 1,
		1, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
	}
	// Rays vars
	r, mx, my, mp, dof                 int
	rx, ry, ra, xo, yo, hx, hy, vx, vy float64
	aTan                               float64 = 0
	nTan                               float64 = 0
	disH                               float64 = 1000000
	disV                               float64 = 1000000
	disT                               float64 = 0
	lineH                              float64 = 0
	lineO                              float64 = 0
	x3d                                float64 = 530 // Current offset to start to draw the 3D map
	x3d_orig                           float64 = 530 // Offset toogle - 0 = 2D MAP on, 530 = 2D MAP Off - set at compile time
	ca                                 float64 = 0
	COLOR_R, COLOR_G, COLOR_B          uint8
)

//fix fisheye effect
func fix_fisheye() {
	ca = player_angle - ra
	if ca < 0 {
		ca = ca + (2 * CONST_PI)
	}
	if ca > 2*CONST_PI {
		ca = ca - (2 * CONST_PI)
	}
	disT = disT * math.Cos(ca)
}

func check_horiz() {
	dof = 0
	disH = 1000000
	hx = player_pos_x
	hy = player_pos_y
	aTan = -1 / math.Tan(ra)
	if ra > CONST_PI {
		ry = float64(((int(player_pos_y) >> 6) << 6)) - float64(0.0001)
		rx = (player_pos_y-ry)*aTan + player_pos_x
		yo = -64
		xo = -yo * aTan
	}
	if ra < CONST_PI {
		ry = float64(((int(player_pos_y) >> 6) << 6)) + float64(64)
		rx = (player_pos_y-ry)*aTan + player_pos_x
		yo = 64
		xo = -yo * aTan
	}
	// Looking left or right, will not hit horizontal lines
	if ra == 0 || ra == CONST_PI {
		rx = player_pos_x
		ry = player_pos_y
		dof = 8
	}
	for dof < 8 {
		mx = (int(rx) >> 6)
		my = (int(ry) >> 6)
		mp = my*mapX + mx
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 { // was == 1
			dof = 8 // hit wall
			hx = rx
			hy = ry
			// ax ay bx by ang
			disH = math.Sqrt((hx-player_pos_x)*(hx-player_pos_x) + (hy-player_pos_y)*(hy-player_pos_y))
		} else {
			rx = rx + xo
			ry = ry + yo
			dof++
		}
	}
}

func check_verti(screen *ebiten.Image) {
	dof = 0
	disV = 1000000
	vx = player_pos_x
	vy = player_pos_y
	nTan = -math.Tan(ra)
	if ra > CONST_PI2 && ra < CONST_PI3 {
		rx = float64(((int(player_pos_x) >> 6) << 6)) - float64(0.0001)
		ry = (player_pos_x-rx)*nTan + player_pos_y
		xo = -64
		yo = -xo * nTan
	}
	if ra < CONST_PI2 || ra > CONST_PI3 {
		rx = float64(((int(player_pos_x) >> 6) << 6)) + float64(64)
		ry = (player_pos_x-rx)*nTan + player_pos_y
		xo = 64
		yo = -xo * nTan
	}
	// staight up or down
	if ra == 0 || ra == CONST_PI {
		rx = player_pos_x
		ry = player_pos_y
		dof = 8
	}
	for dof < 8 {
		mx = (int(rx) >> 6)
		my = (int(ry) >> 6)
		mp = my*mapX + mx
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 { // was == 1
			dof = 8 // hit wall
			vx = rx
			vy = ry
			// ax ay bx by ang
			disV = math.Sqrt((vx-player_pos_x)*(vx-player_pos_x) + (vy-player_pos_y)*(vy-player_pos_y))
		} else {
			rx = rx + xo
			ry = ry + yo
			dof++
		}
	}
	if disV < disH {
		rx = vx
		ry = vy
		disT = disV
		COLOR_G = 0
		COLOR_B = 255
	} else if disH < disV {
		rx = hx
		ry = hy
		disT = disH
		COLOR_G = 0
		COLOR_B = 96
	}
	// Draw Shortest ray based on disT in Orange
	if STATE_SHOW_2D_MAP == 1 {
		ebitenutil.DrawLine(screen, player_pos_x, player_pos_y, rx, ry, color.RGBA{255, 128, 0, 255})
	}
	if ra < 0 {
		ra = ra + (2 * CONST_PI)
	}
	if ra > (2 * CONST_PI) {
		ra = ra - (2 * CONST_PI)
	}
}

func cast_rays(screen *ebiten.Image) {
	ra = player_angle - CONST_DR*36 // numbers of ray / 2
	if ra < 0 {
		ra = ra + (2 * CONST_PI)
	}
	if ra > (2 * CONST_PI) {
		ra = ra - (2 * CONST_PI)
	}

	for ray := 0; ray < 72; ray++ { // numbers of ray casted
		// Horizontal lines
		check_horiz()
		// Vertical lines + smallest line drawn on screen
		check_verti(screen)
		// Fix fisheye effect
		fix_fisheye()

		// Draw 3D lines/map
		lineH = float64(64*320) / disT
		if lineH > float64(320) {
			lineH = float64(320)
		}

		// Dim 2.0 shader - COLOR_R is handled here
		if disT >= 255 {
			COLOR_R = 255
		} else if disT <= 0 {
			COLOR_R = 0
		} else {
			COLOR_R = uint8(float64(disT))
		}

		// Line offset
		lineO = 160 - (lineH / float64(3)) // was 2 (int) but 2.x helps with perspective somehow

		// x, y, ray x8, lineH
		ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(8), lineH+lineO, color.RGBA{COLOR_R, COLOR_G, COLOR_B, 255})
		x3d = x3d + 8
		ra = (ra + CONST_DR) // increment
	}
	x3d = x3d_orig
}

func update(screen *ebiten.Image) error {

	opBackground := &ebiten.DrawImageOptions{}
	opSplash := &ebiten.DrawImageOptions{}
	opwall := &ebiten.DrawImageOptions{}

	// Images options
	opBackground.GeoM.Translate(bg_posx, bg_posy)
	opwall.GeoM.Translate(wall_posx, wall_posy)

	// Draw images
	if boot > 0 {
		boot = boot - 1
		screen.DrawImage(SplashImage, opSplash)
		if boot == 1 {
			// Initial values for PDX/PDY, only applied once
			player_delta_x = math.Cos(player_angle) * 5
			player_delta_y = math.Sin(player_angle) * 5
		}
	} else {
		screen.DrawImage(backgroundImage, opBackground)
		// Draw 2D map
		if STATE_SHOW_2D_MAP == 1 {
			for i := 0; i < mapX; i++ {
				for j := 0; j < mapY; j++ {
					if map_array[i*mapX+j] == 1 {
						opwall := &ebiten.DrawImageOptions{}
						opwall.GeoM.Translate(wall_posx, wall_posy)
						screen.DrawImage(wallImage, opwall)
					}
					if wall_posx < 448 {
						wall_posx += 64
					} else {
						wall_posx = 0
						wall_posy += 64
					}
				}
			}
			wall_posx = 0
			wall_posy = 0

			// Draw player as a rect
			ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 12, 12, color.Black)
			//ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 4, 4, color.RGBA{255, 100, 100, 255})
		}
		// Raycasting
		cast_rays(screen)

		// Draw player line of sight after raycast -> https://pkg.go.dev/github.com/hajimehoshi/ebiten/ebitenutil#DrawLine
		if STATE_SHOW_2D_MAP == 1 {
			ebitenutil.DrawLine(screen, player_pos_x, player_pos_y, player_pos_x+(player_delta_x*5), player_pos_y+(player_delta_y*5), color.RGBA{255, 255, 0, 255})
		}

		// Show debug info
		str = `{{.newline}} {{.engine_version}} {{.newline}} {{.debug_str}} {{.newline}} {{.debug_str2}} {{.newline}} {{.player_x}} {{.newline}} {{.player_y}}`
		str = strings.Replace(str, "{{.engine_version}}", engine_version, -1)
		str = strings.Replace(str, "{{.debug_str}}", debug_str, -1)
		str = strings.Replace(str, "{{.debug_str2}}", debug_str2, -1)
		str = strings.Replace(str, "{{.newline}}", "\n                                                                                                ", -1)
		str = strings.Replace(str, "{{.player_x}}", fmt.Sprintf("Player X : %f", player_pos_x), -1)
		str = strings.Replace(str, "{{.player_y}}", fmt.Sprintf("Player Y : %f", player_pos_y), -1)
		if STATE_SHOW_DEBUG == 1 {
			ebitenutil.DebugPrint(screen, str)
		}
	}

	// Keyboard handling
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		keyStates[ebiten.KeyM]++
	} else {
		keyStates[ebiten.KeyM] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		keyStates[ebiten.KeyF]++
	} else {
		keyStates[ebiten.KeyF] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		keyStates[ebiten.KeyK]++
	} else {
		keyStates[ebiten.KeyK] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		keyStates[ebiten.KeyI]++
	} else {
		keyStates[ebiten.KeyI] = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) { // Z Azerty
		player_pos_x = player_pos_x + player_delta_x
		player_pos_y = player_pos_y + player_delta_y
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) { // S
		player_pos_x = player_pos_x - player_delta_x
		player_pos_y = player_pos_y - player_delta_y
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) { // Q Azerty
		player_angle -= 0.05
		// Reset
		if player_angle <= 0 {
			player_angle = 6.283
		}
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) { // D
		player_angle += 0.05
		// Reset
		if player_angle >= 6.283 {
			player_angle = 0
		}
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
	}
	if IsKeyTriggered(ebiten.KeyM) == true {
		if STATE_SHOW_2D_MAP == 0 {
			STATE_SHOW_2D_MAP = 1
			x3d_orig = 530
		} else {
			STATE_SHOW_2D_MAP = 0
			x3d_orig = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyF) == true {
		if STATE_FULLSCREEN == 0 {
			ebiten.SetFullscreen(true)
			STATE_FULLSCREEN = 1
		} else {
			ebiten.SetFullscreen(false)
			STATE_FULLSCREEN = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyI) == true {
		if STATE_SHOW_DEBUG == 0 {
			STATE_SHOW_DEBUG = 1
		} else {
			STATE_SHOW_DEBUG = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyK) == true {
		os.Exit(0)
	}

	return nil
}

func main() {
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("bg2.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	SplashImage, _, err = ebitenutil.NewImageFromFile("splash.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	wallImage, _, err = ebitenutil.NewImageFromFile("wall.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(engine_version)
	ebiten.Run(update, 1024, 512, 1, engine_version)
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}
