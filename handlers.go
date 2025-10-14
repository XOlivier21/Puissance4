package main

import (
    "html/template"
    "log"
    "net/http"
    "strconv"
)

var jeuActuel = NouveauJeu()
var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
    "Sequence": func(n int) []int {
        result := make([]int, n)
        for i := 0; i < n; i++ {
            result[i] = i
        }
        return result
    },
    "add": func(a, b int) int {
        return a + b
    },
}).ParseFiles("templates/index.html"))

func gestionnaireIndex(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    err := tmpl.ExecuteTemplate(w, "index.html", jeuActuel)
    if err != nil {
        log.Println("Erreur de template:", err)
        http.Error(w, "Erreur interne", http.StatusInternalServerError)
    }
}

func gestionnaireDeposer(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    colonneStr := r.FormValue("colonne")
    colonne, err := strconv.Atoi(colonneStr)
    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    jeuActuel.DeposerPion(colonne)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func gestionnaireReinitialiser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    jeuActuel.Reinitialiser()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}