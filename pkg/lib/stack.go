package lib

type stack struct {
	s []byte
}

func NewStack() *stack {
	return &stack{make([]byte, 0)}
}

func (s *stack) Size() int {
	return len(s.s)
}
func (s *stack) Push(v byte) {
	s.s = append(s.s, v)
}

func (s *stack) Pop() byte {
	l := len(s.s)
	if l == 0 {
		return 0
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res
}
