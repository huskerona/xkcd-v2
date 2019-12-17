package model

type XKCD struct {
	Day string `json:"day"`
	Month string `json:"month"`
	Year string `json:"year"`
	Number int `json:"num"`
	Title string `json:"title"`
	SafeTitle string `json:"safe_title"`
	Transcript string `json:"transcript"`
	ImageURL string `json:"img"`
	ImageAlt string `json:"alt"`
	News string `json:"news"`
	Link string `json:"link"`

	// base64 encoded image downloaded from ImageURL
	Image string
}
