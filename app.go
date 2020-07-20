package main

import (
  "fmt"
  "net/http"
  "net/url"
  "encoding/json"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome go chi\n"))
  })

  r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
    values, _ := url.ParseQuery(r.URL.RawQuery)
    // id, vkey, playid
    id := values["id"][0]
    vkey := values["vkey"][0]
    playid := values["playid"][0]
    w.Write([]byte(fmt.Sprintf("id: %v, vkey: %v, playid: %v, url: %v\n", id, vkey, playid, r.URL.RawQuery)))
  })

  r.Post("/upload.json", func(w http.ResponseWriter, r *http.Request) {
    // body_str, _ := ioutil.ReadAll(r.Body)
    type Ysp struct {
      Vkey string
      Livepid string
      Playurl string
    }
    var ysp Ysp
    _ = json.NewDecoder(r.Body).Decode(&ysp)
    w.Write([]byte(fmt.Sprintf("vkey: %v, livepid: %v, playurl: %v\n", ysp.Vkey, ysp.Livepid, ysp.Playurl)))
  })
  http.ListenAndServe(":3000", r)
}
