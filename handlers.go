package main

import (
    "html/template"
    "log"
    "net/http"
    "strconv"
)

var jeuActuel *Jeu
var tmpl *template.Template

func init() {
    jeuActuel = NouveauJeu(false)
    var err error
    tmpl = template.Must(template.New("").Funcs(template.FuncMap{
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
    }).ParseGlob("templates/*.html"))
    
    if err != nil {
        log.Fatal(err)
    }
}

func gestionnaireMenu(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    err := tmpl.ExecuteTemplate(w, "menu.html", nil)
    if err != nil {
        log.Println("Erreur de template:", err)
        http.Error(w, "Erreur interne", http.StatusInternalServerError)
    }
}

func gestionnaireNouvellePartie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    mode := r.FormValue("mode")
    modePouvoirs := mode == "pouvoirs"
    
    jeuActuel = NouveauJeu(modePouvoirs)
    http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

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
        http.Redirect(w, r, "/jeu", http.StatusSeeOther)
        return
    }

    jeuActuel.DeposerPion(colonne)
    http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

func gestionnairePouvoir(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    pouvoirStr := r.FormValue("pouvoir")
    pouvoir, _ := strconv.Atoi(pouvoirStr)
    
    colonneStr := r.FormValue("colonne")
    
    switch pouvoir {
    case PouvoirBombe, PouvoirDouble:
        colonne, err := strconv.Atoi(colonneStr)
        if err != nil {
            http.Redirect(w, r, "/jeu", http.StatusSeeOther)
            return
        }
        if pouvoir == PouvoirBombe {
            jeuActuel.UtiliserPouvoirBombe(colonne)
        } else {
            jeuActuel.UtiliserPouvoirDouble(colonne)
        }
    case PouvoirAnnulation:
        jeuActuel.UtiliserPouvoirAnnulation()
    }

    http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

func gestionnaireReinitialiser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    jeuActuel.Reinitialiser()
    http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

func gestionnaireRetourMenu(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}