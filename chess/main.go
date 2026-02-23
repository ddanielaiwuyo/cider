package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	boardSize  = 8
	squareSize = 80
)

type Game struct {
	lightSquare *ebiten.Image
	darkSquare  *ebiten.Image
}

func NewGame() *Game {
	g := &Game{}

	// Create light square
	g.lightSquare = ebiten.NewImage(squareSize, squareSize)
	g.lightSquare.Fill(color.RGBA{240, 217, 181, 255})

	// Create dark square
	g.darkSquare = ebiten.NewImage(squareSize, squareSize)
	g.darkSquare.Fill(color.RGBA{181, 136, 99, 255})

	return g
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {

			var square *ebiten.Image
			if (row+col)%2 == 0 {
				square = g.lightSquare
			} else {
				square = g.darkSquare
			}

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(
				float64(col*squareSize),
				float64(row*squareSize),
			)

			screen.DrawImage(square, opts)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return boardSize * squareSize, boardSize * squareSize
}

func main() {
	ebiten.SetWindowTitle("Chess Board")
	ebiten.SetWindowSize(boardSize*squareSize, boardSize*squareSize)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
