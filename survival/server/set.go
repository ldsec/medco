package survivalserver

// Set implements a set of unique string
type Set struct {
	data map[string]struct{}
}

// NewSet set contructor
func NewSet(size int) *Set {
	return &Set{data: make(map[string]struct{}, size)}
}

// Add inserts a new element
func (set *Set) Add(key string) {
	set.data[key] = struct{}{}
}

// Remove removes an element if the element is in the set, otherwise does nothing
func (set *Set) Remove(key string) {
	_, ok := set.data[key]
	if ok {
		delete(set.data, key)
	}
}

// ForEach sequentially operates a function that takes a string as input
func (set *Set) ForEach(instruction func(string)) {
	for key := range set.data {
		instruction(key)
	}

}
