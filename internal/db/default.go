package db

import (
	"github.com/mmuoDev/location-history-api.git/pkg"
	"github.com/mmuoDev/location-history-api.git/pkg/db"
)

//AddLocation adds location data to an order
func AddLocation(storage db.Storage) AddLocationFunc {
	return func(req pkg.Location, orderID string) {
		storage.Insert(orderID, req)
	}
}

//RetrieveHistory retrieves location history for an order
func RetrieveHistory(storage db.Storage) RetrieveHistoryFunc {
	return func(orderID string, max int) []pkg.Location {
		return storage.Retrieve(orderID, max)
	}
}

//DeleteHistory deletes location history for an order
func DeleteHistory(storage db.Storage) DeleteHistoryFunc {
	return func(orderID string) {
		storage.Delete(orderID)
	}
}
