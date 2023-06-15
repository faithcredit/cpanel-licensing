package utils

type Response struct {
	Status      int    `json:"status"`
	Expiry      string `json:"expiry" `
	Brand       string `json:"brand"`
	Url         string `bson:"url"`
	Ip          string `json:"ip"`
	Version     string `json:"version"`
	CurrentDate string `json:"currentDate"`
	StatusCode  int
}
