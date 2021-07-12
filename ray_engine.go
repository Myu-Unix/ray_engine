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
)

var (
	backgroundImage *ebiten.Image
	SplashImage     *ebiten.Image
	wallImage       *ebiten.Image
	keyStates               = map[ebiten.Key]int{}
	show_debug              = 1
	bg_posx         float64 = 0
	bg_posy         float64 = 0
	wall_posx       float64 = 0
	wall_posy       float64 = 0
	CONST_PI        float64 = 3.14159
	CONST_PI2       float64 = CONST_PI / 2
	CONST_PI3       float64 = 3 * (CONST_PI / 2)
	CONST_DR        float64 = 0.0174533 // one radian in degrees
	player_pos_x    float64 = 220
	player_pos_y    float64 = 320
	player_delta_x  float64 = 0
	player_delta_y  float64 = 0
	player_angle            = 4.33 // Initial angle
	mapX            int     = 8
	mapY            int     = 8
	boot                    = 56
	mouse_enabled	int	= 0
	engine_version          = "ray_engine 0.5.7.1"
	debug_str               = "'z/s/q/d' (Azerty) to move, 'k' to exit"
	str             string
	map_array = [64]int{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 1, 0, 0, 0, 0, 1,
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
	x3d float64	= 530 // Offset to start to draw the 3D map
	ca float64 = 0
	COLOR_R, COLOR_G, COLOR_B uint8
)

//fix fisheye effect
func fix_fisheye() {
  ca = player_angle - ra
  if ca < 0 {
   ca = ca + (2 * CONST_PI)
  }
  if ca > 2 * CONST_PI {
    ca = ca - (2* CONST_PI)
  }
  disT = disT * math.Cos(ca)
}

func cast_horiz() {
			dof = 0
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
		if ra == 0 || ra == CONST_PI {
			rx = player_pos_x
			ry = player_pos_y
			dof = 8
		}
		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*mapX + mx
			if (mp > 0 && mp < mapX*mapY) && map_array[mp] == 1 {
				dof = 8 // hit wall
				hx = rx
				hy = ry
				// ax ay bx by ang
				disH = math.Sqrt((hx-player_pos_x)*(hx-player_pos_x) + (hy-player_pos_y)*(hy-player_pos_y))
			} else {
				rx = rx + xo
				ry = ry + yo
				dof = dof + 1
			}
		}
}

func cast_verti(screen *ebiten.Image) {
dof = 0
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
		if ra == 0 || ra == CONST_PI {
			rx = player_pos_x
			ry = player_pos_y
			dof = 8
		}
		for dof < 8 {
			mx = (int(rx) >> 6)
			my = (int(ry) >> 6)
			mp = my*mapX + mx
			if (mp > 0 && mp < mapX*mapY) && map_array[mp] == 1 {
				dof = 8 // hit wall
				vx = rx ; vy = ry
				// ax ay bx by ang
				disV = math.Sqrt((vx-player_pos_x)*(vx-player_pos_x) + (vy-player_pos_y)*(vy-player_pos_y))
			} else {
				rx = rx + xo ; ry = ry + yo ; dof = dof + 1
			}
		}
		if disV < disH {
			rx = vx ; ry = vy ; disT = disV
			COLOR_R = 192
			COLOR_G = 0
			COLOR_B = 0
		}
		if disH < disV {
			rx = hx ; ry = hy ; disT = disH
			COLOR_R = 249
			COLOR_G = 3
			COLOR_B = 0
		}
		// Draw Shortest ray based on disT in Orange
		ebitenutil.DrawLine(screen, player_pos_x, player_pos_y, rx, ry, color.RGBA{255, 128, 0, 255})
		ra = (ra + CONST_DR)
		if ra < 0 {
			ra = ra + (2 * CONST_PI)
		}
		if ra > (2*CONST_PI) {
			ra = ra - (2 * CONST_PI)
		}
}

func cast_rays(screen *ebiten.Image) {
	ra = player_angle - CONST_DR*30
	if ra < 0 {
		ra = ra + 2*CONST_PI
	}
	if ra > 2*CONST_PI {
		ra = ra - 2*CONST_PI
	}
	for ray := 0; ray < 64; ray++ { // numbers of ray casted
		// Horizontal lines
		cast_horiz()
		// Vertical lines + smallest line
		cast_verti(screen)
		// Draw 3D lines/map
		fix_fisheye()

		lineH = float64(64 * 320) / disT
		if lineH > float64(320) {
			lineH = float64(320)
		}
		// "dim" far objects
		if lineH < float64(120) {
			COLOR_R = 97
			COLOR_G = 0
			COLOR_B = 0
		}
		lineO = 160 - (lineH/float64(2.3)) // was 2 (int) but 2.x helps with perspective somehow
		// x, y, rayx8, lineH
		ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(8), lineH+lineO, color.RGBA{COLOR_R, COLOR_G, COLOR_B, 255})

		x3d = x3d + 8
	}
	x3d = 530
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
		// Initial values for PDX/PDY
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
	} else {
		screen.DrawImage(backgroundImage, opBackground)
		// Draw map
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
		ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 8, 8, color.Black)
		ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 4, 4, color.RGBA{255, 100, 100, 255})

		// Raycasting
		cast_rays(screen)

		// Draw player line of sight after raycast
		// https://pkg.go.dev/github.com/hajimehoshi/ebiten/ebitenutil#DrawLine
		ebitenutil.DrawLine(screen, player_pos_x, player_pos_y, player_pos_x+(player_delta_x*5), player_pos_y+(player_delta_y*5), color.RGBA{255, 255, 0, 255})

		// Show debug info
		str = `{{.newline}} {{.engine_version}} {{.newline}} {{.debug_str}} {{.newline}} {{.player_angle}}`
		str = strings.Replace(str, "{{.engine_version}}", engine_version, -1)
		str = strings.Replace(str, "{{.debug_str}}", debug_str, -1)
		str = strings.Replace(str, "{{.newline}}", "\n                                                                                                                                ", -1)
		str = strings.Replace(str, "{{.player_angle}}", fmt.Sprintf("Player angle : %f", player_angle), -1)
		if show_debug == 1 {
			ebitenutil.DebugPrint(screen, str)
		}

	}

	// --- Dodgy AF Mouse support, disabled by default
	if mouse_enabled == 1 {
	x, _ := ebiten.CursorPosition()
	fmt.Println("Mouse X : %f", x)
	player_angle = float64(x)/48 // Influence mouse sensitivity
	// reset
	if player_angle < 0 {
			player_angle = 6.283
		}
	if player_angle > 6.283 {
			player_angle = 0
		}
	player_delta_x = math.Cos(player_angle) * 5
	player_delta_y = math.Sin(player_angle) * 5
	// ---
	}

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
		if player_angle < 0 {
			player_angle = 6.283
		}
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) { // D
		player_angle += 0.05
		// Reset
		if player_angle > 6.283 {
			player_angle = 0
		}
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
	}
	if IsKeyTriggered(ebiten.KeyM) == true {
		if show_debug == 0 {
          show_debug = 1
	    } else {
	      show_debug = 0
	    }
	}
	if IsKeyTriggered(ebiten.KeyF) == true {
		// Bug : cannot go back to normal size afterwards
        ebiten.SetFullscreen(true)
    }
	if IsKeyTriggered(ebiten.KeyK) == true {
		os.Exit(0)
	}

	return nil
}

func main() {
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("bg.png", ebiten.FilterNearest)
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
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	ebiten.Run(update, 1024, 512, 1, engine_version)
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}
