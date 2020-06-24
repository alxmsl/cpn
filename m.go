package cpn

import "sync"

type M struct {
	v interface{}

	lock sync.RWMutex
	path []string
	word []string
}

func NewM(value interface{}) *M {
	return &M{
		v: value,

		//@todo: set this value based on PN longest path size to reduce memory allocations
		path: []string{},
		word: []string{},
	}
}

func (m *M) PassP(n string) {
	m.lock.RLock()
	if len(m.path) == 0 || m.path[len(m.path)-1] != n {
		m.lock.RUnlock()
		m.lock.Lock()
		defer m.lock.Unlock()
		if len(m.path) == 0 || m.path[len(m.path)-1] != n {
			m.path = append(m.path, n)
		}
		return
	}
	m.lock.RUnlock()
}

func (m *M) PassT(n string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.path = append(m.path, n)
	m.word = append(m.word, n)
}

func (m *M) Path() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.path
}

func (m *M) Value() interface{} {
	return m.v
}

func (m *M) Word() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.word
}
