package monitoring

import "sync"

type Observable struct {
	data      map[string]interface{}
	listeners []func()
	mu        sync.RWMutex
}

func NewObservable() *Observable {
	return &Observable{
		data: make(map[string]interface{}),
	}
}

func (o *Observable) Add(key string, value interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.data[key] = value
	o.notify()
}

func (o *Observable) Remove(key string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.data, key)
	o.notify()
}

func (o *Observable) Get(key string) (interface{}, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	value, exists := o.data[key]
	return value, exists
}

func (o *Observable) RegisterListener(listener func()) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listeners = append(o.listeners, listener)
}

func (o *Observable) notify() {
	for _, listener := range o.listeners {
		listener()
	}
}
