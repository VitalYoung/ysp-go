package main

import (
  "os"
  "fmt"
  "log"
  "bufio"
  "context"
  "regexp"
  "net/http"
  "net/url"
  "encoding/json"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "github.com/go-redis/redis"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  ctx := context.Background()
  rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
  })

  r.Get("/ysp.m3u8", func(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("./ysp_videos.txt")
    if err != nil {
      log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    re, _ := regexp.Compile("tvg-id=\"([0-9]*)\"") // tvg-id="600004092"
    out_str := ""
    for scanner.Scan() {
      if re.MatchString(scanner.Text()) {
        livepid := re.FindStringSubmatch(scanner.Text())[1]
        // fmt.Println(livepid)
        val, err := rdb.Get(ctx, livepid).Result()
        if err == redis.Nil {
          fmt.Println("key [" + livepid + "] does not exist\n")
        } else if err != nil {
          fmt.Println(err)
        } else {
          out_str += scanner.Text() + "\n" + val + "\n"
        }
      }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    w.Header().Add("Content-Type", "application/x-mpegURL")
    w.Write([]byte(out_str))
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
    _ = rdb.Set(ctx, ysp.Livepid, ysp.Playurl, 0)
    w.Write([]byte(fmt.Sprintf("success\n")))
  })
  http.ListenAndServe(":3000", r)
}
