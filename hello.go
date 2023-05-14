package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

type sumHandler struct{}

func parseQueryValues(urlStr string) (map[string][]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return u.Query(), nil
}

func (h *sumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	queryValues, err := parseQueryValues(r.URL.String())
	if err != nil {
		fmt.Println("error parsing query values")
		return
	}

	n1, err := strconv.Atoi(queryValues["n1"][0])
	if err != nil {
		fmt.Println("error parsing n1")
		return
	}

	n2, err := strconv.Atoi(queryValues["n2"][0])
	if err != nil {
		fmt.Println("error parsing n2")
		return
	}

	sum := n1 + n2
	fmt.Println("sum is ", sum)
	fmt.Fprintf(w, "sum is %d\n", sum)
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Println("count is ", h.n)
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	http.Handle("/count", &countHandler{})
	http.Handle("/sum", &sumHandler{})
	log.Fatal(http.ListenAndServe(":7001", nil))
}
