# Basically...

Start with [DataAccessLayer.go](./DataAccessLayer.go), it defines a thread safe DAL which takes in a DataStore interface which is also defined within the same file. A new DAL is created with the `NewDataAccessLayer(db DataStore)` method this method injects your DataStore and spins up a goroutine that listens on a channel for `dbRequest`s and acts on the DataStore in a one at a time fashion.

[inMemDataStore.go](./inMemDataStore.go) defines an in memory ephemeral data store  
[jsonDataStore.go](./jsonDataStore.go) defines a persistent datastore that writes and reads data from a JSON file  

[restApi.go](./restApi.go) defines a RESTful API that is served on port 8080. This can be interacted with via [a Python script](./py/main.py)
[website.go](./website.go) defines a poor website served on port 6060.

[main.go](./main.go) coordinates all of this to run at the same time.

# To use

issue command `go run main`
