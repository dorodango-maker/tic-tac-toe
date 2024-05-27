package main

import (
    "fmt"
    "image/color"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Symbol int

const (
    Empty Symbol = iota
    X
    O
)

type Position struct {
    Row, Col int
}

type Game struct {
    board [3][3]Symbol
    turn Symbol
    xImg *ebiten.Image
    oImg *ebiten.Image
    xImgTransparent *ebiten.Image
    oImgTransparent *ebiten.Image
    winner Symbol
    xPositions []Position
    oPositions []Position
    oldestPosition Position
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
        board: [3][3]Symbol{},
        turn: X,
        xImg: xImg,
        oImg: oImg,
        xImgTransparent: xImgTransparent,
        oImgTransparent: oImgTransparent,
        winner: Empty,
        xPositions: []Position{},
        oPositions: []Position{},
        oldestPosition: Position{Row: -1, Col: -1},
    }
}

func loadAndResizeImage(filePath string, width, height int) (*ebiten.Image, error) {
    img, _, err := ebitenutil.NewImageFromFile(filePath)
    if err != nil {
        return nil, err
    }

    resizedImg := ebiten.NewImage(width, height)
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(float64(width)/float64(img.Bounds().Dx()), float64(height)/float64(img.Bounds().Dy()))
    resizedImg.DrawImage(img, op)

    return resizedImg, nil
}

func (g *Game) Update() error {
    g.handleWinnerState()
    g.handleGameProgression()
    return nil
}

func (g *Game) handleWinnerState() {
    if g.winner != Empty && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
        g.resetGame()
    }
}

func (g *Game) resetGame() {
    g.board = [3][3]Symbol{}
    g.turn = X
    g.winner = Empty
    g.xPositions = []Position{}
    g.oPositions = []Position{}
    g.oldestPosition = Position{Row: -1, Col: -1}
}

func (g *Game) handleGameProgression() {
    g.winner = checkWin(g.board)
    if g.winner != Empty {
        return
    }
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        g.processPlayerMove()
    }
}

func checkWin(board [3][3]Symbol) Symbol {
    for i := 0; i < 3; i++ {
        if board[i][0] != Empty && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
            return board[i][0]
        }
    }

    for i := 0; i < 3; i++ {
        if board[0][i] != Empty && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
            return board[0][i]
        }
    }

    if board[0][0] != Empty && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
        return board[0][0]
    }
    if board[0][2] != Empty && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
        return board[0][2]
    }

    return Empty
}

func (g *Game) processPlayerMove() {
    pos := g.getCursorPosition()
    if g.isValidMove(pos) {
        g.addMark(pos)
        g.toggleTurn()
    }
}

func (g *Game) getCursorPosition() Position {
    x, y := ebiten.CursorPosition()
    return Position{Row: y / 100, Col: x / 100}
}

func (g *Game) isValidMove(pos Position) bool {
    return pos.Row < 3 && pos.Col < 3 && g.board[pos.Row][pos.Col] == Empty && !(g.oldestPosition == pos)
}

func (g *Game) addMark(pos Position) {
    if g.turn == X {
        g.updatePositions(&g.xPositions, pos)
    } else {
        g.updatePositions(&g.oPositions, pos)
    }
    g.board[pos.Row][pos.Col] = g.turn
}

func (g *Game) updatePositions(positions *[]Position, pos Position) {
    if len(*positions) == 3 {
        oldest := (*positions)[0]
        g.board[oldest.Row][oldest.Col] = Empty
        *positions = (*positions)[1:]
        g.oldestPosition = pos
    }
    *positions = append(*positions, pos)
}

func (g *Game) toggleTurn() {
    if g.turn == X {
        g.turn = O
    } else {
        g.turn = X
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{0, 0, 0, 255})
    g.grid(screen)
    g.marks(screen)
    g.oldestMark(screen)
    g.winnerMessage(screen)
}

func (g *Game) grid(screen *ebiten.Image) {
    const gridSize = 3
    const cellSize = 100
    for i := 1; i < gridSize; i++ {
        vector.StrokeLine(screen, float32(i*cellSize), 0, float32(i*cellSize), float32(gridSize*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
        vector.StrokeLine(screen, 0, float32(i*cellSize), float32(gridSize*cellSize), float32(i*cellSize), 2, color.RGBA{255, 255, 255, 255}, false)
    }
}

func (g *Game) marks(screen *ebiten.Image) {
    const cellSize = 100
    const iconSize = 75
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(col*cellSize+(cellSize-iconSize)/2), float64(row*cellSize+(cellSize-iconSize)/2))
            if g.board[row][col] == X {
                screen.DrawImage(g.xImg, op)
            } else if g.board[row][col] == O {
                screen.DrawImage(g.oImg, op)
            }
        }
    }
}

func (g *Game) oldestMark(screen *ebiten.Image) {
    if (g.turn == X && len(g.xPositions) == 3) || (g.turn == O && len(g.oPositions) == 3) {
        var oldest Position
        var img *ebiten.Image
        if g.turn == X {
            oldest = g.xPositions[0]
            img = g.xImgTransparent
        } else {
            oldest = g.oPositions[0]
            img = g.oImgTransparent
        }
        g.board[oldest.Row][oldest.Col] = Empty
        transparentMark(screen, oldest.Row, oldest.Col, img)
        g.oldestPosition = oldest
    }
}

func (g *Game) winnerMessage(screen *ebiten.Image) {
    if g.winner != Empty {
        msg := fmt.Sprintf("Player %d wins! Right-click to reset.", g.winner)
        ebitenutil.DebugPrint(screen, msg)
    }
}

func transparentMark(screen *ebiten.Image, row, col int, img *ebiten.Image) {
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
