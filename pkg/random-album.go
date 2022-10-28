package random_album

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/pedrocb/random-album-picker/internal"

	"github.com/PuerkitoBio/goquery"
)

func buildCDFFromRatings(ratings map[int]int, multiplier float64) []float64 {
	// Build CDF from ratings map

	// Calculate total number of ratings and shares
	totalRatings := 0
	totalShares := 0.0
	for rating, numRatings := range ratings {
		totalRatings += numRatings
		totalShares += (math.Pow(multiplier, float64(rating-1)) * float64(numRatings))
	}

	// Initialize cumulative distribution function array
	cdf := make([]float64, totalRatings)

	// First allocate shares for each album so that each rating has <multiplier> times
	// the shares the rating before (this means that an album that has 10*
	// is twice more likely to be picked than another one with 9* if multiplier is 2)
	counter := 0
	for rating := 1; rating <= 10; rating += 1 {
		numRatings := ratings[rating]
		numShares := math.Pow(multiplier, float64(rating-1))
		percentage := numShares / totalShares
		for i := 0; i < numRatings; i++ {
			if counter == 0 {
				cdf[counter] = percentage
			} else {
				cdf[counter] = cdf[counter-1] + percentage
			}
			counter += 1
		}
	}
	return cdf
}

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

func getRandomRatingIndex(ratings map[int]int, multiplier float64) int {
	cdf := buildCDFFromRatings(ratings, multiplier)

	index := internal.GetSampleFromCDF(cdf)
	return index
}

func getAlbumByIndex(user string, index int) (string, string) {
	page := (index / 25) + 1
	albumIndex := index % 25

	// Request RYM profile page
	client := &http.Client{}
	url := fmt.Sprintf("https://rateyourmusic.com/collection/%s/r0.5-5.0,ss.r/%d", user, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", ""
	}
	req.Header.Set("User-Agent", "RYM Go Scraper")
	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", ""
	}
	album := doc.Find(".or_q_albumartist_td").Eq(albumIndex).Text()
	rating, _ := doc.Find(".or_q_rating_date_s img").Eq(albumIndex).Attr("title")

	return album, rating
}

func GetRandomAlbumFromRYM(user string) (string, error) {
	// Get user ratings
	myUserRatings, err := getUserRatings(user)
	if err != nil {
		return "", err
	}
	album, rating := getAlbumByIndex(user, getRandomRatingIndex(myUserRatings, 2.0))
	return fmt.Sprintf("%s - %s", album, rating), nil
}

