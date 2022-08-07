package models

type Filters struct {
	FilePath       string
	TitleType      string
	PrimaryTitle   string
	OriginalTitle  string
	StartYear      string
	EndYear        string
	RuntimeMinutes string
	Genres         string
	MaxApiRequests int
	MaxRunTime     int
	PlotFilter     string
}
