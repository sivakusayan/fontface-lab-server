package data

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

var c = CreateCache()

type Font struct {
	Family     string
	Category   string
	Subsets    []string
	Popularity int
}

type FontFamilyList struct {
	FamilyMetadataList []Font
}

func GetCachedFontFamilyList() FontFamilyList {
	data, found := c.Get("FontFamilyList")

	if !found {
		data = getFontFamilyList()
		c.Set("FontFamilyList", data, 24*time.Hour)
	}

	return data.(FontFamilyList)
}

// Implementation of sort.Interface so we can sort by popularity.
// https://pkg.go.dev/sort#Interface
type ByPopularity []Font

func (a ByPopularity) Len() int           { return len(a) }
func (a ByPopularity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPopularity) Less(i, j int) bool { return a[i].Popularity < a[j].Popularity }

func getFontFamilyList() FontFamilyList {
	response, err := http.Get("https://fonts.google.com/metadata/fonts")
	if err != nil {
		log.Fatal(err)
	}
	byteStream, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data FontFamilyList
	json.Unmarshal(byteStream, &data)

	// The Google Fonts API seems to list a large amount of fonts under the
	// "menu" subset. I'm not sure what it is, but we don't need it, so filter
	// it out to decrease payload size.
	for i := range data.FamilyMetadataList {
		subsets := data.FamilyMetadataList[i].Subsets
		menuSubsetIndex := indexOf(subsets, "menu")
		if menuSubsetIndex > -1 {
			remove(subsets, menuSubsetIndex)
		}
	}

	sort.Sort(ByPopularity(data.FamilyMetadataList))

	return data
}
