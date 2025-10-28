package main

import "sync"

const (
    Lignes   = 6
    Colonnes = 7
    Vide     = 0
    Joueur1  = 1
    Joueur2  = 2
)

// Types de pouvoirs
const (
    PouvoirAucun       = 0
    PouvoirBombe       = 1
    PouvoirDouble      = 2
    PouvoirAnnulation  = 3
)

type Jeu struct {
    Plateau        [Lignes][Colonnes]int
    JoueurActuel   int
    Gagnant        int
    PartieTerminee bool
    ModePouvoirs   bool
    ModeIA         bool
    NiveauIA       int
    IA             *IA
    PouvoisJ1      map[int]int
    PouvoisJ2      map[int]int
    Historique     []struct {
        Ligne   int
        Colonne int
        Joueur  int
    }
    mu sync.Mutex
}

func NouveauJeu(modePouvoirs bool) *Jeu {
    jeu := &Jeu{
        JoueurActuel: Joueur1,
        ModePouvoirs: modePouvoirs,
        ModeIA:       false,
    }
    
    if modePouvoirs {
        jeu.PouvoisJ1 = map[int]int{
            PouvoirBombe:      2,
            PouvoirDouble:     2,
            PouvoirAnnulation: 1,
        }
        jeu.PouvoisJ2 = map[int]int{
            PouvoirBombe:      2,
            PouvoirDouble:     2,
            PouvoirAnnulation: 1,
        }
    }
    
    return jeu
}

func NouveauJeuIA(modePouvoirs bool, niveauIA int) *Jeu {
    jeu := NouveauJeu(modePouvoirs)
    jeu.ModeIA = true
    jeu.NiveauIA = niveauIA
    jeu.IA = NouvelleIA(niveauIA)
    return jeu
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
            
            // Sauvegarder dans l'historique
            j.Historique = append(j.Historique, struct {
                Ligne   int
                Colonne int
                Joueur  int
            }{ligne, colonne, j.JoueurActuel})

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

func (j *Jeu) UtiliserPouvoirBombe(colonne int) bool {
    j.mu.Lock()
    defer j.mu.Unlock()

    if !j.ModePouvoirs || j.PartieTerminee {
        return false
    }

    pouvoirs := j.PouvoisJ1
    if j.JoueurActuel == Joueur2 {
        pouvoirs = j.PouvoisJ2
    }

    if pouvoirs[PouvoirBombe] <= 0 {
        return false
    }

    // Trouver la position du pion à bombarder
    var ligne int = -1
    for l := Lignes - 1; l >= 0; l-- {
        if j.Plateau[l][colonne] != Vide {
            ligne = l
            break
        }
    }

    if ligne == -1 {
        return false
    }

    // Détruire le pion central et les adjacents
    j.Plateau[ligne][colonne] = Vide
    
    // Adjacents (8 directions)
    directions := [][2]int{
        {-1, -1}, {-1, 0}, {-1, 1},
        {0, -1}, {0, 1},
        {1, -1}, {1, 0}, {1, 1},
    }
    
    for _, dir := range directions {
        nl, nc := ligne+dir[0], colonne+dir[1]
        if nl >= 0 && nl < Lignes && nc >= 0 && nc < Colonnes {
            j.Plateau[nl][nc] = Vide
        }
    }

    // Faire tomber les pions
    j.appliquerGravite()

    pouvoirs[PouvoirBombe]--
    j.changerJoueur()
    return true
}

func (j *Jeu) UtiliserPouvoirDouble(colonne int) bool {
    j.mu.Lock()
    defer j.mu.Unlock()

    if !j.ModePouvoirs || j.PartieTerminee {
        return false
    }

    pouvoirs := j.PouvoisJ1
    if j.JoueurActuel == Joueur2 {
        pouvoirs = j.PouvoisJ2
    }

    if pouvoirs[PouvoirDouble] <= 0 {
        return false
    }

    // Déposer 2 pions
    count := 0
    for ligne := Lignes - 1; ligne >= 0 && count < 2; ligne-- {
        if j.Plateau[ligne][colonne] == Vide {
            j.Plateau[ligne][colonne] = j.JoueurActuel
            count++
            
            if j.verifierVictoire(ligne, colonne) {
                j.Gagnant = j.JoueurActuel
                j.PartieTerminee = true
                break
            }
        }
    }

    if count > 0 {
        pouvoirs[PouvoirDouble]--
        if !j.PartieTerminee {
            j.changerJoueur()
        }
        return true
    }

    return false
}

func (j *Jeu) UtiliserPouvoirAnnulation() bool {
    j.mu.Lock()
    defer j.mu.Unlock()

    if !j.ModePouvoirs || j.PartieTerminee || len(j.Historique) < 2 {
        return false
    }

    pouvoirs := j.PouvoisJ1
    if j.JoueurActuel == Joueur2 {
        pouvoirs = j.PouvoisJ2
    }

    if pouvoirs[PouvoirAnnulation] <= 0 {
        return false
    }

    // Annuler les 2 derniers coups
    for i := 0; i < 2 && len(j.Historique) > 0; i++ {
        dernier := j.Historique[len(j.Historique)-1]
        j.Plateau[dernier.Ligne][dernier.Colonne] = Vide
        j.Historique = j.Historique[:len(j.Historique)-1]
    }

    pouvoirs[PouvoirAnnulation]--
    return true
}

func (j *Jeu) appliquerGravite() {
    for col := 0; col < Colonnes; col++ {
        // Collecter tous les pions non vides
        pions := []int{}
        for ligne := 0; ligne < Lignes; ligne++ {
            if j.Plateau[ligne][col] != Vide {
                pions = append(pions, j.Plateau[ligne][col])
            }
        }
        
        // Remplir la colonne du bas vers le haut
        for ligne := 0; ligne < Lignes; ligne++ {
            j.Plateau[ligne][col] = Vide
        }
        
        for i, pion := range pions {
            j.Plateau[Lignes-len(pions)+i][col] = pion
        }
    }
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
    j.Historique = nil
    
    if j.ModePouvoirs {
        j.PouvoisJ1 = map[int]int{
            PouvoirBombe:      2,
            PouvoirDouble:     2,
            PouvoirAnnulation: 1,
        }
        j.PouvoisJ2 = map[int]int{
            PouvoirBombe:      2,
            PouvoirDouble:     2,
            PouvoirAnnulation: 1,
        }
    }
}