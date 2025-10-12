package main

import "net/http"

func main() {
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/play", handlePlay)
    http.HandleFunc("/reset", handleReset)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.ListenAndServe(":8080", nil)
}
