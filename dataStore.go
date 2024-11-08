package main

import (
	"github.com/google/uuid"
	"errors"
)

type Id uuid.UUID
type Title string
type Priority string
type Complete bool

type ToDoItem struct {
	Id
	Title
	Priority
	Complete
}

func ConstructToDoItem(t Title, p Priority, c Complete) ToDoItem {
	return ToDoItem{
		Id(uuid.New()),
		t,
		p,
		c,
	}
}

var ErrCannotAdd = errors.New("cannot add item as it already exists in datastore")
var ErrCannotUpdate = errors.New("cannot update item as it does not exist in datastore")
var ErrCannotDelete = errors.New("cannot delete item as it does not exist in datastore")
var ErrCannotQuery = errors.New("cannot query, as item does not exist in datastore")

type Action int64

const (
	Add Action = iota
	Update
	Delete
	Read
)

// Use for sending requests to DB
// 
// Wait on `errorReturnChan` for errors, if nil is returned then action was successful  
// Wait on `dataReturnChan` if you anticipate data being returned  
// 
// Both channels with have 1 item of data sent down it and then get closed
type DbRequest struct {
	Action
	ToDoItem
	errorReturnChan chan error
	dataReturnChan chan []ToDoItem
}

func (request *DbRequest) Complete(err error, data []ToDoItem) {
	request.errorReturnChan <- err
	close(request.errorReturnChan)
	request.dataReturnChan <- []ToDoItem{}
	close(request.dataReturnChan)

}

type DataStore struct {
	data map[Id]ToDoItem
	requests chan DbRequest
}

func NewDataStore() DataStore {
	dataMap := make(map[Id]ToDoItem)
	db := DataStore{
		dataMap,
		make(chan DbRequest),
	}
	go db.act()
	return db
}

func (d *DataStore) act() {
	for request := range(d.requests) {
		switch request.Action {
		case Add:
			err := d.Add(request.ToDoItem)
			request.Complete(err, []ToDoItem{})
		case Update:
			err := d.Update(request.ToDoItem)
			request.Complete(err, []ToDoItem{})
		case Delete:
			err := d.Delete(request.ToDoItem)
			request.Complete(err, []ToDoItem{})
		case Read:
			data := d.Read()
			request.Complete(nil, data)
		}
	}
}

func (d *DataStore) Add(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if keyExists {
		return ErrCannotAdd
	}
	d.data[dataKey] = item
	return nil
}

func (d *DataStore) Update(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if !keyExists {
		return ErrCannotUpdate
	}
	d.data[dataKey] = item
	return nil
}

func (d *DataStore) Delete(item ToDoItem) error {
	_, keyExists := d.data[item.Id]
	if !keyExists {
		return ErrCannotDelete
	}
	delete(d.data, item.Id)
	return nil
}

func (d DataStore) Read() []ToDoItem {
	var dataSlice []ToDoItem
	for _, item := range(d.data) {
		dataSlice = append(dataSlice, item)
	}
	return dataSlice
}
