package monitoring

import "sync"

type Observable struct {
	data      map[string]interface{}
	listeners []func(key string, value interface{}, action string)
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
	o.notify(key, value, "save")
}

func (o *Observable) Remove(key string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	value := o.data[key]
	delete(o.data, key)
	o.notify(key, value, "remove")
}

func (o *Observable) Get(key string) (interface{}, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	value, exists := o.data[key]
	return value, exists
}

func (o *Observable) RegisterListener(listener func(key string, value interface{}, action string)) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listeners = append(o.listeners, listener)
}

func (o *Observable) notify(key string, value interface{}, action string) {
	for _, listener := range o.listeners {
		listener(key, value, action)
	}
}
