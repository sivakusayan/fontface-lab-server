package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Font struct {
	Family     string
	Category   string
	Subsets    []string
	Popularity int
}

type FontFamilyList struct {
	FamilyMetadataList []Font
}

var c = CreateCache()

func getCachedFontFamilyList(w http.ResponseWriter, r *http.Request) {
	data, found := c.Get("FontFamilyList")

	if !found {
		list := getFontFamilyList()
		jsonList, _ := json.Marshal(list)
		data = string(jsonList)

		c.Set("FontFamilyList", data, 24*time.Hour)
	}

	io.WriteString(w, data)
}

func getFontFamilyList() FontFamilyList {
	response, err := http.Get("https://fonts.google.com/metadata/fonts")
	if err != nil {
		log.Fatal(err)
	}
	byteStream, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newData FontFamilyList
	json.Unmarshal(byteStream, &newData)

	for fontIndex := range newData.FamilyMetadataList {
		subsets := newData.FamilyMetadataList[fontIndex].Subsets
		menuSubsetIndex := IndexOf(subsets, "menu")
		if menuSubsetIndex > -1 {
			Remove(subsets, menuSubsetIndex)
		}
	}

	return newData
}

func main() {
	godotenv.Load()
	http.HandleFunc("/api/font-family-list", getCachedFontFamilyList)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
