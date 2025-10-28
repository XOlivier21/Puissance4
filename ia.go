package main

import (
    "math/rand"
    "time"
)

type IA struct {
    Niveau int // 1 = Facile, 2 = Moyen, 3 = Difficile
}

func NouvelleIA(niveau int) *IA {
    rand.Seed(time.Now().UnixNano())
    return &IA{Niveau: niveau}
}

func (ia *IA) ChoisirCoup(jeu *Jeu) int {
    switch ia.Niveau {
    case 1:
        return ia.coupFacile(jeu)
    case 2:
        return ia.coupMoyen(jeu)
    case 3:
        return ia.coupDifficile(jeu)
    default:
        return ia.coupFacile(jeu)
    }
}

// IA Facile : joue aléatoirement
func (ia *IA) coupFacile(jeu *Jeu) int {
    colonnesDisponibles := []int{}
    for col := 0; col < Colonnes; col++ {
        if jeu.Plateau[0][col] == Vide {
            colonnesDisponibles = append(colonnesDisponibles, col)
        }
    }
    
    if len(colonnesDisponibles) == 0 {
        return -1
    }
    
    return colonnesDisponibles[rand.Intn(len(colonnesDisponibles))]
}

// IA Moyen : bloque l'adversaire et cherche à gagner
func (ia *IA) coupMoyen(jeu *Jeu) int {
    // 1. Chercher un coup gagnant
    for col := 0; col < Colonnes; col++ {
        if ia.peutGagner(jeu, col, Joueur2) {
            return col
        }
    }
    
    // 2. Bloquer l'adversaire
    for col := 0; col < Colonnes; col++ {
        if ia.peutGagner(jeu, col, Joueur1) {
            return col
        }
    }
    
    // 3. Jouer au centre si possible
    if jeu.Plateau[0][3] == Vide {
        return 3
    }
    
    // 4. Sinon jouer aléatoirement
    return ia.coupFacile(jeu)
}

// IA Difficile : utilise minimax simplifié
func (ia *IA) coupDifficile(jeu *Jeu) int {
    meilleurScore := -10000
    meilleurCoup := -1
    
    for col := 0; col < Colonnes; col++ {
        if jeu.Plateau[0][col] == Vide {
            // Simuler le coup
            ligne := ia.simulerCoup(jeu, col, Joueur2)
            if ligne == -1 {
                continue
            }
            
            score := ia.evaluerPosition(jeu, ligne, col)
            
            // Annuler le coup
            jeu.Plateau[ligne][col] = Vide
            
            if score > meilleurScore {
                meilleurScore = score
                meilleurCoup = col
            }
        }
    }
    
    if meilleurCoup == -1 {
        return ia.coupMoyen(jeu)
    }
    
    return meilleurCoup
}

func (ia *IA) peutGagner(jeu *Jeu, colonne int, joueur int) bool {
    // Trouver la ligne où le pion tomberait
    ligne := -1
    for l := Lignes - 1; l >= 0; l-- {
        if jeu.Plateau[l][colonne] == Vide {
            ligne = l
            break
        }
    }
    
    if ligne == -1 {
        return false
    }
    
    // Simuler le coup
    jeu.Plateau[ligne][colonne] = joueur
    gagne := jeu.verifierVictoire(ligne, colonne)
    jeu.Plateau[ligne][colonne] = Vide
    
    return gagne
}

func (ia *IA) simulerCoup(jeu *Jeu, colonne int, joueur int) int {
    for ligne := Lignes - 1; ligne >= 0; ligne-- {
        if jeu.Plateau[ligne][colonne] == Vide {
            jeu.Plateau[ligne][colonne] = joueur
            return ligne
        }
    }
    return -1
}

func (ia *IA) evaluerPosition(jeu *Jeu, ligne int, colonne int) int {
    score := 0
    
    // Vérifier si c'est un coup gagnant
    if jeu.verifierVictoire(ligne, colonne) {
        return 1000
    }
    
    // Préférer le centre
    score += (3 - abs(colonne-3)) * 10
    
    // Compter les alignements de 2 et 3
    score += ia.compterAlignements(jeu, ligne, colonne, Joueur2) * 10
    
    // Bloquer l'adversaire
    score += ia.compterAlignements(jeu, ligne, colonne, Joueur1) * 5
    
    return score
}

func (ia *IA) compterAlignements(jeu *Jeu, ligne int, colonne int, joueur int) int {
    count := 0
    
    // Vérifier toutes les directions
    directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
    
    for _, dir := range directions {
        alignement := 1
        
        // Compter dans une direction
        for i := 1; i < 4; i++ {
            l, c := ligne+dir[0]*i, colonne+dir[1]*i
            if l >= 0 && l < Lignes && c >= 0 && c < Colonnes && jeu.Plateau[l][c] == joueur {
                alignement++
            } else {
                break
            }
        }
        
        // Compter dans l'autre direction
        for i := 1; i < 4; i++ {
            l, c := ligne-dir[0]*i, colonne-dir[1]*i
            if l >= 0 && l < Lignes && c >= 0 && c < Colonnes && jeu.Plateau[l][c] == joueur {
                alignement++
            } else {
                break
            }
        }
        
        if alignement >= 2 {
            count += alignement
        }
    }
    
    return count
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}