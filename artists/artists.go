package artists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type SearchAjaxArtists struct {
	Error               string     `json:"error"`
	TotalRecords        int        `json:"iTotalRecords"`
	TotalDisplayRecords int        `json:"iTotalDisplayRecords"`
	Echo                int        `json:"sEcho"`
	Data                [][]string `json:"aaData"`
}

type SearchArtistsData struct {
	Name string
	URL  string
}

func searchArtistAjax(client http.Client, artist string) ([][]string, error) {

	var searchArtistData [][]string
	artistString := strings.Replace(artist, " ", "+", -1)
	url := fmt.Sprintf("https://www.metal-archives.com/search/ajax-band-search/?field=name&query=%s", artistString)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return searchArtistData, err
	}

	req.Header.Set("User-Agent", "https://github.com/a-castellano/metal-archives-wrapper")

	res, getErr := client.Do(req)
	if getErr != nil {
		return searchArtistData, getErr
	}
	//fmt.Println("_____")
	//fmt.Println(res)
	//fmt.Println("_____")

	body, readErr := ioutil.ReadAll(res.Body)
	//fmt.Println("_____")
	//s := string(body)
	//fmt.Println(s) // ABC€�
	//fmt.Println("_____")
	if readErr != nil {
		return searchArtistData, readErr
	}
	searchArtist := SearchAjaxArtists{}
	jsonErr := json.Unmarshal(body, &searchArtist)
	if jsonErr != nil {
		return searchArtistData, jsonErr
	}
	searchArtistData = searchArtist.Data
	return searchArtistData, nil
}

func SearchArtist(artist string) (SearchArtistsData, error) {

	var artistData SearchArtistsData
	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	data, err := searchArtistAjax(client, artist)

	if err != nil {
		return artistData, err
	} else {
		fmt.Println(data)
		re := regexp.MustCompile(`^<a href=\"([^\"]+)\">([^<]+)</a>`)
		for _, foundArtistData := range data {
			fmt.Println(foundArtistData)
			match := re.FindAllStringSubmatch(foundArtistData[0], -1)
			fmt.Println("_")
			fmt.Println(match[0][1])
			fmt.Println(match[0][2])
		}
	}
	return artistData, nil
}