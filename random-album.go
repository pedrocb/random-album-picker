package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"random-album-picker/internal"

	"github.com/PuerkitoBio/goquery"
)

func getUserRatings(user string) (map[int]int, error) {
	// Request RYM profile page
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://rateyourmusic.com/~%s", user), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "RYM Python Scraper")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Build ratings map {<rating>: <numRatings>}
	ratings := make(map[int]int)
	doc.Find("#musicrating tr").Each(func(j int, tr *goquery.Selection) {
		if j == 0 || j == 11 {
			return
		}
		childs := tr.Find("td")
		rating, err := strconv.ParseFloat(strings.TrimSpace(childs.First().Text()), 64)
		if err != nil {
			fmt.Printf("Got error %s", err.Error())
			return
		}

		numRatings, _ := strconv.ParseFloat(strings.TrimSpace(childs.Next().Text()), 64)

		ratings[int(rating*2)] = int(numRatings)
	})

	return ratings, nil
}

func getRandomRatingIndex(ratings map[int]int) int {
	// Calculate total number of ratings
	totalRatings := 0
	for _, numRatings := range ratings {
		totalRatings += numRatings
	}

	// Initialize array of percentages (each value in the array is the
	// percentage of the album on that index to be picked)
	percentages := make([]float64, totalRatings)

	// First allocate shares for each album so that each rating has double
	// the shares the rating before (this means that an album that has 10*
	// is twice more likely to be picked than another one with 9*)
	counter := 0
	sumPercentages := 0.0
	for rating := 1; rating <= 10; rating += 1 {
		numRatings := ratings[rating]
		for i := 0; i < numRatings; i++ {
			_percentage := math.Pow(2, float64(rating-1))
			percentages[counter] = _percentage
			sumPercentages += _percentage
			counter += 1
		}
	}

	percentages[0] /= sumPercentages
	for index, value := range percentages {
		if index == 0 {
			continue
		}
		percentages[index] = percentages[index-1] + value/sumPercentages
	}

	index := internal.GetSampleFromCDF(percentages)
	return index
}

func getAlbumByIndex(user string, index int) string {
	page := (index / 25) + 1
	albumIndex := index % 25

	// Request RYM profile page
	client := &http.Client{}
	url := fmt.Sprintf("https://rateyourmusic.com/collection/%s/r0.5-5.0,ss.r/%d", user, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "RYM Go Scraper")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ""
	}
	album := doc.Find(".or_q_albumartist_td").Eq(albumIndex).Text()

	return album
}

func main() {
	user := "pedrocb"
	// Get user ratings
	myUserRatings, err := getUserRatings(user)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
		return
	}
	fmt.Println(getAlbumByIndex(user, getRandomRatingIndex(myUserRatings)))
}
