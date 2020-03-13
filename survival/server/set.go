package survivalserver

// Set implements a set of unique string
type Set struct {
	data map[TagID]struct{}
}

// NewSet set contructor
func NewSet(size int) *Set {
	return &Set{data: make(map[TagID]struct{}, size)}
}

// Add inserts a new element
func (set *Set) Add(key TagID) {
	set.data[key] = struct{}{}
}

// Remove removes an element if the element is in the set, otherwise does nothing
func (set *Set) Remove(key TagID) {
	_, ok := set.data[key]
	if ok {
		delete(set.data, key)
	}
}

// ForEach sequentially operates a function that takes a string as input
func (set *Set) ForEach(instruction func(TagID)) {
	for key := range set.data {
		instruction(key)
	}

}
