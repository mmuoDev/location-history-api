package db

import "github.com/mmuoDev/location-history-api.git/pkg"

//AddLocationFunc returns a functionality to add a location data to an order
type AddLocationFunc func(req pkg.Location, orderID string) 

//RetrieveHistory returns a functionality to retrieve location history for an order
type RetrieveHistoryFunc func(orderID string, max int) []pkg.Location

//DeleteHistory returns a functionality to delete location history for an order
type DeleteHistoryFunc func(orderID string) 
