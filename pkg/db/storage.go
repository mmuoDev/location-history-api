package db

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mmuoDev/location-history-api.git/pkg"
)

const (
	defaultTTL = "5"
)

type Data struct {
	Location      []pkg.Location
	LastUpdatedAt int64
}

//Storage represents a data structure that handles storage
//of orders with their location data. There can be multiple location data to an order
type Storage struct {
	d  map[string]*Data
	mu sync.Mutex
}

//New returns a new data storage
func New() (*Storage, error) {
	s := &Storage{
		d: make(map[string]*Data),
	}
	maxTTL := os.Getenv("LOCATION_HISTORY_TTL_SECONDS")
	if maxTTL == "" {
		maxTTL = defaultTTL
	}
	i, err := strconv.Atoi(maxTTL)
	if err != nil {
		return s, err
	}

	go func() {
		for now := range time.Tick(time.Second) {
			s.mu.Lock()
			for k, v := range s.d {
				if now.Unix()-v.LastUpdatedAt > int64(i) {
					delete(s.d, k)
				}
			}
			s.mu.Unlock()
		}
	}()
	return s, nil
}

//Insert appends a location to a specified order
func (s *Storage) Insert(orderID string, l pkg.Location) {
	s.mu.Lock()
	if s.d[orderID] == nil {
		s.d[orderID] = &Data{}
	}
	s.d[orderID] = &Data{
		Location:      append(s.d[orderID].Location, l),
		LastUpdatedAt: time.Now().Unix(),
	}
	s.mu.Unlock()
}

//Retrieve retrieves all location data for an order
func (s *Storage) Retrieve(orderID string, max int) []pkg.Location {
	if s.d[orderID] == nil {
		s.d[orderID] = &Data{}
	}
	ll := s.d[orderID].Location
	if max == 0 {
		return ll
	}
	if len(ll) > 0 {
		s.mu.Lock()
		res := []pkg.Location{}
		c := 0
		for _, l := range ll {
			res = append(res, l)
			c++
			if max == c {
				break
			}
		}
		s.d[orderID].LastUpdatedAt = time.Now().Unix()
		s.mu.Unlock()
		return res
	}
	return ll
}

//Delete deletes location data for an order
func (s *Storage) Delete(orderID string) {
	delete(s.d, orderID)
}
