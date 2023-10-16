package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MappingFunc[T any] func(row *goquery.Selection) T

type Armor struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ReqLevel int    `json:"reqLevel"`
	Def      int    `json:"def"`
	Mdef     int    `json:"mdef"`
	Effect   string `json:"effect"`
}

type Weapon struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ReqLevel    int    `json:"reqLevel"`
	Atk         int    `json:"atk"`
	Matk        int    `json:"matk"`
	WeaponLevel int    `json:"weaponLevel"`
	Effect      string `json:"effect"`
}

func ToInt(s string) int {
	int, err := strconv.Atoi(strings.Trim(s, "\n"))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return int
}

func mapToGarment(row *goquery.Selection) *Armor {
	name := strings.Trim(row.Find("td:nth-child(1)").Text(), "\n")
	reqLevel := ToInt(row.Find("td:nth-child(2)").Text())
	def := ToInt(row.Find("td:nth-child(3)").Text())
	mdef := ToInt(row.Find("td:nth-child(4)").Text())
	effect := strings.Trim(row.Find("td:nth-child(5)").Text(), "\n")

	return &Armor{
		Name:     name,
		Type:     "Garment",
		ReqLevel: reqLevel,
		Def:      def,
		Mdef:     mdef,
		Effect:   effect,
	}
}

func main() {
	url := "https://returntomorroc.com/wiki/index.php/Garment"
	garments := scrape(url, mapToGarment)
	saveAsJson("garments.json", garments)
}

func saveAsJson[T any](name string, items []T) {
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(name, jsonBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The struct has been saved to the file '%s'.", name)
}

func scrape[T any](url string, mapFunc MappingFunc[T]) []T {
	var results []T
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	table := doc.Find("table").First()

	table.Find("tr").Each(func(_ int, rowHtml *goquery.Selection) {
		item := mapFunc(rowHtml)
		results = append(results, item)
	})

	return results
}
