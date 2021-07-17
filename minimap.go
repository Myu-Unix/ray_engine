package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func draw_minimap(screen *ebiten.Image) {
	if STATE_MINIMAP == 1 {
		for k := 0; k < mapX; k++ {
			for l := 0; l < mapY; l++ {
				if map_array[k*mapX+l] == 1 {
					opWallMinimap := &ebiten.DrawImageOptions{}
					opWallMinimap.GeoM.Translate(wall_minimap_posx, wall_minimap_posy)
					screen.DrawImage(wallMiniImage, opWallMinimap)
				} else if map_array[k*mapX+l] == 2 {
					opEnemyMinimap := &ebiten.DrawImageOptions{}
					opEnemyMinimap.GeoM.Translate(wall_minimap_posx, wall_minimap_posy)
					screen.DrawImage(wallEnemyImage, opEnemyMinimap)
				}
				if wall_minimap_posx <= 132 { // 20 (offset) + (16x8 -> 128) - 16
					wall_minimap_posx += 8
				} else {
					wall_minimap_posx = 20 // offset
					wall_minimap_posy += 8
				}
			}
		}
		// Minimap is place on the lower left corner of the screen
		wall_minimap_posx = 20
		wall_minimap_posy = 350

		// Draw player as a rect
		ebitenutil.DrawRect(screen, float64((player_pos_x/2)+wall_minimap_posx), float64((player_pos_y/2)+wall_minimap_posy), 4, 4, color.RGBA{196, 255, 0, 255})
	}
}
