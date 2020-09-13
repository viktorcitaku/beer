package internal

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"hash/fnv"
	"io"
	"log"
	"time"
)

type User struct {
	Email       string
	LastUpdated time.Time
}

// UserBeerPreferences holds additional information for selected beers
type UserBeerPreferences struct {
	UniqueID           int
	BeerID             int
	BeerName           string
	UserEmail          string
	DrunkTheBeerBefore bool
	GotDrunk           bool
	LastTime           string
	Rating             int
	Comment            string
}

// DB ...
var DB *memdb.MemDB

// InitializeDatabase ...
func InitializeDatabase() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"User": {
				Name: "User",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:         "id",
						AllowMissing: false,
						Unique:       true,
						Indexer:      &memdb.StringFieldIndex{Field: "Email"},
					},
				},
			},
			"UserBeerPreferences": {
				Name: "UserBeerPreferences",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:         "id",
						AllowMissing: false,
						Unique:       true,
						Indexer:      &memdb.IntFieldIndex{Field: "UniqueID"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	DB = db

	ExecuteTicker()
}

// ExecuteTicker ...
func ExecuteTicker() {
	go func() {
		for range time.Tick(time.Minute) {
			logString := "\n===== START DB CLEANUP ====\n"

			txn := DB.Txn(true) // here is an issue
			txn.TrackChanges()

			it, err := txn.Get("User", "id")
			if err != nil {
				log.Println(err)
				return
			}

			now := time.Now()

			for obj := it.Next(); obj != nil; obj = it.Next() {
				u := obj.(*User)
				if u != nil && u.LastUpdated.Sub(now) > 5 {
					logString += fmt.Sprintf("Beer Preference for user: %s are going to be deleted!\n", u.Email)
					preferences := UserBeerPreferencesByEmail(u.Email)
					if preferences != nil {
						for _, p := range preferences {
							logString += fmt.Sprintf("BeerID: %d, BeerName: %s, is going to be deleted!\n", p.BeerID, p.BeerName)
							_ = txn.Delete("UserBeerPreferences", p)
						}
					}

					logString += fmt.Sprintf("User: %s  is going to be deleted!\n", u.Email)
					_ = txn.Delete("User", u)
				}
			}

			changes := txn.Changes()
			if changes != nil {
				txn.Commit()
			} else {
				txn.Abort()
			}

			logString += "===== END DB CLEANUP ====\n"
			log.Printf(logString)
		}
	}()
}

// SaveUser ...
func SaveUser(email string) {
	user := &User{
		Email:       email,
		LastUpdated: time.Now(),
	}

	txn := DB.Txn(true)

	if err := txn.Insert("User", user); err != nil {
		log.Println(err)
	} else {
		// Commit the transaction
		txn.Commit()
	}
}

// FindUser ...
func FindUser(email string) *User {
	// Create read-only transaction
	txn := DB.Txn(false)
	defer txn.Abort()

	// Lookup by email
	raw, err := txn.First("User", "id", email)
	if err != nil {
		log.Println(err)
	}

	return raw.(*User)
}

// UserBeerPreferencesByEmail ...
func UserBeerPreferencesByEmail(email string) []*UserBeerPreferences {
	// Create read-only transaction
	txn := DB.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("UserBeerPreferences", "id")
	if err != nil {
		log.Println(err)
		return []*UserBeerPreferences{}
	}

	// Populate list/slice
	var ubp []*UserBeerPreferences
	for obj := it.Next(); obj != nil; obj = it.Next() {
		u := obj.(*UserBeerPreferences)
		if u != nil && u.UserEmail == email {
			ubp = append(ubp, obj.(*UserBeerPreferences))
		}
	}

	return ubp
}

// UserBeerPreferencesByBeerId ...
func UserBeerPreferencesByBeerId(beerId int, email string) *UserBeerPreferences {
	// Create read-only transaction
	txn := DB.Txn(false)
	defer txn.Abort()

	beers := UserBeerPreferencesByEmail(email)
	for _, b := range beers {
		if b != nil && b.BeerID == beerId {
			return b
		}
	}

	return nil
}

// UserBeerPreferencesByUniqueId ...
func UserBeerPreferencesByUniqueId(beerID int, userEmail string) *UserBeerPreferences {
	// Create read-only transaction
	txn := DB.Txn(false)
	defer txn.Abort()

	uniqueID := GenerateUniqueID(beerID, userEmail)

	raw, err := txn.First("UserBeerPreferences", "id", uniqueID)
	if err != nil || raw == nil {
		return nil
	}

	return raw.(*UserBeerPreferences)
}

// SaveUserBeerPreferences ...
func SaveUserBeerPreferences(userBeerPreferences *UserBeerPreferences) {
	txn := DB.Txn(true)

	if userBeerPreferences.UniqueID == 0 {
		// Set the UniqueID
		userBeerPreferences.UniqueID = GenerateUniqueID(userBeerPreferences.BeerID, userBeerPreferences.UserEmail)
	}

	if err := txn.Insert("UserBeerPreferences", userBeerPreferences); err != nil {
		log.Println(err)
	} else {
		// Commit the transaction
		txn.Commit()
	}
}

func GenerateUniqueID(beerID int, userEmail string) int {
	h := fnv.New32a()
	_, _ = io.WriteString(h, userEmail)
	uniqueID := int(h.Sum32()) + beerID
	return uniqueID
}
