# Location History Server

An in-memory location history server.

## Functionalities
- Add Location
- Retrieve Location History
- Delete Location History

## Usage
Clone project and `cd` into project foler

### Starting server
``` bash
$ make run
```  

### Running Tests
``` bash
$ make test
```  
### Postman Collection
The Postman collection can be found in the `postman` folder

## Implementing temporary storage 
I have implemented a system that expires the location history for an `orderId`. The env var - `LOCATION_HISTORY_TTL_SECONDS` specifies the number of seconds location history should stay in a map before it is deleted. 
### The logic 
We have a map of `orderId` to the location history (a struct). In the struct, we have a `lastUpdatedAt` attribute that updates to the current timestamp when the data is added to the location history or retrieved from. Every second, a check is made to see if the seconds between now and the struct's `lastUpdatedAt` is greater than `LOCATION_HISTORY_TTL_SECONDS`, the history is deleted. 
