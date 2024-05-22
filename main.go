package main

import (
    "image/color"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
    board [3][3]int // 0: 空, 1: X, 2: O
    turn  int       // 1: Xのターン, 2: Oのターン
    xImg  *ebiten.Image
    oImg  *ebiten.Image
}

func loadAndResizeImage(filePath string, width, height int) (*ebiten.Image, error) {
    img, _, err := ebitenutil.NewImageFromFile(filePath)
    if err != nil {
        return nil, err
    }

    // 画像をリサイズ
    resizedImg := ebiten.NewImage(width, height)
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(float64(width)/float64(img.Bounds().Dx()), float64(height)/float64(img.Bounds().Dy()))
    resizedImg.DrawImage(img, op)

    return resizedImg, nil
}

func NewGame() *Game {
    const iconSize = 75

    xImg, err := loadAndResizeImage("assets/x.png", iconSize, iconSize)
    if err != nil {
        log.Fatal(err)
    }

    oImg, err := loadAndResizeImage("assets/o.png", iconSize, iconSize)
    if err != nil {
        log.Fatal(err)
    }

    return &Game{
        board: [3][3]int{},
        turn:  1,
        xImg:  xImg,
        oImg:  oImg,
    }
}

func (g *Game) Update() error {
    // クリック位置に応じて、記号を配置
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        row := y / 100
        col := x / 100
        if row < 3 && col < 3 && g.board[row][col] == 0 {
            g.board[row][col] = g.turn
            if g.turn == 1 {
                g.turn = 2
            } else {
                g.turn = 1
            }
        }
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    const gridSize = 3
    const cellSize = 100
    const iconSize = 75

    // 背景を黒で塗りつぶす
    screen.Fill(color.RGBA{0, 0, 0, 255})

    // グリッドを白で描画
    for i := 1; i < gridSize; i++ {
        vector.StrokeLine(screen, float32(i*cellSize), 0, float32(i*cellSize), float32(gridSize*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
        vector.StrokeLine(screen, 0, float32(i*cellSize), float32(gridSize*cellSize), float32(i*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
    }

    // ボードに記号を描画
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(col*cellSize+(cellSize-iconSize)/2), float64(row*cellSize+(cellSize-iconSize)/2))
            if g.board[row][col] == 1 {
                screen.DrawImage(g.xImg, op)
            } else if g.board[row][col] == 2 {
                screen.DrawImage(g.oImg, op)
            }
        }
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 300, 300
}

func main() {
    ebiten.SetWindowSize(300, 300)
    ebiten.SetWindowTitle("Tic-Tac-Toe")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}
