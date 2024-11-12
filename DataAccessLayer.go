package main

import (
	"errors"
)

type DataStore interface {
	create(item ToDoItem) error
	read() []ToDoItem
	update(item ToDoItem) error
	delete(item ToDoItem) error
}

var ErrCannotCreate = errors.New("cannot create item as it already exists in datastore")
var ErrCannotUpdate = errors.New("cannot update item as it does not exist in datastore")
var ErrCannotDelete = errors.New("cannot delete item as it does not exist in datastore")
var ErrCannotQuery = errors.New("cannot query, as item does not exist in datastore")
var ErrUnknownAction = errors.New("Unknown action")

type action int64

const (
	Create action = iota
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

func NewDataStore(data map[Id]ToDoItem) DataAccessLayer {
	db := inMemoryDataStore{
		data,
	}
	dal := DataAccessLayer{
		&db,
		make(chan dbRequest),
	}
	go dal.act()
	return dal
}

func NewEmptyDataStore() DataAccessLayer {
	data := make(map[Id]ToDoItem)
	return NewDataStore(data)
}

type DataAccessLayer struct {
	db DataStore
	requests chan dbRequest
}

func (d DataAccessLayer) Create(item ToDoItem) error {
	errChan := make(chan error)
	dataChan := make(chan []ToDoItem)
	d.requests <- dbRequest{
		Create,
		item,
		errChan,
		dataChan,
	}
	err := <-errChan
	<-dataChan
	return err
}

func (d DataAccessLayer) Update(item ToDoItem) error {
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

func (d DataAccessLayer) Delete(item ToDoItem) error {
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

func (d DataAccessLayer) Read() []ToDoItem {
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

func (request *dbRequest) complete(err error, data []ToDoItem) {
	request.errorReturnChan <- err
	close(request.errorReturnChan)
	request.dataReturnChan <- data
	close(request.dataReturnChan)
}

func (d *DataAccessLayer) act() {
	for request := range d.requests {
		switch request.action {
		case Create:
			err := d.db.create(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Update:
			err := d.db.update(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Delete:
			err := d.db.delete(request.ToDoItem)
			request.complete(err, []ToDoItem{})
		case Read:
			data := d.db.read()
			request.complete(nil, data)
		default:
			request.complete(ErrUnknownAction, []ToDoItem{})
		}
	}
}
