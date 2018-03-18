package main

import(
  "io"
  "net/http"

  "github.com/go-zoo/bone"
)

func main () {
  mux := bone.New()
  mux.Get("/hello", http.HandlerFunc(helloHandler))
  http.ListenAndServe(":8181", mux)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	str := `{"hello": "world"}`
	io.WriteString(w, str)
}
