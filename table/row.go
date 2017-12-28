package table

import (
	"log"
)

type Row struct {
	table     *Table
	values    []interface{}
	strValues []string
}

type Rows []*Row

func (r Rows) Len() int {
	return len(r)
}

func (r Rows) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rows) Less(i, j int) bool {
	sortOrder := r[i].table.sort
	var k int
	for k = 0; k < len(sortOrder)-1; k++ {
		s := sortOrder[k]
		if s.order == SortDesc {
			if gt(r[i].values[s.index], r[j].values[s.index]) {
				return true
			}
			if gt(r[j].values[s.index], r[i].values[s.index]) {
				return false
			}
		} else {
			if lt(r[i].values[s.index], r[j].values[s.index]) {
				return true
			}
			if lt(r[j].values[s.index], r[i].values[s.index]) {
				return false
			}
		}
	}

	s := sortOrder[k]
	if s.order == SortDesc {
		return gt(r[i].values[s.index], r[j].values[s.index])
	}
	return lt(r[i].values[s.index], r[j].values[s.index])
}

func gt(a interface{}, b interface{}) bool {
	return lt(b, a)
}

func lt(a interface{}, b interface{}) bool {
	switch v := a.(type) {
	case int:
		return a.(int) < b.(int)
	case string:
		return a.(string) < b.(string)
	default:
		log.Fatalf("unknown type: %T", v)
	}
	return false
}
