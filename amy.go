package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	LOCAL_HOST = "localhost:8000"
	HOST       = "https://api.edamam.com"
	PORT       = ":443"
	AUTH_ID    = "EDAMAM_ID"
	AUTH_KEY   = "EDAMAM_KEY"
	URL_LEN    = 9
)

type SearchResponse struct {
	Query string `json:"q"`
	Hits  []struct {
		Recipe struct {
			URL             string   `json:"url"`
			Label           string   `json:"label"`
			Img             string   `json:"image"`
			IngredientLines []string `json:"ingredientLines"`
		} `json:"recipe"`
	} `json:"hits"`
}

func createURL(componentArray [URL_LEN]string) string {
	var urlBuffer bytes.Buffer
	for i := 0; i < URL_LEN; i++ {
		urlBuffer.WriteString(componentArray[i])
	}
	return urlBuffer.String()
}

func request(url string) SearchResponse {
	var bodyString string
	var res SearchResponse

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("APPLICATION HTTP ERROR: ")
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println(err2)
		}
		bodyString = string(bodyBytes)
	} else {
		bodyString = resp.Status
	}

	// parse JSON https://eager.io/blog/go-and-json/
	data := []byte(bodyString)
	json.Unmarshal(data, &res)

	return res
}

func getFoodItem() string {
	foods := [8]string{"chicken", "beef", "taco", "crockpot", "pasta", "quick", "easy", "casserole"}
	seeded := rand.New(rand.NewSource(time.Now().UnixNano())) //Seeded Random num https://gobyexample.com/random-numbers
	return foods[seeded.Intn(len(foods))]
}

func getRandomRecipe() SearchResponse {
	var url string = createURL([URL_LEN]string{HOST, PORT,
		"/search?q=", getFoodItem(), "&app_id=", os.Getenv(AUTH_ID),
		"&app_key=", os.Getenv(AUTH_KEY), "&from=1&to=20"})
	fmt.Println("Querying " + url)
	res := request(url)
	return res
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	menu := getRandomRecipe()
	
	tmpl := template.New("menu")
	tmpl, err := tmpl.ParseFiles("templates/menu.html")
	err1 := tmpl.Execute(w, menu)// this doesn't work unless templates dir is in same dir as app binary

	if err != nil || err1 != nil {
		fmt.Println("error: ", err1, err)
	}

}

func main() {
	fmt.Println("Launching server on: " + LOCAL_HOST)
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(LOCAL_HOST, nil)
}
