/*
ray_engine by Myu-Unix - Sept 2020
Inspired by 3DSage's C/OpenGL raycasting engine : https://www.youtube.com/watch?v=gYRrGTC7GtA
*/

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"image/color"
	_ "image/png"
	"log"
	"math"
	"os/exec"
	"strings"
)

type game struct{}

var (
	STATE_SHOW_DEBUG    = 1
	STATE_SHOW_2D_MAP   = 0
	STATE_FULLSCREEN    = 0
	STATE_MOUSE_SUPPORT = 1
	STATE_YSHEARING     = 0
	STATE_COLLISION     = 0
	STATE_SOUND         = 0
	SHOW_MINIMAP        = 1
	backgroundImage     *ebiten.Image
	splashImage         *ebiten.Image
	wallImage           *ebiten.Image
	wallMiniImage       *ebiten.Image
	gunImage            *ebiten.Image
	crossHairImage      *ebiten.Image
	fireImage           *ebiten.Image
	keyStates                   = map[ebiten.Key]int{}
	bg_posx             float64 = 0
	bg_posy             float64 = 0
	wall_posx           float64 = 0
	wall_posy           float64 = 0
	wall_minimap_posx   float64 = 20
	wall_minimap_posy   float64 = 0
	CONST_PI            float64 = 3.1415926535
	CONST_PI2           float64 = CONST_PI / 2
	CONST_PI3           float64 = (3 * CONST_PI) / 2
	CONST_DR            float64 = 0.0174533 // one radian in degrees
	player_pos_x        float64 = 115
	player_pos_y        float64 = 220
	player_delta_x      float64 = 0
	player_delta_y      float64 = 0
	player_angle        float64 = 0 // Initial angle
	mapX                int     = 16
	mapY                int     = 16
	max_dof             int     = 16
	boot                        = 42
	mpv_run             []byte
	engine_version      = "ray_engine 0.6.5"
	debug_str           = "'Arrow to move, 'k' to exit, 'i' for debug info, 'n' toogle Y-shearing"
	debug_str2          = "'m' to hide 2D map and enter GUN mode, 'f' for fullscreen, j toogle minimap"
	str                 string
	/*map_array           = [64]int{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 1, 1, 0, 0, 1,
		1, 0, 0, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
	}*/
	map_array = [256]int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
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
	x3d                                float64 = 0 // Current offset to start to draw the 3D map
	x3d_orig                           float64 = 0 // Offset toogle - 0 = 2D MAP off, 530 = 2D map on - set at compile time
	ca                                 float64 = 0
	COLOR_R, COLOR_G, COLOR_B          uint8
	// Mouse vars
	x int = 0
	y int = 0
	// Gun vars
	gunx          float64 = 512
	guny          float64 = 350
	reset_gun_pos         = 0
)

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
		dof = max_dof
	}
	for dof < max_dof {
		mx = (int(rx) >> 6)
		my = (int(ry) >> 6)
		mp = my*mapX + mx
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 { // was == 1
			dof = max_dof // hit wall
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
		dof = max_dof
	}
	for dof < max_dof {
		mx = (int(rx) >> 6)
		my = (int(ry) >> 6)
		mp = my*mapX + mx
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 { // was == 1
			dof = max_dof // hit wall
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
	// Draw shortest ray based on disT in Orange
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
	ra = player_angle - (CONST_DR * 32) // numbers of rays / 2 - was 30
	if ra < 0 {
		ra = ra + (2 * CONST_PI)
	}
	if ra > (2 * CONST_PI) {
		ra = ra - (2 * CONST_PI)
	}

	for ray := 0; ray < 64; ray++ { // numbers of rays casted - was 60
		// Check horizontal lines
		check_horiz()
		// Check vertical lines + determine smallest line and draw on 2D map
		check_verti(screen)
		// Fix fisheye effect
		fix_fisheye()

		// Simple Shader - COLOR_R is handled here
		if disT >= 255 {
			COLOR_R = 255
		} else if disT <= 0 {
			COLOR_R = 0
		} else {
			COLOR_R = uint8(float64(disT))
		}

		// Simplistic and buggy collision detection FIXME
		if disT <= 16 {
			STATE_COLLISION = 1
		} else {
			STATE_COLLISION = 0
		}

		// Draw 3D lines/map section
		lineH = float64(64*320) / disT
		if lineH > float64(320) {
			lineH = float64(320)
		}

		// Basic Line offset
		if STATE_YSHEARING == 0 {
			lineO = 160 - (lineH / float64(3)) // was 2 (int) but 2.x helps with perspective somehow
		}

		// Dynamic line offset/Y Shearing - ALPHA TODO
		if STATE_YSHEARING == 1 {
			if y >= 96 && y <= 320 {
				lineO = float64(y) - (lineH / float64(3)) // was 2 (int) but 2.x helps with perspective somehow
			} else if y >= 320 {
				lineO = float64(320) - (lineH / float64(3)) // was 2 (int) but 2.x helps with perspective somehow
			} else if y <= 96 {
				lineO = float64(96) - (lineH / float64(3)) // was 2 (int) but 2.x helps with perspective somehow
			}
		}

		// x, y, ray@16px x64 = 1024, lineH
		ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(16), lineH+lineO, color.RGBA{COLOR_R, COLOR_G, COLOR_B, 255})
		x3d = x3d + 16
		ra = (ra + CONST_DR) // increment
	}
	x3d = x3d_orig
}

func (g *game) Draw(screen *ebiten.Image) {

	opBackground := &ebiten.DrawImageOptions{}
	opSplash := &ebiten.DrawImageOptions{}
	opwall := &ebiten.DrawImageOptions{}
	opGun := &ebiten.DrawImageOptions{}
	opCrosshair := &ebiten.DrawImageOptions{}
	opFire := &ebiten.DrawImageOptions{}

	// Images options
	opBackground.GeoM.Translate(bg_posx, bg_posy)
	opwall.GeoM.Translate(wall_posx, wall_posy)
	opGun.GeoM.Translate(gunx, guny)
	opCrosshair.GeoM.Translate(412, 240)
	opFire.GeoM.Translate(512, 385)

	// Reset user pos after firing
	if reset_gun_pos > 0 {
		reset_gun_pos = reset_gun_pos - 1
	} else {
		gunx = 512
		guny = 350
	}

	// Draw images
	if boot > 0 {
		//boot = boot - 1
		screen.DrawImage(splashImage, opSplash)
		if boot == 0 {
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
					if wall_posx < 960 { // 1024-64
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
			ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 4, 4, color.RGBA{255, 100, 100, 255})
		}

		// Raycasting
		cast_rays(screen)

		// Draw minimap
		draw_minimap(screen)

		// Draw player line of sight after raycast -> https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2/ebitenutil#DrawLine
		if STATE_SHOW_2D_MAP == 1 {
			ebitenutil.DrawLine(screen, player_pos_x, player_pos_y, player_pos_x+(player_delta_x*5), player_pos_y+(player_delta_y*5), color.RGBA{255, 255, 0, 255})
		} else {
			// Draw Gun and crosshair
			screen.DrawImage(gunImage, opGun)
			screen.DrawImage(crossHairImage, opCrosshair)
		}

		// Show debug info
		str = `{{.newline}} {{.engine_version}} {{.newline}} {{.debug_str}} {{.newline}} {{.debug_str2}} {{.newline}} {{.player_x}} {{.newline}} {{.player_y}} {{.newline}} {{.player_a}} {{.newline}} {{.state_ys}} {{.newline}} {{.state_ms}} {{.newline}} {{.state_co}}`
		str = strings.Replace(str, "{{.engine_version}}", engine_version, -1)
		str = strings.Replace(str, "{{.debug_str}}", debug_str, -1)
		str = strings.Replace(str, "{{.debug_str2}}", debug_str2, -1)
		str = strings.Replace(str, "{{.newline}}", "\n                                                                                                ", -1)
		str = strings.Replace(str, "{{.player_x}}", fmt.Sprintf("Player X : %f", player_pos_x), -1)
		str = strings.Replace(str, "{{.player_y}}", fmt.Sprintf("Player Y : %f", player_pos_y), -1)
		str = strings.Replace(str, "{{.player_a}}", fmt.Sprintf("Player Angle : %f", player_angle), -1)
		str = strings.Replace(str, "{{.state_ys}}", fmt.Sprintf("Y Shearing : %d", STATE_YSHEARING), -1)
		str = strings.Replace(str, "{{.state_ms}}", fmt.Sprintf("Mouse enabled : %d", STATE_MOUSE_SUPPORT), -1)
		str = strings.Replace(str, "{{.state_co}}", fmt.Sprintf("Collision : %d", STATE_COLLISION), -1)
		if STATE_SHOW_DEBUG == 1 {
			ebitenutil.DebugPrint(screen, str)
		}
	}

	/* Keyboard - FIXME When running in goroutine, it somewhat breaks the mouse support.
	not possible to move left/right either */
	keyboard_handling()

	// Mouse support is ALPHA
	if STATE_MOUSE_SUPPORT == 1 {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		x, y = ebiten.CursorPosition()
		player_angle = float64(x) / 360 // Breaks left/right on keyboard
		player_delta_x = math.Cos(player_angle) * 5
		player_delta_y = math.Sin(player_angle) * 5
		// Mouse buttons
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			go gun_sound()
			screen.DrawImage(fireImage, opFire)
			gunx = gunx + 1
			guny = guny - 1
			reset_gun_pos = 1
		}
	}
}

