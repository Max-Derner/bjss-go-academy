package main

import (
	"errors"
	"github.com/google/uuid"
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
var ErrUnknownAction = errors.New("Unknown action")

type action int64

const (
	Add action = iota
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
type dbRequest struct {
	action
	ToDoItem
	errorReturnChan chan error
	dataReturnChan  chan []ToDoItem
}

func (request *dbRequest) complete(err error, data []ToDoItem) {
	request.errorReturnChan <- err
	close(request.errorReturnChan)
	request.dataReturnChan <- data
	close(request.dataReturnChan)

}

type dataStore struct {
	data     map[Id]ToDoItem
	requests chan dbRequest
}

func NewDataStore(data map[Id]ToDoItem) dataStore {
	db := dataStore{
		data,
		make(chan dbRequest),
	}
	go db.act()
	return db
}

func NewEmptyDataStore() dataStore {
	data := make(map[Id]ToDoItem)
	return NewDataStore(data)
}

func (d *dataStore) act() {
	for request := range d.requests {
		switch request.action {
		case Add:
			err := d.add(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Update:
			err := d.update(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Delete:
			err := d.delete(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Read:
			data := d.read()
			request.complete(nil, data)
		default:
			request.complete(ErrUnknownAction, []ToDoItem{})
		}
	}
}

func (d *dataStore) add(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if keyExists {
		return ErrCannotAdd
	}
	d.data[dataKey] = item
	return nil
}

func (d dataStore) Add(item ToDoItem) error {
	errChan := make(chan error)
	dataChan := make(chan []ToDoItem)
	d.requests <- dbRequest{
		Add,
		item,
		errChan,
		dataChan,
	}
	err := <-errChan
	<-dataChan
	return err
}

func (d *dataStore) update(item ToDoItem) error {
	dataKey := item.Id
	_, keyExists := d.data[dataKey]
	if !keyExists {
		return ErrCannotUpdate
	}
	d.data[dataKey] = item
	return nil
}

func (d dataStore) Update(item ToDoItem) error {
	errChan := make(chan error)
	dataChan := make(chan []ToDoItem)
	d.requests <- dbRequest{
		Update,
		item,
		errChan,
		dataChan,
	}
	err := <-errChan
	<-dataChan
	return err
}

func (d *dataStore) delete(item ToDoItem) error {
	_, keyExists := d.data[item.Id]
	if !keyExists {
		return ErrCannotDelete
	}
	delete(d.data, item.Id)
	return nil
}

func (d dataStore) Delete(item ToDoItem) error {
	errChan := make(chan error)
	dataChan := make(chan []ToDoItem)
	d.requests <- dbRequest{
		Delete,
		item,
		errChan,
		dataChan,
	}
	err := <-errChan
	<-dataChan
	return err
}

func (d dataStore) read() []ToDoItem {
	var dataSlice []ToDoItem
	for _, item := range d.data {
		dataSlice = append(dataSlice, item)
	}
	return dataSlice
}

func (d dataStore) Read() []ToDoItem {
	errChan := make(chan error)
	dataChan := make(chan []ToDoItem)
	d.requests <- dbRequest{
		Read,
		ToDoItem{},
		errChan,
		dataChan,
	}
	<-errChan
	data := <-dataChan
	return data
}
