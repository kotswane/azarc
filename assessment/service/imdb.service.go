package service

import (
	"azarc/models"
	"sync"
)

type ImdbService interface {
	Uncompress() (int64, error)
	ReadFile(mu *sync.Mutex, wg *sync.WaitGroup, ch chan<- string)
	filterData(data models.DataFields) bool
	formatData(data string) models.DataFields
	doImdbApiCall(data models.DataFields) models.Plot
	ProcessData(ch <-chan string, wg *sync.WaitGroup)
}
