package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/madraceee/url-shortner/internal/database"
)

func shortURL(longURL string) (uuid.UUID, string) {
	uuid := uuid.New()

	encoded := base64.StdEncoding.EncodeToString(uuid[:])

	if len(encoded) > 7 {
		encoded = encoded[:7]
	}

	return uuid, encoded
}

type payload struct {
	URL string `json:"url"`
}

func (apiCfg *apiConfig) handleUrlShorten(w http.ResponseWriter, r *http.Request) {

	godotenv.Load(".env")
	pageURL := os.Getenv("URL")

	type parameters struct {
		LongURL string `json:"longURL"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error Decoding JSON")
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	//Check Whether URL is present in database
	record, err := apiCfg.DB.FindExistingRecordUsingLongURL(r.Context(), params.LongURL)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		uuid, shortUrl := shortURL(params.LongURL)
		formattedShortURL := fmt.Sprintf("%v/%v", pageURL, shortUrl)
		apiCfg.DB.AddNewRecord(r.Context(), database.AddNewRecordParams{
			Uuid:     uuid,
			Longurl:  params.LongURL,
			Shorturl: shortUrl,
		})

		respondWithJSON(w, 200, payload{
			URL: formattedShortURL,
		})
		return
	}

	if err != nil {
		log.Println("Error while fetching results:", err)
		respondWithError(w, 500, err.Error())
		return
	}

	//If record exsists
	formattedShortURL := fmt.Sprintf("%v/%v", pageURL, record.Shorturl)
	respondWithJSON(w, 200, payload{
		URL: formattedShortURL,
	})
}

func (apiCfg *apiConfig) handleFetchShortUrl(w http.ResponseWriter, r *http.Request) {

	shortURL := chi.URLParam(r, "shortURL")

	record, err := apiCfg.DB.FindExistingRecordUsingShortURL(r.Context(), shortURL)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		respondWithError(w, 400, "Given Short URL does not Exists")
		return
	}
	if err != nil {
		log.Printf("Could not retrieve record from DB:%s", err)
	}

	// Check whether https is there or not hhtp is present
	// If not presen then add

	if !strings.Contains(record.Longurl, "https://") && !strings.Contains(record.Longurl, "http://") {
		record.Longurl = "http://" + record.Longurl
	}

	http.Redirect(w, r, record.Longurl, int(302))
}
