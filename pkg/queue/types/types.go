package types

type HandlerFunc func(id string, values map[string]interface{}) error

type Values map[string]interface{}