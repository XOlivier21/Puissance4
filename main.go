package main

import (
    "log"
    "net/http"
)

func main() {
    // Routes principales
    http.HandleFunc("/", gestionnaireMenu)
    http.HandleFunc("/nouvelle-partie", gestionnaireNouvellePartie)
    http.HandleFunc("/jeu", gestionnaireIndex)
    http.HandleFunc("/deposer", gestionnaireDeposer)
    http.HandleFunc("/pouvoir", gestionnairePouvoir)
    http.HandleFunc("/reinitialiser", gestionnaireReinitialiser)
    http.HandleFunc("/menu", gestionnaireRetourMenu)

    // Fichiers statiques (CSS)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Serveur démarre sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}