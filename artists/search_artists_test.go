// +build integration_tests unit_tests

package artists

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSearchArtistAjaxNoArtists(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 0,
	"iTotalDisplayRecords": 0,
	"sEcho": 0,
	"aaData": [
		]
}
	`))}}}

	data, err := searchArtistAjax(client, "AnyArtist")

	if err != nil {
		t.Errorf("TestClientNoArtists shouldn't fail.")
	}

	if len(data) != 0 {
		t.Errorf("TestClientNoArtists should return empty data.")
	}

}

func TestSearchArtistAjaxBrokenJson(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 0,
	"iTotalDisplayRecords": 0,
	"sEcho": 0,
	"aaData": [
}
	`))}}}

	_, err := searchArtistAjax(client, "AnyArtist")

	if err == nil {
		t.Errorf("TestBrokenJson should fail because JSON response is broken.")
	}
}

func TestSearchArtistAjaxOneArtist(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 1,
	"iTotalDisplayRecords": 1,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Satyricon/341\">Satyricon</a>  <!-- 12.348988 -->" ,
			"Black Metal" ,
			"Norway"     		]
				]
}
	`))}}}

	data, err := searchArtistAjax(client, "AnyArtist")

	if err != nil {
		t.Errorf("TestClientNoArtists shouldn't fail.")
	}

	if len(data) != 1 {
		t.Errorf("TestClientNoArtists should return one entry only.")
	}
}

func TestSearchArtistAjaxMoreThanOneArtist(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 3,
	"iTotalDisplayRecords": 3,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Burzum/88\">Burzum</a>  <!-- 11.432714 -->" ,
			"Black Metal, Ambient" ,
			"Norway"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Down_to_Burzum/3540435931\">Down to Burzum</a>  <!-- 5.716357 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Krimparturr/21151\">Krimparturr</a> (<strong>a.k.a.</strong> Krimpartûrr Bürzum Shi-Hai) <!-- 1.2505064 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				]
}
	`))}}}

	data, err := searchArtistAjax(client, "AnyArtist")

	if err != nil {
		t.Errorf("TestSearchArtistAjaxMoreThanOneArtist shouldn't fail.")
	}

	if len(data) != 3 {
		t.Errorf("TestSearchArtistAjaxMoreThanOneArtist should return three entries.")
	}
}

func TestSearchArtistErrored(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 0,
	"iTotalDisplayRecords": 0,
	"sEcho": 0,
	"aaData": [
}
	`))}}}

	data, extraData, err := SearchArtist(client, "AnyArtist")

	if err == nil {
		t.Errorf("TestSearchArtistAjaxMoreThanOneArtist should fail.")
	}

	if data.Name != "" {
		t.Errorf("Retrieved artist name should be empty.")
	}

	if data.URL != "" {
		t.Errorf("Retrieved artist URL should be empty.")
	}

	if data.ID != "" {
		t.Errorf("Retrieved artist id should be empty.")
	}

	if data.Genre != "" {
		t.Errorf("Retrieved artist Genre should be empty.")
	}

	if data.Country != "" {
		t.Errorf("Retrieved artist Country should be empty.")
	}

	if len(extraData) != 0 {
		t.Errorf("Retrieved extra data should be empty.")
	}
}

func TestSearchArtistNotFound(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 0,
	"iTotalDisplayRecords": 0,
	"sEcho": 0,
	"aaData": [
		]
}
	`))}}}

	data, extraData, err := SearchArtist(client, "AnyArtist")

	if err == nil {
		t.Errorf("TestSearchArtistNotFound should fail.")
	}

	if err.Error() != "No artist was found." {
		t.Errorf("TestSearchArtistNotFound error should be 'No artist was found.'")
	}

	if data.Name != "" {
		t.Errorf("Retrieved artist name should be empty.")
	}

	if data.URL != "" {
		t.Errorf("Retrieved artist URL should be empty.")
	}

	if data.ID != "" {
		t.Errorf("Retrieved artist id should be empty.")
	}

	if data.Genre != "" {
		t.Errorf("Retrieved artist Genre should be empty.")
	}

	if data.Country != "" {
		t.Errorf("Retrieved artist Country should be empty.")
	}

	if len(extraData) != 0 {
		t.Errorf("Retrieved extra data should be empty.")
	}
}

