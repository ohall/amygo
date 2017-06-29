package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "bytes"
    "math/rand"
    "encoding/json"
    "time"
)

const HOST = "https://api.edamam.com"
const PORT = ":443"
const AUTH_ID = "EDAMAM_ID"
const AUTH_KEY = "EDAMAM_KEY"
const URL_LEN = 9

type SearchResponse struct {
    Query string `json:"q"`
    Hits [] struct {
        Recipe struct {
            URI string `json:"uri"`
            Label string `json:"label"`
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

func getFoodItem() string  {
    foods := [8]string{"chicken", "beef", "taco", "crockpot", "pasta", "quick", "easy", "casserole"}
    seeded := rand.New(rand.NewSource(time.Now().UnixNano())) //Seeded Random num https://gobyexample.com/random-numbers
    return foods[seeded.Intn(len(foods))]
}

func main() {
    var url string = createURL([URL_LEN]string{HOST,PORT,
        "/search?q=",getFoodItem(),"&app_id=",os.Getenv(AUTH_ID),
        "&app_key=",os.Getenv(AUTH_KEY),"&from=1&to=20"})

    res := request(url)

    fmt.Println("Some recipies:")
    for i := 0; i < len(res.Hits); i++ {
        fmt.Println(res.Hits[i].Recipe.Label)
        fmt.Println(res.Hits[i].Recipe.URI)
        fmt.Println()
    }
}
