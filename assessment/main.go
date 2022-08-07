package main

import (
	"azarc/models"
	"azarc/service"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	filePath := flag.String("filePath", "", "a string")
	titleType := flag.String("titleType", "", "a string")
	primaryTitle := flag.String("primaryTitle", "", "a string")
	originalTitle := flag.String("originalTitle", "", "a string")
	startYear := flag.String("startYear", "", "an int")
	endYear := flag.String("endYear", "", "an int")
	runtimeMinutes := flag.String("runtimeMinutes", "", "an int")
	genres := flag.String("genres", "", "a string")
	maxApiRequests := flag.Int("maxApiRequests", 0, "an int")
	maxRunTime := flag.Int("maxRunTime", 0, "an int")

	flag.Parse()

	var filters = models.Filters{
		FilePath:       *filePath,
		TitleType:      *titleType,
		PrimaryTitle:   *primaryTitle,
		OriginalTitle:  *originalTitle,
		StartYear:      *startYear,
		EndYear:        *endYear,
		RuntimeMinutes: *runtimeMinutes,
		Genres:         *genres,
		MaxApiRequests: *maxApiRequests,
		MaxRunTime:     *maxRunTime,
	}

	if *filePath == "" {
		log.Fatalln("filePath cannot be empty")
	}

	if *primaryTitle == "" {
		log.Fatalln("primaryTitle cannot be empty")
	}

	imdbService := service.NewImdb(filters)
	resultCode, err := imdbService.Uncompress()
	if resultCode == -1 {
		log.Fatalln("Error: ", err)
	}

	var mu sync.Mutex
	dataChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go imdbService.ReadFile(&mu, &wg, dataChannel)
	go imdbService.ProcessData(dataChannel, &wg)

	wg.Wait()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
