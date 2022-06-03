package cpn

import (
	"sync"
	"time"
)

type M struct {
	c time.Time

	VV []*V

	lock sync.RWMutex
	path []*E
	word []string
}

type V struct {
	Name  string
	Value interface{}
}

type E struct {
	T time.Time
	N string
}

func NewM(value interface{}) *M {
	return &M{
		c: time.Now(),
		VV: append([]*V{}, &V{
			Value: value,
		}),

		//@todo: set this value based on PN longest path size to reduce memory allocations
		path: []*E{},
		word: []string{},
	}
}

func (m *M) History() []*E {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return append([]*E{{m.c, ""}}, m.path...)
}

func (m *M) PassP(n string) {
	m.lock.RLock()
	if len(m.path) == 0 || m.path[len(m.path)-1].N != n {
		m.lock.RUnlock()
		m.lock.Lock()
		defer m.lock.Unlock()
		if len(m.path) == 0 || m.path[len(m.path)-1].N != n {
			m.path = append(m.path, &E{time.Now(), n})
		}
		return
	}
	m.lock.RUnlock()
}

func (m *M) PassT(n string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.path = append(m.path, &E{time.Now(), n})
	m.word = append(m.word, n)
}

func (m *M) Path() []*E {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.path
}

func (m *M) SetValue(value interface{}) {
	m.VV = append(m.VV, &V{
		Name:  m.word[len(m.word)-1],
		Value: value,
	})
}

func (m *M) Value() interface{} {
	return m.VV[len(m.word)-1].Value
}

func (m *M) Word() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.word
}