func (g *game) Update() error {
	return nil
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 512
}

func main() {
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("bg_ceiling_floor.png")
	if err != nil {
		log.Fatal(err)
	}
	splashImage, _, err = ebitenutil.NewImageFromFile("splash_dev.png")
	if err != nil {
		log.Fatal(err)
	}
	wallImage, _, err = ebitenutil.NewImageFromFile("wall.png")
	if err != nil {
		log.Fatal(err)
	}
	wallMiniImage, _, err = ebitenutil.NewImageFromFile("wall8.png")
	if err != nil {
		log.Fatal(err)
	}
	gunImage, _, err = ebitenutil.NewImageFromFile("gun2.png")
	if err != nil {
		log.Fatal(err)
	}
	crossHairImage, _, err = ebitenutil.NewImageFromFile("crosshair.png")
	if err != nil {
		log.Fatal(err)
	}
	fireImage, _, err = ebitenutil.NewImageFromFile("fire.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(engine_version)
	ebiten.SetWindowTitle(engine_version)
	ebiten.SetWindowSize(1024, 512)
	g := &game{}
	ebiten.RunGame(g)
}

func IsKeyTriggered(key ebiten.Key) bool {
	return keyStates[key] == 1
}

// CHANGEME
func gun_sound() {
	var errA error
	mpv_cmd := fmt.Sprintf("mpv --really-quiet --volume=42 gun.mp3")
	mpv_run, errA = exec.Command("bash", "-c", mpv_cmd).Output()
	if errA != nil {
		fmt.Printf("Error playing sound\n")
	}
}
