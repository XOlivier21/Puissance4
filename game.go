package main

type Cell int

const (
    Empty Cell = iota
    Player1
    Player2
)

type Game struct {
    Grid        [6][7]Cell
    CurrentTurn Cell
    GameOver    bool
    Winner      Cell
}

func NewGame() *Game {
    return &Game{CurrentTurn: Player1}
}

func (g *Game) Play(column int) bool {
    if g.GameOver || column < 0 || column > 6 {
        return false
    }

    for i := 5; i >= 0; i-- {
        if g.Grid[i][column] == Empty {
            g.Grid[i][column] = g.CurrentTurn
            g.checkWin(i, column)
            if !g.GameOver {
                g.switchTurn()
            }
            return true
        }
    }
    return false
}

func (g *Game) switchTurn() {
    if g.CurrentTurn == Player1 {
        g.CurrentTurn = Player2
    } else {
        g.CurrentTurn = Player1
    }
}

func (g *Game) checkWin(row, col int) {
    player := g.Grid[row][col]
    dirs := [][2]int{
        {0, 1},  // horizontal
        {1, 0},  // vertical
        {1, 1},  // diagonale ↘
        {1, -1}, // diagonale ↙
    }

    for _, d := range dirs {
        count := 1
        for i := 1; i < 4; i++ {
            r, c := row+i*d[0], col+i*d[1]
            if r < 0 || r > 5 || c < 0 || c > 6 || g.Grid[r][c] != player {
                break
            }
            count++
        }
        for i := 1; i < 4; i++ {
            r, c := row-i*d[0], col-i*d[1]
            if r < 0 || r > 5 || c < 0 || c > 6 || g.Grid[r][c] != player {
                break
            }
            count++
        }
        if count >= 4 {
            g.GameOver = true
            g.Winner = player
            return
        }
    }

    // Vérifie si la grille est pleine (match nul)
    full := true
    for _, row := range g.Grid {
        for _, c := range row {
            if c == Empty {
                full = false
                break
            }
        }
    }
    if full {
        g.GameOver = true
        g.Winner = Empty
    }
}
