package sparse

type pair struct {
	idx int
	v interface{}
}

type sparseArray struct {
	dense  []pair
	sparse []int
	n      int
}

func (s *sparseArray) add(p pair) {
	s.dense[s.n] = p
	s.sparse[p.idx] = s.n
	s.n++
}

