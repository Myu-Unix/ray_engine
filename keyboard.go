package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func keyboard_handling() {

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
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		keyStates[ebiten.KeyL]++
	} else {
		keyStates[ebiten.KeyL] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		keyStates[ebiten.KeyI]++
	} else {
		keyStates[ebiten.KeyI] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		keyStates[ebiten.KeyJ]++
	} else {
		keyStates[ebiten.KeyJ] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		keyStates[ebiten.KeyEnter]++
	} else {
		keyStates[ebiten.KeyEnter] = 0
	}
	if STATE_COLLISION == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) { // Z Azerty
			player_pos_x = player_pos_x + player_delta_x
			player_pos_y = player_pos_y + player_delta_y
			// small bobbing
			if guny >= 350 && guny < 386 {
				guny += 12
			} else {
				guny = 350
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) { // S
			player_pos_x = player_pos_x - player_delta_x
			player_pos_y = player_pos_y - player_delta_y
			// small bobbing
			if guny >= 314 && guny <= 350 {
				guny -= 12
			} else {
				guny = 350
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) { // Q Azerty
			player_pos_x = player_pos_x - player_strafe_delta_x
			player_pos_y = player_pos_y - player_strafe_delta_y
			// small bobbing
			if gunx > 476 && gunx <= 512 {
				gunx -= 12
			} else {
				gunx = 512
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) { // D
			player_pos_x = player_pos_x + player_strafe_delta_x
			player_pos_y = player_pos_y + player_strafe_delta_y
			// small bobbing
			if gunx >= 512 && gunx < 548 {
				gunx += 12
			} else {
				gunx = 512
			}
		}
	}
	if IsKeyTriggered(ebiten.KeyM) {
		if STATE_SHOW_2D_MAP == 0 {
			STATE_SHOW_2D_MAP = 1
			STATE_MINIMAP = 0
			x3d_orig = 530
		} else {
			STATE_SHOW_2D_MAP = 0
			STATE_MINIMAP = 1
			x3d_orig = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyJ) {
		if STATE_MINIMAP == 0 {
			STATE_MINIMAP = 1
		} else {
			STATE_MINIMAP = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyL) {
		if STATE_SCANLINES == 0 {
			STATE_SCANLINES = 1
		} else {
			STATE_SCANLINES = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyF) {
		if STATE_FULLSCREEN == 0 {
			ebiten.SetFullscreen(true)
			STATE_FULLSCREEN = 1
		} else {
			ebiten.SetFullscreen(false)
			STATE_FULLSCREEN = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyI) {
		if STATE_SHOW_DEBUG == 0 {
			STATE_SHOW_DEBUG = 1
		} else {
			STATE_SHOW_DEBUG = 0
		}
	}
	if IsKeyTriggered(ebiten.KeyK) {
		os.Exit(0)
	}
	if IsKeyTriggered(ebiten.KeyEnter) {
		boot = 0
	}
}
