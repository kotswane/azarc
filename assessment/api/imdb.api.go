package api

import (
	"azarc/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ClientAPI struct {
	httpClient *http.Client
}

func NewImdbAPI() ImdbAPI {
	httpClient := &http.Client{Timeout: time.Duration(15) * time.Second}
	return &ClientAPI{
		httpClient: httpClient,
	}
}

func (c ClientAPI) QueryIMDB(title string) (models.Plot, error) {
	var body []byte
	// 
	var uri = fmt.Sprintf("http://www.omdbapi.com/?t=%s&apikey=8bc41d20", title)
	request, err := http.NewRequest(http.MethodGet, uri, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var plot models.Plot
	json.Unmarshal(payload, &plot)
	return plot, nil
}
