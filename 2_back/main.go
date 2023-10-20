package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "crypto.db"
const apiURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"

type CryptoCurrency struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}

func fetchCryptoCurrencies() ([]CryptoCurrency, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var currencies []CryptoCurrency
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

func createDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cryptocurrencies (
			id TEXT PRIMARY KEY,
			name TEXT,
			symbol TEXT,
			market_cap REAL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func updateDatabasePeriodically() {
	for {
		currencies, err := fetchCryptoCurrencies()
		if err != nil {
			log.Println("Failed to fetch data from API:", err)
		} else {
			db, err := createDatabase()
			if err != nil {
				log.Println("Failed to create database:", err)
			} else {
				for _, currency := range currencies {
					_, err := db.Exec("REPLACE INTO cryptocurrencies (id, name, symbol, market_cap) VALUES (?, ?, ?, ?)",
						currency.ID, currency.Name, currency.Symbol, currency.CurrentPrice)
					if err != nil {
						log.Println("Failed to insert data into database:", err)
					}
				}
				db.Close()
			}
		}

		time.Sleep(10 * time.Minute)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	db, err := createDatabase()
	if err != nil {
		log.Println("Failed to create database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var currencies []CryptoCurrency
	rows, err := db.Query("SELECT id, name, symbol, market_cap FROM cryptocurrencies")
	if err != nil {
		log.Println("Failed to query database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var currency CryptoCurrency
		err := rows.Scan(&currency.ID, &currency.Name, &currency.Symbol, &currency.CurrentPrice)
		if err != nil {
			log.Println("Failed to scan row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		currencies = append(currencies, currency)
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Println("Failed to parse HTML template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, currencies)
}

func main() {
	go updateDatabasePeriodically()

	http.HandleFunc("/", handleRequest)
	log.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
