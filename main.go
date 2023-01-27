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
	Id           int `json:"id"`
	Entrepreneur string `json:"entrepreneur"`
	PitchTitle 	 string `json:"pitchTitle"`
	PitchIdea    string `json:"pitchIdea"`
	AskAmount    float64 `json:"askAmount"`
	Equity       float64 `json:"equity"`
	Offers       []Offer `json:"offers"`
}

type Offer struct {
	Id       int `json:"id"`
	Investor string `json:"investor"`
	Amount   float64 `json:"amount"`
	Equity   float64 `json:"equity"`
	Comment  string `json:"comment"`
	PitchId  int `json:"pitch_id"`
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
	json.NewEncoder(w).Encode(id)
}

//POST
func postOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params:= mux.Vars(r)
	pitch_id := params["pitch_id"]
	var offer Offer
	_ = json.NewDecoder(r.Body).Decode(&offer)
	
	sqlStatement := `
	INSERT INTO offerdetails (investor, amount, equity, comment, pitch_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	err := db.QueryRow(sqlStatement, offer.Investor, offer.Amount, offer.Equity, offer.Comment, pitch_id).Scan(&id)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(id)
}

//GET
func getAllPitches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sqlStatement = `SELECT * FROM pitchdetails ORDER BY id DESC`
	var sqlStatementForOffers = `SELECT * FROM offerdetails WHERE pitch_id = $1`
	allPitches, _ := db.Query(sqlStatement)
	var pitchesWithOffers []Pitch
	
	for allPitches.Next() {
		//var err error
		var pitch Pitch
		_ = allPitches.Scan(&pitch.Id, &pitch.Entrepreneur, &pitch.PitchTitle, &pitch.PitchIdea, &pitch.AskAmount, &pitch.Equity)
		//fmt.Println("pitch", pitch)
		offersOfPitchId, _ := db.Query(sqlStatementForOffers, pitch.Id)
		//fmt.Println("offersofPitchId and _", offersOfPitchId, err)
		for offersOfPitchId.Next() {
			var offer Offer
			_ = offersOfPitchId.Scan(&offer.Id, &offer.Investor, &offer.Amount, &offer.Equity, &offer.Comment, &offer.PitchId)
			//fmt.Println("offer", offer)
			pitch.Offers = append(pitch.Offers, offer)
		}
		pitchesWithOffers = append(pitchesWithOffers, pitch)
	}
	json.NewEncoder(w).Encode(pitchesWithOffers)
}

//GET
func getPitch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params:= mux.Vars(r)
	pitch_id := params["pitch_id"]
	var sqlStatement = `SELECT * FROM pitchdetails WHERE id = $1`
	pitch := db.QueryRow(sqlStatement, pitch_id)
	var p Pitch
	_ = pitch.Scan(&p.Id, &p.Entrepreneur, &p.PitchTitle, &p.PitchIdea, &p.AskAmount, &p.Equity)
	json.NewEncoder(w).Encode(p)
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