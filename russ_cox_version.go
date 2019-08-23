package NFA

// in contrast to princeton version, the nfa is represented with
// different approach, the chars are the edges in this version,
// whereas the chars are nodes in princeton version

// compiling to NFA
type state struct {
	c         int // c < 256, c = 256, c = 257
	out, out1 *state
	lastList  int
}

type frag struct {
	start *state   // start points to the start state for the frag
	out   *ptrList // out is a list of pointers points to *state pointers that are not yet connected to anything
}

type ptrList struct {
}

// list1 creates a new ptr list containing a single ptr outp
func list1(outp **state) *ptrList      {}
func append_(l1, l2 *ptrList) *ptrList {}
func patch(l *ptrList, s *state)       {}
