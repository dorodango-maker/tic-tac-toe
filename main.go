package main

import (
    "image/color"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{}

func (g *Game) Update() error {
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    const gridSize = 3
    const cellSize = 100

    // 背景を黒で塗りつぶす
    screen.Fill(color.RGBA{0, 0, 0, 255})

    // グリッドを白で描画
    for i := 1; i < gridSize; i++ {
        vector.StrokeLine(screen, float32(i*cellSize), 0, float32(i*cellSize), float32(gridSize*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
        vector.StrokeLine(screen, 0, float32(i*cellSize), float32(gridSize*cellSize), float32(i*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 300, 300
}

func main() {
    ebiten.SetWindowSize(300, 300)
    ebiten.SetWindowTitle("Tic-Tac-Toe")
    if err := ebiten.RunGame(&Game{}); err != nil {
        log.Fatal(err)
    }
}
