package zenfo

// Venue represents a venue entry for DB store and web API
type Venue struct {
	Name    string  `json:"name"`
	Addr    string  `json:"addr"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
	Website string  `json:"website"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
}
