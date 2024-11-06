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

var ErrCannotAdd = errors.New("cannot add item as it already exists in datastore")
var ErrCannotUpdate = errors.New("cannot update item as it does not exist in datastore")
var ErrCannotDelete = errors.New("cannot delete item as it does not exist in datastore")
var ErrCannotQuery = errors.New("cannot query, as item does not exist in datastore")

func ConstructToDoItem(t Title, p Priority, c Complete) ToDoItem {
	return ToDoItem{
		Id(uuid.New()),
		t,
		p,
		c,
	}
}

func NewDataStore() DataStore {
	dataMap := make(map[Id]ToDoItem)
	return DataStore{
		dataMap,
	}
}

type DataStore struct {
	data map[Id]ToDoItem
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

func (d *DataStore) Delete(i Id) error {
	_, keyExists := d.data[i]
	if !keyExists {
		return ErrCannotDelete
	}
	delete(d.data, i)
	return nil
}

func (d DataStore) Read() []ToDoItem {
	var dataSlice []ToDoItem
	for _, item := range(d.data) {
		dataSlice = append(dataSlice, item)
	}
	return dataSlice
}

func (d *DataStore) Query(t Title) (error, ToDoItem) {
	for _, item := range(d.data) {
		if item.Title == t {return nil, item}
	}
	return ErrCannotQuery, ToDoItem{}
}