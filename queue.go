package queue

import (
	"sync"
)

func NewQueue() *Queue {
	return &Queue{
		lock:    new(sync.Mutex),
		data:    make([]one, 0),
		order:   false,
		compare: nil,
	}
}

func NewOrderQueue(fn func(string, string) bool) *Queue {
	return &Queue{
		lock:    new(sync.Mutex),
		data:    make([]one, 0),
		order:   true,
		compare: fn,
	}
}

func (q *Queue) Push(k string, d interface{}, s State) {

	if k == "" || s == 0 {
		return
	}

	q.lock.Lock()
	defer q.lock.Unlock()
	ne := one{
		key:   k,
		data:  d,
		state: s,
	}
	if q.order {
		for i := len(q.data) - 1; i > 0; i-- {
			if q.data[i].key == k {
				return
			}
			if q.compare(q.data[i].key, k) {
				l := q.data[i:]
				q.data = append(q.data[:i], ne)
				q.data = append(q.data, l...)
				return
			}
		}
		q.data = append([]one{ne}, q.data...)
	} else {
		q.data = append(q.data, ne)
	}
}

func (q *Queue) Update(k string, s State) {
	if k == "" || s == 0 {
		return
	}

	q.lock.Lock()
	defer q.lock.Unlock()
	for i := 0; i < len(q.data); i++ {
		if q.data[i].key == k {
			q.data[i].state = s
			return
		}
	}
}

func (q *Queue) GetState(k string) State {
	q.lock.Lock()
	defer q.lock.Unlock()
	for i := 0; i < len(q.data); i++ {
		if q.data[i].key == k {
			return q.data[i].state
		}
	}
	return 0
}

func (q *Queue) GetData(k string) interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	for i := 0; i < len(q.data); i++ {
		if q.data[i].key == k {
			return q.data[i].data
		}
	}
	return nil
}

func (q *Queue) Next(s State) (string, interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if s == 0 {
		return q.data[0].key, q.data[0].data
	}
	for i := 0; i < len(q.data); i++ {
		if q.data[i].state == s {
			return q.data[i].key, q.data[i].data
		}
	}
	return "", nil
}

func (q *Queue) Pull() (string, interface{}, State) {
	q.lock.Lock()
	defer q.lock.Unlock()
	k := q.data[0].key
	d := q.data[0].data
	s := q.data[0].state
	q.data = q.data[1:]
	return k, d, s
}

func (q *Queue) PullByKey(k string) (interface{}, State) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for i := 0; i < len(q.data); i++ {
		if q.data[i].key == k {
			d := q.data[i].data
			s := q.data[i].state
			q.data = append(q.data[:i], q.data[i+1:]...)
			return d, s
		}
	}
	return nil, 0
}

func (q *Queue) Delete(ks ...string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(ks) > 0 {
		for _, k := range ks {
			for i := 0; i < len(q.data); i++ {
				if q.data[i].key == k {
					q.data = append(q.data[:i], q.data[i+1:]...)
					break
				}
			}
		}
	}
}
