package main

func New() interface{} {
	return &Plugin{}
}

type Plugin struct{}