func TestSearchArtistNotMatch(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 3,
	"iTotalDisplayRecords": 3,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Burzum/88\">Burzum</a>  <!-- 11.432714 -->" ,
			"Black Metal, Ambient" ,
			"Norway"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Down_to_Burzum/3540435931\">Down to Burzum</a>  <!-- 5.716357 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Krimparturr/21151\">Krimparturr</a> (<strong>a.k.a.</strong> Krimpartûrr Bürzum Shi-Hai) <!-- 1.2505064 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				]
}

	`))}}}

	data, extraData, err := SearchArtist(client, "AnyArtist")

	if err == nil {
		t.Errorf("TestSearchArtistNotMatch should fail.")
	}

	if err.Error() != "No artist was found." {
		t.Errorf("TestSearchArtistNotMatch error should be 'No artist was found.'")
	}

	if data.Name != "" {
		t.Errorf("Retrieved artist name should be empty.")
	}

	if data.URL != "" {
		t.Errorf("Retrieved artist URL should be empty.")
	}

	if data.ID != "" {
		t.Errorf("Retrieved artist id should be empty.")
	}

	if data.Genre != "" {
		t.Errorf("Retrieved artist Genre should be empty.")
	}

	if data.Country != "" {
		t.Errorf("Retrieved artist Country should be empty.")
	}

	if len(extraData) != 0 {
		t.Errorf("Retrieved extra data should be empty.")
	}

}

func TestSearchArtistMatch(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 3,
	"iTotalDisplayRecords": 3,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Burzum/88\">Burzum</a>  <!-- 11.432714 -->" ,
			"Black Metal, Ambient" ,
			"Norway"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Down_to_Burzum/3540435931\">Down to Burzum</a>  <!-- 5.716357 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Krimparturr/21151\">Krimparturr</a> (<strong>a.k.a.</strong> Krimpartûrr Bürzum Shi-Hai) <!-- 1.2505064 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				]
}

	`))}}}

	data, extraData, err := SearchArtist(client, "Burzum")

	if err != nil {
		t.Errorf("TestSearchArtistMatch shouldn't fail.")
	}

	if data.Name != "Burzum" {
		t.Errorf("Retrieved artist name should be 'Burzum'.")
	}

	if data.URL != "https://www.metal-archives.com/bands/Burzum/88" {
		t.Errorf("Retrieved artist URL should be 'https://www.metal-archives.com/bands/Burzum/88'.")
	}

	if data.ID != "88" {
		t.Errorf("Retrieved artist id should be 88.")
	}

	if data.Genre != "Black Metal, Ambient" {
		t.Errorf("Retrieved artist Genre should be 'Black Metal, Ambient'.")
	}

	if data.Country != "Norway" {
		t.Errorf("Retrieved artist Country should be 'Norway'.")
	}

	if len(extraData) != 0 {
		t.Errorf("Retrieved extra data should be empty.")
	}
}

