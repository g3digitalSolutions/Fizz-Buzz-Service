package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Strategy Pattern
type Rule interface {
	Apply(n int) (string, bool)
}

type FizzRule struct{}

func (f FizzRule) Apply(n int) (string, bool) {
	if n%3 == 0 && n%5 != 0 {
		return "Fizz", true
	}
	return "", false
}

type BuzzRule struct{}

func (b BuzzRule) Apply(n int) (string, bool) {
	if n%5 == 0 && n%3 != 0 {
		return "Buzz", true
	}
	return "", false
}

type FizzBuzzRule struct{}

func (fb FizzBuzzRule) Apply(n int) (string, bool) {
	if n%3 == 0 && n%5 == 0 {
		return "FizzBuzz", true
	}
	return "", false
}

func fizzBuzzHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/fizzbuzz/")
	upperBound, err := strconv.Atoi(path) //Get upperbound
	if err != nil {                       //validate
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		return
	}

	//Get rules
	rules := []Rule{FizzBuzzRule{}, FizzRule{}, BuzzRule{}}

	//Initialize result
	json_result := map[string][]int{
		"Fizz":     {},
		"Buzz":     {},
		"FizzBuzz": {},
	}

	//Go through the elements and put in correct rule
	for i := 1; i <= upperBound; i++ {
		for _, rule := range rules {
			if key, ok := rule.Apply(i); ok {
				json_result[key] = append(json_result[key], i)
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(json_result); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func main() {
	//http.HandleFunc("/fizzbuzz/", fizzBuzzHandler)

	//fmt.Println("Running on localhost:8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))
  port := os.Getenv("PORT") // Heroku/Render will set this
	if port == "" {
		port = "8080" // fallback for local dev
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go web app!")
	})

	log.Printf("Running on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil)) // Bind to 0.0.0.0
}
