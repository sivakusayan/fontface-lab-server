package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/data"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	http.HandleFunc("/api/font-family-list", onFontFamilyListRequest)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func onFontFamilyListRequest(w http.ResponseWriter, r *http.Request) {
	fontFamilyList := data.GetCachedFontFamilyList()
	json, _ := json.Marshal(fontFamilyList)

	// Cache for 6 hours
	w.Header().Set("Cache-Control", "public, max-age=21600")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(json))
}