func TestSearchArtistMatchLowercase(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 3,
	"iTotalDisplayRecords": 3,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Burzum/88\">Burzum</a>  <!-- 11.432714 -->" ,
			"Black Metal, Ambient" ,
			"Norway"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Down_to_Burzum/3540435931\">Down to Burzum</a>  <!-- 5.716357 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Krimparturr/21151\">Krimparturr</a> (<strong>a.k.a.</strong> Krimpartûrr Bürzum Shi-Hai) <!-- 1.2505064 -->" ,
			"Black Metal" ,
			"Brazil"     		]
				]
}

	`))}}}

	data, extraData, err := SearchArtist(client, "burzum")

	if err != nil {
		t.Errorf("TestSearchArtistMatchLowercase shouldn't fail.")
	}

	if data.Name != "Burzum" {
		t.Errorf("Retrieved artist name should be 'Burzum'.")
	}

	if data.URL != "https://www.metal-archives.com/bands/Burzum/88" {
		t.Errorf("Retrieved artist URL should be 'https://www.metal-archives.com/bands/Burzum/88'.")
	}

	if data.ID != "88" {
		t.Errorf("Retrieved artist id should be 88.")
	}

	if data.Genre != "Black Metal, Ambient" {
		t.Errorf("Retrieved artist Genre should be 'Black Metal, Ambient'.")
	}

	if data.Country != "Norway" {
		t.Errorf("Retrieved artist Country should be 'Norway'.")
	}

	if len(extraData) != 0 {
		t.Errorf("Retrieved extra data should be empty.")
	}

}

func TestSearchArtistMultipleMatches(t *testing.T) {
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 5,
	"iTotalDisplayRecords": 5,
	"sEcho": 0,
	"aaData": [
				[
			"<a href=\"https://www.metal-archives.com/bands/Hypocrisy/96\">Hypocrisy</a>  <!-- 10.740315 -->" ,
			"Death Metal (early), Melodic Death Metal (later)" ,
			"Sweden"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Hypocrisy/56165\">Hypocrisy</a>  <!-- 10.740315 -->" ,
			"Power/Thrash Metal" ,
			"United States"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Sermon_of_Hypocrisy/7033\">Sermon of Hypocrisy</a>  <!-- 5.3701577 -->" ,
			"Black Metal" ,
			"United Kingdom"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/The_Polo_Hypocrisy/47897\">The Polo Hypocrisy</a> (<strong>a.k.a.</strong> T.P.H.) <!-- 5.3701577 -->" ,
			"Melodic Death Metal with Hardcore elements" ,
			"Canada"     		]
				,
						[
			"<a href=\"https://www.metal-archives.com/bands/Torture_of_Hypocrisy/3540316100\">Torture of Hypocrisy</a> (<strong>a.k.a.</strong> ToH) <!-- 5.3701577 -->" ,
			"Technical Thrash Metal" ,
			"Poland"     		]
				]
}
	`))}}}

	data, extraData, err := SearchArtist(client, "Hypocrisy")

	if err != nil {
		t.Errorf("TestSearchArtistMultipleMatches shouldn't fail.")
	}

	if data.Name != "Hypocrisy" {
		t.Errorf("Retrieved artist name should be 'Hypocrisy'.")
	}

	if data.URL != "https://www.metal-archives.com/bands/Hypocrisy/96" {
		t.Errorf("Retrieved artist URL should be 'https://www.metal-archives.com/bands/Hypocrisy/96'.")
	}

	if data.ID != "96" {
		t.Errorf("Retrieved artist id should be 96.")
	}

	if data.Genre != "Death Metal (early), Melodic Death Metal (later)" {
		t.Errorf("Retrieved artist Genre should be 'Death Metal (early), Melodic Death Metal (later)'.")
	}

	if data.Country != "Sweden" {
		t.Errorf("Retrieved artist Country should be 'Sweden'.")
	}

	if len(extraData) != 1 {
		t.Errorf("Retrieved extra data should have one extra item only.")
	}

	if extraData[0].Name != "Hypocrisy" {
		t.Errorf("Retrieved artist name should be 'Hypocrisy'.")
	}

	if extraData[0].URL != "https://www.metal-archives.com/bands/Hypocrisy/56165" {
		t.Errorf("Retrieved artist URL should be 'https://www.metal-archives.com/bands/Hypocrisy/56165'.")
	}

	if extraData[0].ID != "56165" {
		t.Errorf("Retrieved artist id should be 56165.")
	}

	if extraData[0].Genre != "Power/Thrash Metal" {
		t.Errorf("Retrieved artist Genre should be 'Power/Thrash Metal'.")
	}

	if extraData[0].Country != "United States" {
		t.Errorf("Retrieved artist Country should be 'United States'.")
	}

}
