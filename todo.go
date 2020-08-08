package main

type Todo struct {
	ID   int    `json:"id,omitempty"`
	Todo string `json:"todo,omitempty"`
	Done bool   `json:"done,omitempty"`
}
