package survivalserver

type Set struct {
	data map[string]struct{}
}

func NewSet(size int) *Set {
	return &Set{data: make(map[string]struct{}, size)}
}

func (set *Set) Add(key string) {
	set.data[key] = struct{}{}
}

func (set *Set) Remove(key string) {
	_, ok := set.data[key]
	if ok {
		delete(set.data, key)
	}
}

func (set *Set) ForEach(instruction func(string)) {
	for key := range set.data {
		instruction(key)
	}

}
