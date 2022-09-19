package cpn

import (
	"sync"
	"time"
)

// M is an abstraction to define a token in PN
type M struct {
	c time.Time
	// v contains the current mark value
	v interface{}
	// vv contains mark values related to places
	vv []*v

	lock sync.RWMutex
	// path contains all edges - both places and transitions - are passed by the mark
	path []*E
	// word contains all transitions are passed by the mark
	word []string
}

type E struct {
	T time.Time
	N string
}

// The v struct represents a mark's value written from the specific place
type v struct {
	p *P
	v interface{}
}

func NewM(value interface{}) *M {
	return &M{
		c: time.Now(),
		v: value,

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

// passP is called when the mark passed place in the net
func (m *M) passP(p *P) {
	m.lock.RLock()
	if len(m.path) == 0 || m.path[len(m.path)-1].N != p.name {
		m.lock.RUnlock()
		m.lock.Lock()
		defer m.lock.Unlock()
		if len(m.path) == 0 || m.path[len(m.path)-1].N != p.name {
			m.path = append(m.path, &E{time.Now(), p.name})
			if m.v != nil {
				m.vv = append(m.vv, &v{p, m.v})
				m.v = nil
			}
		}
		return
	} else if m.v != nil {
		m.lock.RUnlock()
		m.lock.Lock()
		defer m.lock.Unlock()
		m.vv = append(m.vv, &v{p, m.v})
		m.v = nil
		return
	}
	m.lock.RUnlock()
}

// passT is called when the mark passed transition in the net
func (m *M) passT(t *T) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.path = append(m.path, &E{time.Now(), t.name})
	m.word = append(m.word, t.name)
}

func (m *M) Path() []*E {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.path
}

// SetValue sets the current value to the mark
func (m *M) SetValue(v interface{}) {
	m.v = v
}

// Value returns the current value.
// If the current value is nil, then tries to get the latest written value by using ValueByPlace function
func (m *M) Value() interface{} {
	if m.v != nil {
		return m.v
	}
	return m.ValueByPlace("", 0)
}

// ValueByPlace returns the mark's value specified from the place with name and deep.
// If name is empty, then function returns just latest written value.
// Parameter `deep` is used to define how deep place should be used, what is actual for nets with loop
func (m *M) ValueByPlace(name string, deep int) interface{} {
	if name == "" && len(m.vv) == 0 {
		return m.v
	}

	var c int
	for i := len(m.vv) - 1; i >= 0; i -= 1 {
		if name == "" {
			return m.vv[i].v
		}
		if m.vv[i].p.name == name {
			if c == deep {
				return m.vv[i].v
			}
			c += 1
		}
	}
	return nil
}

func (m *M) Word() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.word
}
