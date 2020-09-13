// Package api Beers API
//
// Title: Beers API
// Schemes: http
// Host: localhost
// BasePath: /v1
// Version: 0.0.1
// License: MIT http://opensource.org/licenses/MIT
// Contact: Viktor Citaku<viktor.citaku@gmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strings"
	"viktorcitaku.dev/beer/internal"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// SaveUser REST API
// swagger:route POST /api/save-user saveUser
//
// Saves user by given email.
//
// responses:
//   default: description: no body
//   200: description: no body
func SaveUser(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "BEER_SESSION")
	// Set some session values.
	session.Values["email"] = r.Form["email"][0]
	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the user in DB
	internal.SaveUser(r.Form["email"][0])

	w.WriteHeader(http.StatusOK)
}

// BeersPayload ...
// swagger:response getBeersResponse
type GetBeersResponse struct {
	// in: body
	Body []BeersPayload
}

// BeersPayload ...
// swagger:model beersPayload
type BeersPayload struct {
	// required: true
	ID int `json:"id"`
	// required: true
	Name string `json:"name"`
}

// GetBeers REST API
// swagger:route GET /api/beers getBeers
//
// Returns the list of beers.
//
// produces:
// - application/json
// responses:
//   default: getBeersResponse
func GetBeers(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "BEER_SESSION")
	fmt.Printf("Email: %s \n", session.Values["email"].(string))

	pf, b, f := validateQueryString(r)

	res, err := http.Get("https://api.punkapi.com/v2/beers")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var beers []BeersPayload

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {

		var body []map[string]interface{}

		err := json.NewDecoder(res.Body).Decode(&body)
		defer res.Body.Close()

		if err != nil {
			log.Println(err)
		}

		for _, obj := range body {

			pfMatch, bMatch, fMatch := filterBeers(pf, b, f, obj)

			if pfMatch && bMatch && fMatch {
				var bp BeersPayload
				bp.ID = int(obj["id"].(float64))
				bp.Name = obj["name"].(string)

				beers = append(beers, bp)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(beers)
	} else {
		log.Println("The content isn't JSON! => " + res.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func filterBeers(pf string, b string, f string, obj map[string]interface{}) (bool, bool, bool) {

	pfMatch := -1
	if pf != "" {
		foodPairings := obj["food_pairing"].([]interface{}) // this is insane
		for _, s := range foodPairings {
			for _, s1 := range strings.Split(pf, "_") {
				if strings.Contains(s.(string), s1) {
					pfMatch = 1
				}
			}
		}
	} else {
		// We assume that pairing food was empty so we include all beers
		pfMatch = 0
	}

	bMatch := -1
	if b != "" && obj["ibu"] != nil {
		ibu := obj["ibu"].(float64)
		switch b {
		case "high":
			if highIbuScale(ibu) {
				bMatch = 1
			}
		case "medium":
			if mediumIbuScale(ibu) {
				bMatch = 1
			}
		case "low":
			if lowIbuScale(ibu) {
				bMatch = 1
			}
		}
	} else {
		// We assume that bitterness was empty so we include all beers
		bMatch = 0
	}

	fMatch := -1
	if f != "" {
		// This line is insane (This might lead to errors as well)
		temp := (((obj["method"]).(map[string]interface{})["fermentation"]).(map[string]interface{})["temp"]).(map[string]interface{})["value"].(float64)
		switch f {
		case "top_fermented":
			if isAleBeer(temp) {
				fMatch = 1
			}
		case "bottom_fermented":
			if isLagerBeer(temp) {
				fMatch = 1
			}
		}
	} else {
		// We assume that fermentation was empty so we include all beers
		fMatch = 0
	}

	// The reason for -1 is to be sure that it was a mismatch, 0 means value was not provided and 1 value matched
	return pfMatch != -1, bMatch != -1, fMatch != -1
}

// validates and returns a triple
func validateQueryString(r *http.Request) (string, string, string) {
	var finalPairingFood string
	pairingFood := r.URL.Query()["pairing_food"]
	if pairingFood != nil && len(pairingFood) > 0 {
		finalPairingFood = pairingFood[0]
	}

	var finalBitterness string
	bitterness := r.URL.Query()["bitterness"]
	if bitterness != nil && len(bitterness) > 0 {
		finalBitterness = bitterness[0]
	}

	var finalFermentation string
	fermentation := r.URL.Query()["fermentation"]
	if fermentation != nil && len(fermentation) > 0 {
		finalFermentation = fermentation[0]
	}

	return finalPairingFood, finalBitterness, finalFermentation
}

func highIbuScale(ibu float64) bool {
	return ibu > 80 && ibu <= 120
}

func mediumIbuScale(ibu float64) bool {
	return ibu > 40 && ibu <= 80
}

func lowIbuScale(ibu float64) bool {
	return ibu >= 0 && ibu <= 40
}

func isAleBeer(temp float64) bool {
	return temp >= 10 && temp <= 25
}

func isLagerBeer(temp float64) bool {
	return temp >= 7 && temp <= 15
}

// SaveBeerPayload ...
// swagger:model saveBeerPayload
type SaveBeerPayload struct {
	// id
	// required: true
	ID int `json:"id"`
	// name
	// required: true
	Name string `json:"name"`
}

// SaveBeer REST API
// swagger:route POST /api/save-beer saveBeers
//
// Saves beer by given ID and Name.
//
// parameters:
// + name: saveBeerPayload
//   in: body
//   type: saveBeerPayload
//   required: true
// responses:
//   default: description: no body
//   200: description: no body
func SaveBeer(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "BEER_SESSION")
	fmt.Printf("Email: %s \n", session.Values["email"].(string))

	var sbp SaveBeerPayload

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {

		err := json.NewDecoder(r.Body).Decode(&sbp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// We check if an existing user is present
		user := *internal.FindUser(session.Values["email"].(string))
		if user == (internal.User{}) {
			log.Printf("User: %s was not present! New will be created.", session.Values["email"].(string))
		}
		// The save function does update at the same time
		internal.SaveUser(session.Values["email"].(string))

		// Search for UserBeerPreferences
		ubp := internal.UserBeerPreferencesByBeerId(sbp.ID, session.Values["email"].(string))
		if ubp == nil {
			log.Printf("UserBeerPreferences: %+v was not present! New one will be created.", ubp)
			ubp = &internal.UserBeerPreferences{} // create new object
		}

		ubp.BeerID = sbp.ID
		ubp.BeerName = sbp.Name
		ubp.UserEmail = session.Values["email"].(string)

		internal.SaveUserBeerPreferences(ubp)

		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("The content isn't JSON! => " + r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetUserBeerPreferencesResponse ...
// swagger:response getUserBeerPreferencesResponse
type GetUserBeerPreferencesResponse struct {
	// in: body
	Body []UserBeerPreferencesPayload
}

// UserBeerPreferencesPayload ...
// swagger:model userBeerPreferencesPayload
type UserBeerPreferencesPayload struct {
	// id
	// required: true
	ID int `json:"id"`
	// name
	// required: true
	Name               string `json:"name"`
	DrunkTheBeerBefore bool   `json:"drunk_before"`
	GotDrunk           bool   `json:"got_drunk"`
	LastTime           string `json:"last_time"`
	Rating             int    `json:"rating"`
	Comment            string `json:"comment"`
}

// GetUserBeerPreferences REST API
// swagger:route GET /api/user-beer-preferences getUserBeerPreferences
//
// Returns the list of UserBeerPreferences.
//
// produces:
// - application/json
// responses:
//   default: getUserBeerPreferencesResponse
//   200: getUserBeerPreferencesResponse
func GetUserBeerPreferences(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "BEER_SESSION")
	fmt.Printf("Email: %s \n", session.Values["email"].(string))

	payload := make([]UserBeerPreferencesPayload, 0)
	preferences := internal.UserBeerPreferencesByEmail(session.Values["email"].(string))
	if preferences == nil {
		log.Println("No preferences found!")
	} else {
		for _, p := range preferences {
			payload = append(payload, mapEntityToPayloadForUserBeerPreferences(p))
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}

func mapEntityToPayloadForUserBeerPreferences(e *internal.UserBeerPreferences) UserBeerPreferencesPayload {
	if e == nil {
		panic("Entity UserBeerPreferences is nil")
	}

	payload := UserBeerPreferencesPayload{
		ID:                 e.BeerID,
		Name:               e.BeerName,
		DrunkTheBeerBefore: e.DrunkTheBeerBefore,
		GotDrunk:           e.GotDrunk,
		LastTime:           e.LastTime,
		Rating:             e.Rating,
		Comment:            e.Comment,
	}

	return payload
}

// NOTE: The below Go Swagger doesn't generate correctly the request body, so it is required to be done manually

// UpdateUserBeerPreferences REST API
// swagger:route POST /api/update-user-beer-preferences updateUserBeerPreferences
//
// Updates beer preferences by given UserBeerPreferences array.
//
// parameters:
// + name: userBeerPreferencesPayload
//   in: body
//   required: true
//   schema:
//     type: array
//     items:
//       type: array
//       "$ref": "#/definitions/userBeerPreferencesPayload"
// responses:
//   default: description: no body
//   200: description: no body
func UpdateUserBeerPreferences(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "BEER_SESSION")
	fmt.Printf("Email: %s \n", session.Values["email"].(string))

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var payload []map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)

		internal.SaveUser(session.Values["email"].(string))

		for _, p := range payload {

			beerID := int(p["id"].(float64))
			userEmail := session.Values["email"].(string)

			ubp := internal.UserBeerPreferencesByUniqueId(beerID, userEmail)
			if ubp == nil {
				ubp = &internal.UserBeerPreferences{}
			}

			ubp.BeerID = beerID
			ubp.BeerName = p["name"].(string)
			ubp.UserEmail = userEmail
			ubp.DrunkTheBeerBefore = p["drunk_before"].(bool)
			ubp.GotDrunk = p["got_drunk"].(bool)
			ubp.LastTime = p["last_time"].(string)
			ubp.Rating = int(p["rating"].(float64))
			ubp.Comment = p["comment"].(string)

			internal.SaveUserBeerPreferences(ubp)
		}
	} else {
		log.Println("The content isn't JSON! => " + r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
