package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func draw_minimap(screen *ebiten.Image) {
	if SHOW_MINIMAP == 1 {
		// Draw minimap
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
				if wall_minimap_posx < 140 { // 20 (offset) + (16x8 -> 128) - 16
					wall_minimap_posx += 8
				} else {
					wall_minimap_posx = 20 // offset
					wall_minimap_posy += 8
				}
			}
		}
		wall_minimap_posx = 20
		wall_minimap_posy = 350

		// Draw player as a rect
		ebitenutil.DrawRect(screen, float64((player_pos_x/8)+20), float64((player_pos_y/8)+350), 4, 4, color.RGBA{196, 255, 0, 255})
		//ebitenutil.DrawRect(screen, float64(player_pos_x/4), float64(player_pos_y/4), 2, 2, color.RGBA{255, 100, 100, 255})
	}
}
