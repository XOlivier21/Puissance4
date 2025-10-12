package main

import (
    "html/template"
    "net/http"
    "strconv"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var game = NewGame()

type TemplateData struct {
    Grid        [6][7]int
    CurrentTurn int
    GameOver    bool
    Winner      int
}

func convertGame(g *Game) TemplateData {
    var data TemplateData
    data.CurrentTurn = int(g.CurrentTurn)
    data.GameOver = g.GameOver
    data.Winner = int(g.Winner)

    for i := 0; i < 6; i++ {
        for j := 0; j < 7; j++ {
            data.Grid[i][j] = int(g.Grid[i][j])
        }
    }
    return data
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    data := convertGame(game)
    tmpl.Execute(w, data)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost || game.GameOver {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    col, err := strconv.Atoi(r.FormValue("column"))
    if err == nil {
        game.Play(col)
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleReset(w http.ResponseWriter, r *http.Request) {
    game = NewGame()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
