package workers

import "sync"

type JobsMap struct {
	mx sync.RWMutex
	m  map[string]string
}

func (m *JobsMap) Load(key string) (string, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *JobsMap) Store(key string, value string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m[key] = value
}
