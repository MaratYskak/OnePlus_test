package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jszwec/csvutil"
)

type InstagramData struct {
	Rank       string `csv:"Rank"`
	Influencer string `csv:"Influencer"`
	Category   string `csv:"Category"`
	Followers  string `csv:"Followers"`
	Country    string `csv:"Country"`
	Authentic  string `csv:"Eng. (Auth.)"`
	AvgEng     string `csv:"Eng. (Avg.)"`
}

func main() {
	url := "https://hypeauditor.com/top-instagram-all-russia/"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("instagram_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data []InstagramData

	doc.Find("div.table div.row").Each(func(index int, row *goquery.Selection) {
		rank := row.Find("div.row-cell.rank span").Text()
		influencer := row.Find("div.row-cell.contributor div.contributor-wrap a.contributor div.contributor__title").Text()
		categoryElements := row.Find("div.row-cell.category div.tag div.tag__content")
		categories := []string{}
		categoryElements.Each(func(index int, category *goquery.Selection) {
			categories = append(categories, category.Text())
		})
		followers := row.Find("div.row-cell.subscribers").Text()
		country := row.Find("div.row-cell.audience").Text()
		authentic := row.Find("div.row-cell.authentic").Text()
		engagement := row.Find("div.row-cell.engagement").Text()

		data = append(data, InstagramData{
			Rank:       rank,
			Influencer: influencer,
			Category:   strings.Join(categories, ","),
			Followers:  followers,
			Country:    country,
			Authentic:  authentic,
			AvgEng:     engagement,
		})
	})

	csvData, err := csvutil.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(csvData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные успешно спарсены и сохранены в файл instagram_data.csv")
}
