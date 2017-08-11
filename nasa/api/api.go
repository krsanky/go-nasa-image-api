package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/krsanky/lg"
)

var base_url = "https://images-api.nasa.gov"

type Wrapper struct {
	Collection CollectionType `json:"collection"`
}

type CollectionType struct {
	Href     string                 `json:"href"`
	Version  string                 `json:"version"`
	Metadata map[string]interface{} `json:"metadata"`
	Items    []Item                 `json:"items"`
	Error    interface{}            `json:"error"`
}

type Item struct {
	Href  string     `json:"href"`
	Links []Link     `json:"links"`
	Data  []DataType `json:"data"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Render string `json:"render"`
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
}

type DataType struct {
	MediaType   string `json:"media_type"`
	NasaId      string `json:"nasa_id"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Center      string `json:"center"`
	DateCreated string `json:"date_created"`
}

func getJson() []byte {
	dat, err := ioutil.ReadFile("file1.json")
	if err != nil {
		panic(err)
	}

	return dat
}

func GetSearch(values url.Values) ([]byte, error) {

	media_type := url.QueryEscape(values.Get("media_type"))
	search_string := url.QueryEscape(values.Get("search_string"))
	description := url.QueryEscape(values.Get("description"))
	keywords := url.QueryEscape(values.Get("keywords"))
	url1 := fmt.Sprintf("%s/search?media_type=%s", base_url, media_type)

	gt0 := false
	if search_string != "" {
		gt0 = true
		url1 = fmt.Sprintf("%s&q=%s", url1, search_string)
	}
	if description != "" {
		gt0 = true
		url1 = fmt.Sprintf("%s&description=%s", url1, description)
	}
	if keywords != "" {
		gt0 = true
		url1 = fmt.Sprintf("%s&keywords=%s", url1, keywords)
	}
	lg.Log.Printf("url1:%s", url1)

	var body []byte
	if !gt0 {
		return []byte{}, errors.New("supply at least one search field")
	} else {
		var err error
		res, err := http.Get(url1)
		if err != nil {
			return []byte{}, err
		}
		body, err = ioutil.ReadAll(res.Body)

		res.Body.Close()
		if err != nil {
			return []byte{}, err
		}
	}
	//lg.Log.Println(string(body))
	return body, nil
}

func GetDetail(nasa_id string) ([]byte, error) {
	var err error
	url_ := fmt.Sprintf("%s/asset/%s", base_url, nasa_id)
	res, err := http.Get(url_)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return body, err
}

func JsonTest(bs []byte) {
	var wrapper Wrapper
	err := json.Unmarshal(bs, &wrapper)
	if err != nil {
		panic(err)
	}
	lg.Log.Printf("wrapper.collection.version:%s\n", wrapper.Collection.Version)
	lg.Log.Printf("wrapper.collection.version.metadata[total_hits]:%g\n", wrapper.Collection.Metadata["total_hits"])
	lg.Log.Printf("wrapper.collection.href:%s\n", wrapper.Collection.Href)
	//lg.Log.Printf("wrapper.collection.items:%v\n", wrapper.Collection.Items)
	lg.Log.Printf("wrapper.collection.items[0].data[0].media_type:%s\n", wrapper.Collection.Items[0].Data[0].MediaType)
}
