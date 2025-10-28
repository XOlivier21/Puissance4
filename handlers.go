package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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

func GestionnaireMenu(w http.ResponseWriter, r *http.Request) {
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

func GestionnaireNouvellePartie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	mode := r.FormValue("mode")
	modePouvoirs := mode == "pouvoirs" || mode == "ia-pouvoirs"
	modeIA := mode == "ia" || mode == "ia-pouvoirs"

	if modeIA {
		niveauStr := r.FormValue("niveau")
		niveau, _ := strconv.Atoi(niveauStr)
		if niveau < 1 || niveau > 3 {
			niveau = 2
		}
		jeuActuel = NouveauJeuIA(modePouvoirs, niveau)
	} else {
		jeuActuel = NouveauJeu(modePouvoirs)
	}

	http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

func GestionnaireIndex(w http.ResponseWriter, r *http.Request) {
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

func GestionnaireDeposer(w http.ResponseWriter, r *http.Request) {
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

	// Si c'est le mode IA et que le jeu n'est pas terminé, l'IA joue
	if jeuActuel.ModeIA && !jeuActuel.PartieTerminee && jeuActuel.JoueurActuel == Joueur2 {
		time.Sleep(500 * time.Millisecond) // Petit délai pour l'effet
		coupIA := jeuActuel.IA.ChoisirCoup(jeuActuel)
		if coupIA != -1 {
			jeuActuel.DeposerPion(coupIA)
		}
	}

	http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

func GestionnairePouvoir(w http.ResponseWriter, r *http.Request) {
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

func GestionnaireReinitialiser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	jeuActuel.Reinitialiser()
	http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

// GestionnaireDeposerAPI gère les requêtes AJAX depuis le client (fetch '/api/deposer').
func GestionnaireDeposerAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	colonneStr := r.FormValue("colonne")
	colonne, err := strconv.Atoi(colonneStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false})
		return
	}

	// Tenter de déposer le pion
	moved := jeuActuel.DeposerPion(colonne)

	// Si mode IA et la partie n'est pas terminée et c'est au tour de l'IA, la faire jouer
	if jeuActuel.ModeIA && !jeuActuel.PartieTerminee && jeuActuel.JoueurActuel == Joueur2 {
		time.Sleep(500 * time.Millisecond)
		coupIA := jeuActuel.IA.ChoisirCoup(jeuActuel)
		if coupIA != -1 {
			jeuActuel.DeposerPion(coupIA)
		}
	}

	// Construire le plateau sous forme de [][]int pour l'encodage JSON
	plateau := make([][]int, Lignes)
	for i := 0; i < Lignes; i++ {
		plateau[i] = make([]int, Colonnes)
		for j := 0; j < Colonnes; j++ {
			plateau[i][j] = jeuActuel.Plateau[i][j]
		}
	}

	// Préparer la réponse
	resp := map[string]interface{}{
		"success":        moved,
		"plateau":        plateau,
		"partieTerminee": jeuActuel.PartieTerminee,
		"gagnant":        jeuActuel.Gagnant,
		"joueurActuel":   jeuActuel.JoueurActuel,
	}

	if jeuActuel.ModePouvoirs {
		resp["pouvoisJ1"] = jeuActuel.PouvoisJ1
		resp["pouvoisJ2"] = jeuActuel.PouvoisJ2
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func GestionnaireRetourMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
