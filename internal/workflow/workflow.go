package workflow

import (
	"github.com/mmuoDev/location-history-api.git/internal/db"
	"github.com/mmuoDev/location-history-api.git/pkg"
)

//AddLocationFunc returns a functionality to add a location data to an order
type AddLocationFunc func(req pkg.Location, orderID string)

//RetrieveHistory returns a functionality to retrieve location history for an order
type RetrieveHistoryFunc func(orderID string, max int) pkg.RetrieveLocationsResponse

//DeleteHistory returns a functionality to delete location history for an order
type DeleteHistoryFunc func(orderID string)

//AddLocation adds location to an order
func AddLocation(addLocation db.AddLocationFunc) AddLocationFunc {
	return func(req pkg.Location, orderID string) {
		addLocation(req, orderID)
	}
}

//RetrieveHistory retrieves location history for an order
func RetrieveHistory(retrieveLocations db.RetrieveHistoryFunc) RetrieveHistoryFunc {
	return func(orderID string, max int) pkg.RetrieveLocationsResponse {
		ll := retrieveLocations(orderID, max)
		return pkg.RetrieveLocationsResponse{
			OrderID: orderID,
			History: ll,
		}
	}
}

//DeleteHistory deletes location history for an order
func DeleteHistory(deleteHistory db.DeleteHistoryFunc) DeleteHistoryFunc {
	return func(orderID string) {
		deleteHistory(orderID)
	}
}