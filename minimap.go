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
				}
				if wall_minimap_posx < 250 { // 20 (offset) + (16x16 -> 256) - 16
					wall_minimap_posx += 16
				} else {
					wall_minimap_posx = 20
					wall_minimap_posy += 16
				}
			}
		}
		wall_minimap_posx = 20
		wall_minimap_posy = 20

		// Draw player as a rect
		ebitenutil.DrawRect(screen, float64((player_pos_x/4)+20), float64((player_pos_y/4)+20), 6, 6, color.RGBA{196, 255, 0, 255})
		//ebitenutil.DrawRect(screen, float64(player_pos_x/4), float64(player_pos_y/4), 2, 2, color.RGBA{255, 100, 100, 255})
	}
}
