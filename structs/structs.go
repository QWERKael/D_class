package structs

type IntList map[int]struct{}

func NewIntList() IntList {
	return map[int]struct{}{}
}

func (il *IntList) Del(i int) {
	delete(*il, i)
}

func (il *IntList) Add(i int) {
	(*il)[i] = struct{}{}
}

func (il *IntList) Len(i int) int {
	return len(*il)
}
