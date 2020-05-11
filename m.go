package cpn

type M struct {
	v    interface{}
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

func (m *M) Path() []string {
	return m.path
}

func (m *M) Value() interface{} {
	return m.v
}

func (m *M) Word() []string {
	return m.word
}
