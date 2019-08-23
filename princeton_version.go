package NFA

import (
	"bufio"
	"fmt"
	"os"
)

// consider three basic operations in terms of an
// abstract machine that can search patterns in a
// text string

// concatenating: {AB}

// or: use vertical bar to denote the op A|B -> {A, B}
// AB | BCD -> {AB, BCD}

// closure: allow part of the pattern to be repeated arbitrarily
// AB* specifies the language consisting of strings with an A
// followed by 0 or more Bs; A*B specifies the language consisting
// of 0 or more As followed by a B.
// the empty string, we denote by e, is found in every text string

// parentheses: we use parentheses to override default precedence
// C(AD|B)K -> {CADK, CBK},

// definition
// a regex is either
// empty
// a single char
// a regex enclosed in parentheses
// two or more concatenated regex
// two or more regex separated by |
// a regex followed by the closure operator *

// continued
// the empty re represents the empty set of strings, with 0 elem

// a char represents the set of strings with 1 elem, itself

// an re enclosed in parentheses represent the same set of strings
// as re without

// the re consisting of two concatenated res represent the cross
// product of the set of the strings represented by the individual
// components (all possible strings that can be formed by taking
// one string from each and concatenate them, in the same order
// as the res)

// the res consisting of the or of two res represents the union
// of the set represented by the individual components

// the re consisting of the closure of an re represent e or the
// union of the sets represented by concatenating any number of
// copies of the re

// shortcuts
// set-of-characters descriptors
// the dot is a wildcard the represents any single char: A.B

// specified set: a sequence of chars within a square brackets represent
// any one of those chars [ABCD]
// range: enclosed in [] separated by -, [A-Z]
// complement: enclosed in [] preceded by ^
// those are simply shortcuts for a sequence of or

// closure shortcuts
// +: at least 1 copy
// ?: 0 or 1 copy
// a cnt or range within {} to specify a given number of copies

// escape sequence

// substr search: to search for a pat in a text string txt is to check
// whether txt is in the language described by the pattern, .*pat.*

// validity checking:
// multiple of three (0|1(01*0)*1)*

// non-deterministic finite-state automata
// the finite-state automaton for KMP changes from state to state by
// looking at a char from the text string and then changing to another
// state. depending on the char. The automaton reports a match if and
// only if it reaches the accept state: each state transition is
// completely determined by the next char in the text

// the NFA corresponding to an re of length M has exactly one state per
// pattern char, starts at state 0, and has (virtual) accept state M

// states corresponding to a char from the alphabet have an outgoing edge
// that goes to the state corresponding to the next char in the pattern

// states corresponding to the metacharacters (, ), |, and * have at least
// one outgoing edge, which may go to any other state

// some states have multiple outgoing edges, but no states have more than
// one outgoing black edge

// basic diff from DFAs
// chars appear in nodes, not in edges
// NFA organizes a text string only after explicitly reading all its chars
// whereas our DFA recognizes a pattern in a text without necessarily reading
// all the char

// the rules for moving from one state to another:
// if the current state corresponds to a char in alphabet and the current
// char in the text string matches the char, the automaton can scan past
// the char in the text string and take the black transition to the next
// state, we refer to such a transition as a match transition

// the automaton can follow any red edge to another state without scanning
// any text char. we refer to such a transition as an empty transition

// we say an NFA recognizes a text string if and only if there is some
// sequence of transitions that scan all the text chars and ends in the
// accept state when started at the beginning of the text in state 0.

// representation:
// the natural representation of empty-transitions is a digraph

// to simulate NFA, we keep track of the set of states of that could
// possibly be encountered while the automaton is examining the cur
// input char.
// The key computation is like multi-src reachability
// to init the set, we find the sets of states reachable via empty
// transitions from state 0. For each such state, we check whether
// a match transition for the first input char is possible

// (empty-transition -> match transition) -> ... -> (empty-transition -> match transition) -> empty-transition
// dfs:
// starting from the vertex s, explore the graph as deeply as possible
// then backtrack

// for each char in N, we iterate through a set of states of size no
// more than M and run a DFS on the digraph of empty transition,
// the worst-case time for each dfs is O(M)

// 1. try the first edge out of s, towards some node v
// 2. continue from v until you reach a dead end, that is
// a node whose neighborhood have all been explored
// 3. backtrack the first node with an unexplored neighbor
// and repeat 2

