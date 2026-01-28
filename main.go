package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>running...</title>
</head>
<body>
  <h1>running...</h1>
</body>
</html>`)
	})

	log.Println("Server started: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
