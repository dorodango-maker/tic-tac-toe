package main

import (
    "fmt"
    "image/color"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Position struct {
    Row, Col int
}

type Game struct {
    board [3][3]int // 0: 空, 1: X, 2: O
    turn int        // 1: Xのターン, 2: Oのターン
    xImg *ebiten.Image
    oImg *ebiten.Image
    xImgTransparent *ebiten.Image
    oImgTransparent *ebiten.Image
    winner int       // 0: まだ勝者なし, 1: Xの勝利, 2: Oの勝利
    xPositions []Position
    oPositions []Position
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

    xImgTransparent, err := loadAndResizeImage("assets/x_transparent.png", iconSize, iconSize)
    if err != nil {
        log.Fatal(err)
    }

    oImgTransparent, err := loadAndResizeImage("assets/o_transparent.png", iconSize, iconSize)
    if err != nil {
        log.Fatal(err)
    }

    return &Game{
        board: [3][3]int{},
        turn: 1,
        xImg: xImg,
        oImg: oImg,
        xImgTransparent: xImgTransparent,
        oImgTransparent: oImgTransparent,
        winner: 0,
        xPositions: []Position{},
        oPositions: []Position{},
    }
}

func checkWin(board [3][3]int) int {
    // 横のチェック
    for i := 0; i < 3; i++ {
        if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
            return board[i][0]
        }
    }
    // 縦のチェック
    for i := 0; i < 3; i++ {
        if board[0][i] != 0 && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
            return board[0][i]
        }
    }
    // 斜めのチェック
    if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
        return board[0][0]
    }
    if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
        return board[0][2]
    }

    return 0
}

func (g *Game) addMark(row, col int) {
    pos := Position{Row: row, Col: col}
    if g.turn == 1 {
        if len(g.xPositions) == 3 {
            oldest := g.xPositions[0]
            g.board[oldest.Row][oldest.Col] = 0
            g.xPositions = g.xPositions[1:]
        }
        g.xPositions = append(g.xPositions, pos)
    } else {
        if len(g.oPositions) == 3 {
            oldest := g.oPositions[0]
            g.board[oldest.Row][oldest.Col] = 0
            g.oPositions = g.oPositions[1:]
        }
        g.oPositions = append(g.oPositions, pos)
    }
    g.board[row][col] = g.turn
}

func (g *Game) Update() error {
    if g.winner != 0 {
        return nil
    }

    // 勝利判定
    g.winner = checkWin(g.board)
    if g.winner != 0 {
        log.Printf("Player %d wins!", g.winner)
        return nil
    }

    // クリック位置に応じて、記号を配置
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        row := y / 100
        col := x / 100
        if row < 3 && col < 3 && g.board[row][col] == 0 {
            g.addMark(row, col)
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

    // 古い記号を半透明な画像で描画
    if g.turn == 1 && len(g.xPositions) == 3 {
        oldest := g.xPositions[0]
        // 元の記号を消す
        g.board[oldest.Row][oldest.Col] = 0
        drawTransparentMark(screen, oldest.Row, oldest.Col, g.xImgTransparent)
    } else if g.turn == 2 && len(g.oPositions) == 3 {
        oldest := g.oPositions[0]
        g.board[oldest.Row][oldest.Col] = 0
        drawTransparentMark(screen, oldest.Row, oldest.Col, g.oImgTransparent)
    }

    // 勝者のメッセージを表示
    if g.winner != 0 {
        msg := fmt.Sprintf("Player %d wins!", g.winner)
        ebitenutil.DebugPrint(screen, msg)
    }
}

func drawTransparentMark(screen *ebiten.Image, row, col int, img *ebiten.Image) {
    const cellSize = 100
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(col*cellSize+(cellSize-75)/2), float64(row*cellSize+(cellSize-75)/2))
    screen.DrawImage(img, op)
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
