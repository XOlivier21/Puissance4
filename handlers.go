package main

import (
	"encoding/json"
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
	// si le mode est 'robot', activer le flag VsRobot
	if mode == "robot" {
		jeuActuel.VsRobot = true
	}
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

// Structures et handlers pour API JSON (utilisés par le JS côté client)
type apiDeposerReq struct {
	Colonne int `json:"colonne"`
}

type GameState struct {
	Plateau        [][]int     `json:"plateau"`
	JoueurActuel   int         `json:"joueurActuel"`
	PartieTerminee bool        `json:"partieTerminee"`
	Gagnant        int         `json:"gagnant"`
	ModePouvoirs   bool        `json:"modePouvoirs"`
	VsRobot        bool        `json:"vsRobot"`
	PouvoisJ1      map[int]int `json:"pouvoisJ1,omitempty"`
	PouvoisJ2      map[int]int `json:"pouvoisJ2,omitempty"`
}

func buildGameState() GameState {
	gs := GameState{
		Plateau:        make([][]int, Lignes),
		JoueurActuel:   jeuActuel.JoueurActuel,
		PartieTerminee: jeuActuel.PartieTerminee,
		Gagnant:        jeuActuel.Gagnant,
		ModePouvoirs:   jeuActuel.ModePouvoirs,
		VsRobot:        jeuActuel.VsRobot,
	}
	for i := 0; i < Lignes; i++ {
		gs.Plateau[i] = make([]int, Colonnes)
		for j := 0; j < Colonnes; j++ {
			gs.Plateau[i][j] = jeuActuel.Plateau[i][j]
		}
	}
	if jeuActuel.ModePouvoirs {
		gs.PouvoisJ1 = jeuActuel.PouvoisJ1
		gs.PouvoisJ2 = jeuActuel.PouvoisJ2
	}
	return gs
}

// API: déposer un pion sans rechargement. Si mode vsRobot activé, le robot joue immédiatement après.
func gestionnaireAPIDeposer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req apiDeposerReq
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// tenter de déposer pour le joueur courant
	jeuActuel.DeposerPion(req.Colonne)

	// si vsRobot activé et la partie n'est pas terminée, faire jouer le robot
	if jeuActuel.VsRobot && !jeuActuel.PartieTerminee {
		col := jeuActuel.ChoisirCoupRobot()
		if col >= 0 {
			jeuActuel.DeposerPion(col)
		}
	}

	// retourner l'état actuel
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildGameState())
}

// API: réinitialiser la partie et renvoyer état
func gestionnaireAPIReinitialiser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	jeuActuel.Reinitialiser()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildGameState())
}

// API: renvoyer l'état actuel (GET)
func gestionnaireAPIEtat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildGameState())
}
