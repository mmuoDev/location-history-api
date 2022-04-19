package db

import "github.com/mmuoDev/location-history-api.git/pkg"

//Storage represents a data structure that handles storage
//of orders with their location data. There can be multiple location data to an order
type Storage struct {
	data map[string][]pkg.Location
}

//New returns a new data storage
func New() *Storage {
	return &Storage{
		data: make(map[string][]pkg.Location),
	}
}

//Insert appends a location to a specified order
func (s *Storage) Insert(orderID string, l pkg.Location) {
	s.data[orderID] = append(s.data[orderID], l)
}

//Retrieve retrieves all location data for an order
func (s *Storage) Retrieve(orderID string, max int) []pkg.Location {
	ll := s.data[orderID]
	if max == 0 {
		return ll
	}
	if len(ll) > 0 {
		res := []pkg.Location{}
		c := 0
		for _, l := range ll {
			res = append(res, l)
			c++
			if max == c {
				break
			}
		}
		return res
	}
	return ll
}

//Delete deletes location data for an order
func (s *Storage) Delete(orderID string) {
	if _, ok := s.data[orderID]; ok {
		s.data[orderID] = make([]pkg.Location, 0)
	}
}
