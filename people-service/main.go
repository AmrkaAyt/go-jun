package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv" // Import dotenv package
)

type Person struct {
	gorm.Model
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

var db *gorm.DB

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	dbURL := os.Getenv("postgres://postgres:773504ok@host:5432/go_junior\n")
	db, err = gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate the schema
	db.AutoMigrate(&Person{})
}

func fetchDataFromAPI(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	// Extract the relevant data from the API response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return ""
	}

	value, ok := data["value"].(string)
	if !ok {
		return ""
	}

	return value
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/people", getPeople).Methods("GET")
	r.HandleFunc("/people/{id}", getPerson).Methods("GET")
	r.HandleFunc("/people", createPerson).Methods("POST")
	r.HandleFunc("/people/{id}", updatePerson).Methods("PUT")
	r.HandleFunc("/people/{id}", deletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Person model remains the same

func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)
	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.First(&person, params["id"])
	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	// Enrich data
	age := fetchDataFromAPI("http://age-service:8081/enrich/age?name=" + person.Name)
	genderNationality := fetchDataFromAPI("http://gender-service:8082/enrich/gender-nationality?name=" + person.Name)

	person.Age, _ = strconv.Atoi(age)

	// Handle gender and nationality as map[string]interface{}
	var genderData map[string]interface{}
	if err := json.Unmarshal([]byte(genderNationality), &genderData); err == nil {
		person.Gender, _ = genderData["gender"].(string)
		person.Nationality, _ = genderData["nationality"].(string)
	}

	// Save to database
	db.Create(&person)
	json.NewEncoder(w).Encode(person)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.First(&person, params["id"])
	json.NewDecoder(r.Body).Decode(&person)

	// Enrich data
	age := fetchDataFromAPI("http://age-service:8081/enrich/age?name=" + person.Name)
	genderNationality := fetchDataFromAPI("http://gender-service:8082/enrich/gender-nationality?name=" + person.Name)

	person.Age, _ = strconv.Atoi(age)

	// Handle gender and nationality as map[string]interface{}
	var genderData map[string]interface{}
	if err := json.Unmarshal([]byte(genderNationality), &genderData); err == nil {
		person.Gender, _ = genderData["gender"].(string)
		person.Nationality, _ = genderData["nationality"].(string)
	}

	// Save to database
	db.Save(&person)
	json.NewEncoder(w).Encode(person)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.First(&person, params["id"])
	db.Delete(&person)
}
