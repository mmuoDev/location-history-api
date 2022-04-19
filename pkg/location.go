package pkg

//Location represents data points that points to a location
type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

//RetrieveLocationsResponse represent response for retrieving location history for an order
type RetrieveLocationsResponse struct {
	OrderID string     `json:"order_id"`
	History []Location `json:"history"`
}

//QueryParams represents query params needed to filter location data
type QueryParams struct {
	Max int `json:"max"`
}
