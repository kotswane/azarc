package service

import (
	"azarc/api"
	"azarc/models"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

type imdbService struct {
	uncompressedFile string
	filters          models.Filters
	mu               sync.Mutex
}

func NewImdb(filters models.Filters) ImdbService {
	return &imdbService{filters: filters}
}

func (i *imdbService) ReadFile(mu *sync.Mutex, wg *sync.WaitGroup, ch chan<- string) {

	fmt.Println("Filters: ", i.filters)
	file, err := os.Open(i.uncompressedFile)

	if err != nil {
		file.Close()
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	fmt.Println("loading file ", i.filters.FilePath)
	for scanner.Scan() {
		mu.Lock()
		ch <- scanner.Text()
		mu.Unlock()
	}
	defer file.Close()
	wg.Done()
}

func (i *imdbService) ProcessData(data <-chan string, wg *sync.WaitGroup) {
	count := 0
	apiCallCount := 1
	for line := range data {

		if count > 0 {
			dataFields := i.formatData(line)
			correctFilters := i.filterData(dataFields)
			if correctFilters {

				if i.filters.MaxApiRequests > 0 {
					if apiCallCount <= i.filters.MaxApiRequests {
						fmt.Println("Limited API calls: ")
						response := i.doImdbApiCall(dataFields)
						apiCallCount++
						if response.IMDB_ID != "" {
							fmt.Println("API response: ", response)
						}
					}
					if apiCallCount > i.filters.MaxApiRequests {
						os.Exit(0)
					}
				} else {
					fmt.Println("Unlimited API calls: ")
					response := i.doImdbApiCall(dataFields)
					if response.IMDB_ID != "" {
						fmt.Println("API response: ", response)
					}
				}

			}
		}
		count++
	}
	wg.Done()
}

func (i *imdbService) filterData(data models.DataFields) bool {

	primaryTitle, _ := regexp.Compile(i.filters.PrimaryTitle)
	originalTitle, _ := regexp.Compile(i.filters.OriginalTitle)
	titleType, _ := regexp.Compile(i.filters.TitleType)
	genres, _ := regexp.Compile(i.filters.Genres)
	runtimeMinutes, _ := regexp.Compile(i.filters.RuntimeMinutes)
	endYear, _ := regexp.Compile(i.filters.EndYear)
	startYear, _ := regexp.Compile(i.filters.StartYear)

	if !primaryTitle.MatchString(data.PrimaryTitle) && i.filters.PrimaryTitle != "" {
		return false
	}

	if !originalTitle.MatchString(data.OriginalTitle) && i.filters.OriginalTitle != "" {
		return false
	}

	if !titleType.MatchString(data.TitleType) && i.filters.TitleType != "" {
		return false
	}

	if !genres.MatchString(data.Genres) && i.filters.Genres != "" {
		return false
	}

	if !runtimeMinutes.MatchString(data.RuntimeMinutes) && i.filters.RuntimeMinutes != "" {
		return false
	}

	if !endYear.MatchString(data.EndYear) && i.filters.EndYear != "" {
		return false
	}

	if !startYear.MatchString(data.StartYear) && i.filters.StartYear != "" {
		return false
	}
	return true
}

func (i *imdbService) formatData(data string) models.DataFields {
	var dataFields models.DataFields
	dataSlice := strings.Split(data, "\t")
	dataFields.Tconst = dataSlice[0]
	dataFields.TitleType = dataSlice[1]
	dataFields.PrimaryTitle = dataSlice[2]
	dataFields.OriginalTitle = dataSlice[3]
	dataFields.IsAdult = dataSlice[4]
	dataFields.StartYear = dataSlice[5]
	dataFields.EndYear = dataSlice[6]
	dataFields.RuntimeMinutes = dataSlice[6]
	dataFields.Genres = dataSlice[7]
	return dataFields
}

func (i *imdbService) Uncompress() (int64, error) {

	if _, err := os.Stat(i.filters.FilePath); errors.Is(err, os.ErrNotExist) {
		return -1, err
	}
	fmt.Println("uncompressing file ", i.filters.FilePath)
	gzipFile, err := os.Open(i.filters.FilePath)
	if err != nil {
		return -1, err
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return -1, err
	}
	defer gzipReader.Close()

	newFile := strings.Replace("gz", i.filters.FilePath, "", 1)
	outfileWriter, err := os.Create(newFile)
	if err != nil {
		return -1, err
	}
	defer outfileWriter.Close()
	i.uncompressedFile = newFile

	response, err := io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return -1, err
	}
	return response, err
}

func (i *imdbService) doImdbApiCall(data models.DataFields) models.Plot {

	response, err := api.NewImdbAPI().QueryIMDB(data.PrimaryTitle)
	if err != nil {
		fmt.Println(err)
	}
	return response
}
