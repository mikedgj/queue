package queue

import (
	"sync"
)

type State int

type Queue struct {
	lock    *sync.Mutex
	data    []one
	order   bool
	compare func(a, b string) bool
}

type one struct {
	key   string
	state State
	data  interface{}
}
