package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1340"
	dbname   = "xharktankGo"
)

var db *sql.DB

type Pitch struct {
	Entrepreneur string `json:"entrepreneur"`
	PitchTitle 	 string `json:"pitchTitle"`
	PitchIdea    string `json:"pitchIdea"`
	AskAmount    int `json:"askAmount"`
	Equity       int `json:"equity"`
}

type Offer struct {
	Investor string `json:"investor"`
	Amount int `json:"amount"`
	Equity int `json:"equity"`
	Comment string `json:"comment"`
	PitchId int `json:"pitch_id"`
}


//POST
func postPitch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var pitch Pitch
	
	_ = json.NewDecoder(r.Body).Decode(&pitch)
	sqlStatement := `
	INSERT INTO pitchdetails (entrepreneur, pitchtitle, pitchidea, askamount, equity)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	err := db.QueryRow(sqlStatement, pitch.Entrepreneur, pitch.PitchTitle, pitch.PitchIdea, pitch.AskAmount, pitch.Equity).Scan(&id)
	if err != nil {
		panic(err)
	}
	//fmt.Println("New record ID is:", id)
	json.NewEncoder(w).Encode(id)
}

//POST
func postOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var offer Offer
	_ = json.NewDecoder(r.Body).Decode(&offer)
	
	sqlStatement := `
	INSERT INTO offerdetails (investor, amount, equity, comment, pitch_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	err := db.QueryRow(sqlStatement, offer.Investor, offer.Amount, offer.Equity, offer.Comment, offer.PitchId).Scan(&id)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(id)
}

//GET
func getAllPitches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
}

//GET
func getPitch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func main () {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	  "password=%s dbname=%s sslmode=disable",
	  host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
	  panic(err)
	}
	defer db.Close()
  
	err = db.Ping()
	if err != nil {
	  panic(err)
	}
	fmt.Println("Successfully connected!")


	r := mux.NewRouter()

	r.HandleFunc("/pitches", postPitch).Methods("POST")
	r.HandleFunc("/pitches", getAllPitches).Methods("GET")
	r.HandleFunc("/pitches/{pitch_id}/makeOffer", postOffer).Methods("POST")
	r.HandleFunc("/pitches/{pitch_id}", getPitch).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}