// backreference like \1 and \2 match the string matched by the
// previous parenthesised expression, and only that string
// like (cat | dog)\1 matches catcat or dogdog

// dfa
// in any state, each possible input leads to at most one new state

// nfa
// in any state, each possible input leads to zero or more new states

// there are multiple ways to translate regexp to NFA
// the NFA for a regular expression is built up from partial NFAs for each
// subexpression, with a different construction for each operator
// the partial NFAs have no matching states

// support (, ), |, *, +, [ABCD], [A-Z]
type graph map[int]map[int]bool

func (g graph) addEdge(from, to int) {
	g[from][to] = true
}

type intStack []int

func (s *intStack) push(i int) {
	*s = append(*s, i)
}

func (s *intStack) pop() int {
	n := len(*s)
	x := (*s)[n-1]
	*s = (*s)[:n-1]
	return x
}

func (s *intStack) top() int {
	return (*s)[len(*s)-1]
}

type NFA struct {
	g  graph
	re []byte
	m  map[int]int
}

var metaChar = map[byte]bool{
	'(': true,
	'|': true,
	'[': true,
	')': true,
	']': true,
	'*': true,
	'+': true,
	'?': true,
}

// this implementation assumes that the input regexp is right format
// TODO: add format checking, and if wrong fomat detected, return error
func buildNFA(regex string) *NFA {
	g := make(graph)
	ops := make(intStack, 0)
	re := []byte(regex)
	M := len(regex)
	rightIdx := make(map[int]int) //[abc], map[idx of a]idx of ']'
	for i := 0; i < M; i++ {
		lp := i // left parenthesis
		switch re[i] {
		case '(', '|', '[':
			ops.push(i)
		case ')':
			orOpIdxes := make(map[int]bool)
			for re[ops.top()] == '|' {
				or := ops.pop()
				orOpIdxes[or] = true
			}
			lp = ops.pop()
			for k := range orOpIdxes {
				g.addEdge(k, i)
				g.addEdge(lp, k+i)
			}
		case ']':
			lp = ops.pop()
			for j := lp + 1; j < i; j++ {
				g.addEdge(lp, j)
				// If a match occurs while checking the characters in this set, the NFA will go to
				// the right square bracket state.
				rightIdx[j] = i
				if re[j+1] == '-' {
					j += 2
				}
			}
		}
		if i < M-1 {
			switch re[i+1] {
			case '*':
				g.addEdge(lp, i+1) // zero or one
				g.addEdge(i+1, lp) // more than one
			case '+':
				g.addEdge(i+1, lp) // more than one
			case '?':
				g.addEdge(lp, i+1) // zero or one
			}
		}
		if metaChar[re[i]] {
			g.addEdge(i, i+1)
		}
	}
	return &NFA{
		g:  g,
		re: re,
		m:  rightIdx,
	}
}

func multipleSrcReach(g graph, srcs []int) map[int]bool {
	marked := make(map[int]bool)
	var dfs func(s int)
	dfs = func(s int) {
		marked[s] = true
		for k := range g[s] {
			if !marked[k] {
				dfs(k)
			}
		}
	}
	for _, src := range srcs {
		if !marked[src] {
			dfs(src)
		}
	}
	return marked
}

func (nfa *NFA) recognize(txt string) bool {
	M := len(nfa.re)
	pc := multipleSrcReach(nfa.g, []int{0})
	for i := 0; i < len(txt); i++ {
		var match []int
		for k := range pc {
			if k < M {
				if _, ok := nfa.m[k]; ok {
					// in a range
					if nfa.re[k+1] == '-' {
						left, right := nfa.re[k], nfa.re[k+2]
						if left <= txt[i] && txt[i] <= right {
							match = append(match, nfa.m[k])
						}
					} else if nfa.re[k] == txt[i] || nfa.re[k] == '.' {
						match = append(match, nfa.m[k])
					}
				} else if nfa.re[k] == txt[i] || nfa.re[k] == '.' {
					match = append(match, k+1)
				}
			}
			pc = multipleSrcReach(nfa.g, match)
			if len(pc) == 0 {
				return false
			}
		}
	}
	for k := range pc {
		if k == M {
			return true
		}
	}
	return false
}

func grep(pat string) {
	regexp := "(.*" + pat + ".*)"
	nfa := buildNFA(regexp)
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanLines)
	for in.Scan() {
		txt := in.Text()
		if nfa.recognize(txt) {
			fmt.Println(txt)
		}
	}
}
