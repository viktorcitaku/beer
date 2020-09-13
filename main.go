package main

import (
	"log"
	"net/http"
	"os"
	"viktorcitaku.dev/beer/api"
	"viktorcitaku.dev/beer/internal"
)

//go:generate swagger generate spec
func main() {

	val, ok := os.LookupEnv("STATIC_FILES")
	if !ok {
		log.Printf("%s not set\n", "STATIC_FILES")
		val = "./web"
	} else {
		log.Printf("%s=%s\n", "STATIC_FILES", val)
	}

	internal.InitializeDatabase()

	http.HandleFunc("/api/save-user", internal.Chain(api.SaveUser, internal.Method("POST"), internal.Logging()))
	http.HandleFunc("/api/beers", internal.Chain(api.GetBeers, internal.Method("GET"), internal.Logging()))
	http.HandleFunc("/api/save-beer", internal.Chain(api.SaveBeer, internal.Method("POST"), internal.Logging()))
	http.HandleFunc("/api/user-beer-preferences", internal.Chain(api.GetUserBeerPreferences, internal.Method("GET"), internal.Logging()))
	http.HandleFunc("/api/update-user-beer-preferences", internal.Chain(api.UpdateUserBeerPreferences, internal.Method("POST"), internal.Logging()))
	http.Handle("/", internal.ChainExt(http.FileServer(http.Dir(val)), internal.Logging()))

	_ = http.ListenAndServe(":8080", nil)
}
