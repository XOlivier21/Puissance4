package main

import "sync"

const (
    Lignes   = 6
    Colonnes = 7
    Vide     = 0
    Joueur1  = 1
    Joueur2  = 2
)

type Jeu struct {
    Plateau        [Lignes][Colonnes]int
    JoueurActuel   int
    Gagnant        int
    PartieTerminee bool
    mu             sync.Mutex
}

func NouveauJeu() *Jeu {
    return &Jeu{JoueurActuel: Joueur1}
}

func (j *Jeu) DeposerPion(colonne int) bool {
    j.mu.Lock()
    defer j.mu.Unlock()

    if j.PartieTerminee || colonne < 0 || colonne >= Colonnes {
        return false
    }

    for ligne := Lignes - 1; ligne >= 0; ligne-- {
        if j.Plateau[ligne][colonne] == Vide {
            j.Plateau[ligne][colonne] = j.JoueurActuel

            if j.verifierVictoire(ligne, colonne) {
                j.Gagnant = j.JoueurActuel
                j.PartieTerminee = true
            } else if j.plateauPlein() {
                j.PartieTerminee = true
            } else {
                j.changerJoueur()
            }
            return true
        }
    }
    return false
}

func (j *Jeu) changerJoueur() {
    if j.JoueurActuel == Joueur1 {
        j.JoueurActuel = Joueur2
    } else {
        j.JoueurActuel = Joueur1
    }
}

func (j *Jeu) verifierVictoire(ligne, colonne int) bool {
    joueur := j.Plateau[ligne][colonne]
    return j.verifierDirection(ligne, colonne, 0, 1, joueur) ||
        j.verifierDirection(ligne, colonne, 1, 0, joueur) ||
        j.verifierDirection(ligne, colonne, 1, 1, joueur) ||
        j.verifierDirection(ligne, colonne, 1, -1, joueur)
}

func (j *Jeu) verifierDirection(ligne, colonne, deltaLigne, deltaColonne, joueur int) bool {
    compte := 1

    for i := 1; i < 4; i++ {
        l, c := ligne+deltaLigne*i, colonne+deltaColonne*i
        if l < 0 || l >= Lignes || c < 0 || c >= Colonnes || j.Plateau[l][c] != joueur {
            break
        }
        compte++
    }

    for i := 1; i < 4; i++ {
        l, c := ligne-deltaLigne*i, colonne-deltaColonne*i
        if l < 0 || l >= Lignes || c < 0 || c >= Colonnes || j.Plateau[l][c] != joueur {
            break
        }
        compte++
    }

    return compte >= 4
}

func (j *Jeu) plateauPlein() bool {
    for colonne := 0; colonne < Colonnes; colonne++ {
        if j.Plateau[0][colonne] == Vide {
            return false
        }
    }
    return true
}

func (j *Jeu) Reinitialiser() {
    j.mu.Lock()
    defer j.mu.Unlock()
    j.Plateau = [Lignes][Colonnes]int{}
    j.JoueurActuel = Joueur1
    j.Gagnant = 0
    j.PartieTerminee = false
}