package pn

type M struct {
	v    int
	path []string
	word []string
}

func NewM(value int) *M {
	return &M{
		v: value,

		//@todo: set this value based on PN longest path size to reduce memory allocations
		path: []string{},
		word: []string{},
	}
}

func (m *M) Path() []string {
	return m.path
}

func (m *M) Value() int {
	return m.v
}

func (m *M) Word() []string {
	return m.word
}
