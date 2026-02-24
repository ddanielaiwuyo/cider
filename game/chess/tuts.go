package chess

import (
	"fmt"
	"image/color"
	// "math/rand/v2"
	"github.com/hajimehoshi/ebiten/v2"
)

// this probably represents the state of our game
type Game struct {
	counter            uint64
	winsizex, winsizey float64
}

func (g *Game) Update() error {
	if g.counter >= 1 {
		g.counter = 0
		return nil
	}

	g.counter += 1

	return nil
}

const square_size = 80
const row_length = 8

var baseSquare = ebiten.NewImage(square_size, square_size)

func (g *Game) Draw(screen *ebiten.Image) {
	// img := ebiten.NewImage(80, 80)
	var randomch uint8 = 212
	baseSquare.Fill(color.RGBA{241, randomch, 98, 11})
	for row := range row_length {
		for col := range row_length {
			opts := &ebiten.DrawImageOptions{}
			x := row * (square_size)
			y := col * (square_size)

			if (row+col)%2 == 0 {
				baseSquare.Fill(color.RGBA{240, 217, 181, 255})
			} else {
				baseSquare.Fill(color.RGBA{181, 136, 99, 255})
			}

			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(baseSquare, opts)
		}
	}
}

func (g *Game) Layout(outWidth, outHeight int) (int, int) {
	return ebiten.WindowSize()
}

const win_size = square_size * 8

func Run() error {
	ebiten.SetWindowSize(win_size, win_size)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := &Game{}
	g.counter = 0
	g.winsizex = win_size
	g.winsizey = win_size
	if err := ebiten.RunGame(g); err != nil {
		return fmt.Errorf(" couldn't run game: %w", err)
	}

	return nil
}
