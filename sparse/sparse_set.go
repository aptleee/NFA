package sparse

// https://research.swtch.com/sparse

// the sleaziness of uninitialized data access is offset
// by performance improvements: some important operations
// change from linear to constant

// develop a technique to initialize an entry of a matrix
// to zero the first time it's accessed, thereby eliminating
// O(V^2) time to init an adjacency matrix

// One problem with trading more space for less time is that
// initializing the space can itself take a great deal of time.
// Show how to circumvent this problem by designing a technique
// to initialize an entry of a vector to zero the first time
// it is accessed. Your scheme should use constant time for
// initialization and each vector access; you may use extra space
// proportional to the size of the vector. Because this method
// reduces initialization time by using even more space, it
// should be considered only when space is cheap, time is dear,
// and the vector is sparse.

type sparseSet struct {
	dense, sparse []int
	n             int
}

func (s *sparseSet) add(i int) {
	s.dense[s.n] = i
	s.sparse[i] = s.n
	s.n++
}

func (s *sparseSet) has(i int) bool {
	return s.sparse[i] < s.n && s.dense[s.sparse[i]] == i
}

func (s *sparseSet) iter(f func(i int)) {
	for i := 0; i < s.n; i++ {
		f(s.dense[i])
	}
}

func (s *sparseSet) clear() {
	s.n = 0
}

// in contrast with bit vector, all operations are faster
// the only problem is the space cost: two words replace a bit

// a situation where sparse array are better choice is work
// queue-based graph traversal algorithm
// iteration over sparse sets visit elements in the order they
// are inserted, so that new entries added during the iteration
// will be visited later in the same iteration




