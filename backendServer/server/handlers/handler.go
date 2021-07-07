package handlers

import "user-data-management/backendServer/db"

//TODO: specify interface

type handler struct {
	Dbw db.DBWrapper
}

func NewHandler(dataSourceName string) *handler {
	handler := handler{}
	handler.Dbw = db.DBWrapper{}
	handler.Dbw.Init(dataSourceName)
	return &handler
}
