package api

type ELicznikMeasurement struct {
	Date  string  `json:"Date"`
	Hour  int     `json:"Hour,string"`
	Power float64 `json:"EC,string"`
}

type ELicznikData struct {
	Ok   int `json:"ok"`
	Dane struct {
		FromGrid []ELicznikMeasurement `json:"chart"`
		FeedIn   []ELicznikMeasurement `json:"OZE"`
	} `json:"dane"`
}
