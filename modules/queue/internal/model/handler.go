package model

type HandlerFunc func(id string, values map[string]interface{}) error
