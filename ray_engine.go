/*
ray_engine by Myu-Unix - Sept 2020 - Jan 2022
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
	STATE_COLLISION     = 0
	STATE_MINIMAP       = 1
	STATE_SCANLINES     = 0
	backgroundImage     *ebiten.Image
	scanlinesImage      *ebiten.Image
	splashImage         *ebiten.Image
	wallImage           *ebiten.Image
	wallMiniImage       *ebiten.Image
	wallEnemyImage      *ebiten.Image
	gunImage            *ebiten.Image
	//Enemy3DImage        *ebiten.Image
	crossHairImage        *ebiten.Image
	fireImage             *ebiten.Image
	keyStates                     = map[ebiten.Key]int{}
	bg_posx               float64 = 0
	bg_posy               float64 = 0
	wall_posx             float64 = 0
	wall_posy             float64 = 0
	wall_minimap_posx     float64 = 20
	wall_minimap_posy     float64 = 350
	CONST_PI              float64 = 3.1415926535
	CONST_PI2             float64 = CONST_PI / 2
	CONST_PI3             float64 = (3 * CONST_PI) / 2
	CONST_DR              float64 = 0.0174533 // one radian in degrees
	player_pos_x          float64 = 36
	player_pos_y          float64 = 50
	player_delta_x        float64 = 0
	player_delta_y        float64 = 0
	player_strafe_delta_x float64 = 0
	player_strafe_delta_y float64 = 0
	player_angle          float64 = 0
	mapX                  int     = 16
	mapY                  int     = 16
	max_dof               int     = 16
	boot                          = 42
	mpv_run               []byte
	engine_version        = "ray_engine 0.8.0"
	debug_str             = "'Arrow : move, 'k' : exit, 'i' : debug info, 'l' : scanlines"
	debug_str2            = "'m' : Gun mode/2D map mode, 'f' : fullscreen, j : toogle minimap"
	str                   string
	map_array             = [256]int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
		1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
		1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1,
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
	ray                                int     = 0
	COLOR_R, COLOR_G, COLOR_B          uint8
	// Mouse vars
	x int = 0
	y int = 0
	// Gun vars
	gunx          float64 = 512 // 512 ith gun1
	guny          float64 = 350 // 350 with gun1
	reset_gun_pos         = 0
	// Ballistic vars
	center_x        float64 = 480 // 1024/2 - 32 (which is half the crosshair size)
	center_y        float64 = 224 // 512/2 - 32 (which is half the crosshair size)
	override_colors         = 0
	ENEMY_SIGHT             = 0
	enemy_mp        int     = 0
	MP_120_BLIT     int     = 0
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

func check_horiz(screen *ebiten.Image) {
	dof = 0
	disH = 1000000
	hx = player_pos_x
	hy = player_pos_y
	aTan = -1 / math.Tan(ra)
	if ra > CONST_PI {
		ry = float64(((int(player_pos_y) >> 4) << 4)) - float64(0.0001)
		rx = (player_pos_y-ry)*aTan + player_pos_x
		yo = -16
		xo = -yo * aTan
	}
	if ra < CONST_PI {
		ry = float64(((int(player_pos_y) >> 4) << 4)) + float64(16)
		rx = (player_pos_y-ry)*aTan + player_pos_x
		yo = 16
		xo = -yo * aTan
	}
	// Looking left or right, will not hit horizontal lines
	if ra == 0 || ra == CONST_PI {
		rx = player_pos_x
		ry = player_pos_y
		dof = max_dof
	}
	for dof < max_dof {
		mx = (int(rx) >> 4)
		my = (int(ry) >> 4)
		mp = my*mapX + mx
		// Regular walls values are '1', enemies are '2'
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 {
			dof = max_dof // hit wall
			hx = rx
			hy = ry
			// ax ay bx by ang
			disH = math.Sqrt(((hx - player_pos_x) * (hx - player_pos_x)) + ((hy - player_pos_y) * (hy - player_pos_y)))
		} else {
			rx = rx + xo
			ry = ry + yo
			dof++
		}
	}
	// Ballistics. ray 31 is the middle ray
	if ray == 31 {
		if mp > 0 && mp < mapX*mapY && map_array[mp] == 2 {
			ebitenutil.DebugPrint(screen, "\n// DEBUG : Enemy in sight")
			ENEMY_SIGHT = 1 // Mark "destructible"
			enemy_mp = mp   // Mark "destructible"
		} else {
			ENEMY_SIGHT = 0
			enemy_mp = 0
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
		rx = float64(((int(player_pos_x) >> 4) << 4)) - float64(0.0001)
		ry = (player_pos_x-rx)*nTan + player_pos_y
		xo = -16
		yo = -xo * nTan
	}
	if ra < CONST_PI2 || ra > CONST_PI3 {
		rx = float64(((int(player_pos_x) >> 4) << 4)) + float64(16)
		ry = (player_pos_x-rx)*nTan + player_pos_y
		xo = 16
		yo = -xo * nTan
	}
	// staight up or down
	if ra == 0 || ra == CONST_PI {
		rx = player_pos_x
		ry = player_pos_y
		dof = max_dof
	}
	for dof < max_dof {
		mx = (int(rx) >> 4)
		my = (int(ry) >> 4)
		mp = my*mapX + mx
		// Regular walls values are '1', enemies are '2'
		if mp > 0 && mp < mapX*mapY && map_array[mp] > 0 { // was == 1
			dof = max_dof // hit wall
			vx = rx
			vy = ry
			// ax ay bx by ang
			disV = math.Sqrt(((vx - player_pos_x) * (vx - player_pos_x)) + ((vy - player_pos_y) * (vy - player_pos_y)))
		} else {
			rx = rx + xo
			ry = ry + yo
			dof++
		}
	}
	// Ballistics. ray 31 is the middle ray
	if ray == 31 {
		if mp > 0 && mp < mapX*mapY && map_array[mp] == 2 {
			//ebitenutil.DebugPrint(screen, "\n// DEBUG : Enemy in sight")
			ENEMY_SIGHT = 1 // Mark "destructible"
			enemy_mp = mp   // Mark "destructible"
		} else { // Cannot uncomment those otherwise enemies detection is broken
			//ENEMY_SIGHT = 0
			//enemy_mp = 0
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

	for ray = 0; ray < 64; ray++ { // numbers of rays casted - was 60
		// Check horizontal lines
		check_horiz(screen)
		// Check vertical lines + determine smallest line and draw on 2D map
		check_verti(screen)
		// Fix fisheye effect
		fix_fisheye()

		// Simple Shading - COLOR_R is handled here
		if disT >= 128 {
			COLOR_R = 255
		} else if disT <= 0 {
			COLOR_R = 0
		} else {
			COLOR_R = uint8(float64(disT * 2))
		}
		//COLOR_R = 0

		// Simplistic and buggy collision detection FIXME
		if disT <= 4 {
			STATE_COLLISION = 1
		} else {
			STATE_COLLISION = 0
		}

		// Draw lines on 3D map section
		lineH = float64(16*800) / disT
		if lineH > float64(800) {
			lineH = float64(800)
		}

		// Basic Line offset
		lineO = 256 - float64(lineH/2)

		// x, y, ray@16px x64 = 1024, lineH
		if mp > 0 && mp < mapX*mapY && map_array[mp] == 2 {
			ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(16), lineH, color.RGBA{255, 255, 255, 255})
		} else {
			ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(16), lineH, color.RGBA{COLOR_R, COLOR_G, COLOR_B, 240})
			//ebitenutil.DrawRect(screen, float64(x3d), float64(lineO), float64(16), lineH+lineO, color.RGBA{COLOR_R, COLOR_G, COLOR_B, 240})
		}

		// Blit image on '2' blocks
		/*if mp > 0 && mp < mapX*mapY && map_array[mp] == 2 {
			if MP_120_BLIT == 0 {
				x3d_save := x3d + 32
				lineO_save := lineO + 72
				opX := &ebiten.DrawImageOptions{}
				opX.GeoM.Translate(x3d_save, lineO_save)
				screen.DrawImage(Enemy3DImage, opX)
				MP_120_BLIT = 1
			}
		}*/

		x3d = x3d + 16
		ra = (ra + CONST_DR) // increment
	}
	x3d = x3d_orig
	MP_120_BLIT = 0
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
	opCrosshair.GeoM.Translate(center_x, center_y)
	opFire.GeoM.Translate(512, 385)

	// Reset gun position after firing
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
		if STATE_SCANLINES == 1 {
			screen.DrawImage(scanlinesImage, opBackground)
		}
		if boot == 0 {
			// Initial values for PDX/PDY, only applied once
			player_delta_x = math.Cos(player_angle) * 2
			player_delta_y = math.Sin(player_angle) * 2
			// Add 90 degrees in radians to get the right angle of player_angle
			player_strafe_delta_x = math.Cos(player_angle+(math.Pi/2)) * 2
			player_strafe_delta_y = math.Sin(player_angle+(math.Pi/2)) * 2
		}
	} else {
		screen.DrawImage(backgroundImage, opBackground)

		// 2D Map if enabled
		show2DMap(screen)

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

		if STATE_SCANLINES == 1 {
			screen.DrawImage(scanlinesImage, opBackground)
		}

		// Show debug info
		str = `{{.below}} {{.newline}} {{.engine_version}} {{.newline}} {{.debug_str}} {{.newline}} {{.debug_str2}} {{.newline}} {{.player_x}} {{.newline}} {{.player_y}} {{.newline}} {{.player_a}} {{.newline}} {{.state_emp}} {{.newline}} {{.state_ms}} {{.newline}} {{.state_co}}`
		str = strings.Replace(str, "{{.engine_version}}", engine_version, -1)
		str = strings.Replace(str, "{{.debug_str}}", debug_str, -1)
		str = strings.Replace(str, "{{.debug_str2}}", debug_str2, -1)
		str = strings.Replace(str, "{{.below}}", "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n", -1)
		str = strings.Replace(str, "{{.newline}}", "\n                        ", -1)
		str = strings.Replace(str, "{{.player_x}}", fmt.Sprintf("Player X : %f", player_pos_x), -1)
		str = strings.Replace(str, "{{.player_y}}", fmt.Sprintf("Player Y : %f", player_pos_y), -1)
		str = strings.Replace(str, "{{.player_a}}", fmt.Sprintf("Player Angle : %f", player_angle), -1)
		str = strings.Replace(str, "{{.state_ms}}", fmt.Sprintf("Mouse enabled : %d", STATE_MOUSE_SUPPORT), -1)
		str = strings.Replace(str, "{{.state_emp}}", fmt.Sprintf("Enemy block : %d", enemy_mp), -1)
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
		player_delta_x = math.Cos(player_angle) * 2
		player_delta_y = math.Sin(player_angle) * 2
		// Add 90 degrees in radians to get the right angle of player_angle
		player_strafe_delta_x = math.Cos(player_angle+(math.Pi/2)) * 2
		player_strafe_delta_y = math.Sin(player_angle+(math.Pi/2)) * 2
		// Mouse buttons
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Ballistics
			if ENEMY_SIGHT == 1 && enemy_mp != 0 {
				//ebitenutil.DebugPrint(screen, "\n**** HIT !! ****")
				map_array[enemy_mp] = 0 // Destroy block
				ENEMY_SIGHT = 0
			}
			go gun_sound()
			screen.DrawImage(fireImage, opFire)
			gunx = gunx + 4
			guny = guny - 2
			reset_gun_pos = 4
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
	backgroundImage, _, err = ebitenutil.NewImageFromFile("imgs/bg_ceiling_floor.png")
	if err != nil {
		log.Fatal(err)
	}
	scanlinesImage, _, err = ebitenutil.NewImageFromFile("imgs/scanlines.png")
	if err != nil {
		log.Fatal(err)
	}
	splashImage, _, err = ebitenutil.NewImageFromFile("imgs/splash_08.png")
	if err != nil {
		log.Fatal(err)
	}
	wallImage, _, err = ebitenutil.NewImageFromFile("imgs/wall.png")
	if err != nil {
		log.Fatal(err)
	}
	wallMiniImage, _, err = ebitenutil.NewImageFromFile("imgs/wall8.png")
	if err != nil {
		log.Fatal(err)
	}
	wallEnemyImage, _, err = ebitenutil.NewImageFromFile("imgs/enemy8.png")
	if err != nil {
		log.Fatal(err)
	}
	gunImage, _, err = ebitenutil.NewImageFromFile("imgs/gun2.png")
	if err != nil {
		log.Fatal(err)
	}
	/*Enemy3DImage, _, err = ebitenutil.NewImageFromFile("imgs/enemy_test2.png")
	if err != nil {
		log.Fatal(err)
	}*/
	crossHairImage, _, err = ebitenutil.NewImageFromFile("imgs/crosshair.png")
	if err != nil {
		log.Fatal(err)
	}
	fireImage, _, err = ebitenutil.NewImageFromFile("imgs/fire.png")
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
	mpv_cmd := fmt.Sprintf("mpv --really-quiet --volume=42 sound/gun.mp3")
	mpv_run, errA = exec.Command("bash", "-c", mpv_cmd).Output()
	if errA != nil {
		fmt.Printf("Error playing sound\n")
	}
}

func show2DMap(screen *ebiten.Image) {
	// Draw 2D map
	if STATE_SHOW_2D_MAP == 1 {
		for i := 0; i < mapX; i++ {
			for j := 0; j < mapY; j++ {
				if map_array[i*mapX+j] > 0 {
					opwall := &ebiten.DrawImageOptions{}
					opwall.GeoM.Translate(wall_posx, wall_posy)
					screen.DrawImage(wallImage, opwall)
				}
				if wall_posx < 240 { // 256 (16*16) - 16 (wall size in px). Was 960 = 1024-64
					wall_posx += 16
				} else {
					wall_posx = 0
					wall_posy += 16
				}
			}
		}
		wall_posx = 0
		wall_posy = 0

		// Draw player as a rect
		ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 12, 12, color.Black)
		ebitenutil.DrawRect(screen, float64(player_pos_x), float64(player_pos_y), 4, 4, color.RGBA{255, 100, 100, 255})
	}

}
