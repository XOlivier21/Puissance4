package main

import (
    "log"
    "net/http"
)

func main() {
    // Routes principales
    http.HandleFunc("/", GestionnaireMenu)
    http.HandleFunc("/nouvelle-partie", GestionnaireNouvellePartie)
    http.HandleFunc("/jeu", GestionnaireIndex)
    http.HandleFunc("/deposer", GestionnaireDeposer)
    http.HandleFunc("/api/deposer", GestionnaireDeposerAPI) // Nouvelle route API
    http.HandleFunc("/pouvoir", GestionnairePouvoir)
    http.HandleFunc("/reinitialiser", GestionnaireReinitialiser)
    http.HandleFunc("/menu", GestionnaireRetourMenu)

    // Fichiers statiques (CSS)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Serveur démarre sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}