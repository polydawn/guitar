package format

//Implement's golang's sorting interface for Headers.

import (
	"sort"
)

type HeaderSorter struct {
	Headers []*Header
	By      func(h1, h2 *Header) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *HeaderSorter) Len() int {
	return len(s.Headers)
}

// Swap is part of sort.Interface.
func (s *HeaderSorter) Swap(i, j int) {
	s.Headers[i], s.Headers[j] = s.Headers[j], s.Headers[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *HeaderSorter) Less(i, j int) bool {
	return s.By(s.Headers[i], s.Headers[j])
}

//Generic sort function
func SortHeaders(headers []*Header, by func(h1, h2 *Header) bool) {
	sortHeaders := &HeaderSorter {
		Headers: headers,
		By: by,
	}

	sort.Sort(sortHeaders)
}

//Convenience function
func SortHeadersByName(headers []*Header) {
	byName := func(h1, h2 *Header) bool {
		return h1.Name < h2.Name
	}

	SortHeaders(headers, byName)
}
