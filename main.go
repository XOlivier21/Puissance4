package main

import (
    "log"
    "net/http"
)

func main() {
    // Routes principales
    http.HandleFunc("/", gestionnaireIndex)
    http.HandleFunc("/deposer", gestionnaireDeposer)
    http.HandleFunc("/reinitialiser", gestionnaireReinitialiser)

    // Fichiers statiques (CSS)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Serveur démarre sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}