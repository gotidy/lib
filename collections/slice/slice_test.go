package slice

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestIndexExists(t *testing.T) {
	s := []string{"0", "1", "2", "3", "4", "5"}
	expected := 2
	if i := Index(s, "2"); i != expected {
		t.Errorf("expected %d, actual %d", expected, i)
	}
}

func TestIndexNotExists(t *testing.T) {
	s := []string{"0", "1", "2", "3", "4", "5"}
	if i := Index(s, "10"); i != -1 {
		t.Errorf("expected -1, actual %d", i)
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name     string
		s        []string
		pos      int
		v        string
		expected []string
	}{
		{
			name:     "Insert at start with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      0,
			v:        "10",
			expected: []string{"10", "0", "1", "2", "3", "4", "5"},
		},
		{
			name:     "Insert at end with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      6,
			v:        "10",
			expected: []string{"0", "1", "2", "3", "4", "5", "10"},
		},
		{
			name:     "Insert in middle with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      4,
			v:        "10",
			expected: []string{"0", "1", "2", "3", "10", "4", "5"},
		},
		{
			name:     "Insert at start with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5", ""}[:6],
			pos:      0,
			v:        "10",
			expected: []string{"10", "0", "1", "2", "3", "4", "5"},
		},
		{
			name:     "Insert at end with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5", ""}[:6],
			pos:      6,
			v:        "10",
			expected: []string{"0", "1", "2", "3", "4", "5", "10"},
		},
		{
			name:     "Insert in middle with reallocate",
			s:        []string{"0", "1", "2", "3", "4", "5", ""}[:6],
			pos:      4,
			v:        "10",
			expected: []string{"0", "1", "2", "3", "10", "4", "5"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if res := Insert(test.s, test.pos, test.v); !reflect.DeepEqual(res, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, res)
			}
		})
	}
}

func TestInsertPanic(t *testing.T) {
	defer func() { _ = recover() }()
	_ = Insert([]string{"0", "1", "2", "3", "4", "5"}, 7, "10")
	t.Errorf("expected panic")
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		s        []string
		expected []string
	}{
		{
			name:     "even",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			expected: []string{"5", "4", "3", "2", "1", "0"},
		},
		{
			name:     "odd",
			s:        []string{"0", "1", "2", "3", "4", "5", "6"},
			expected: []string{"6", "5", "4", "3", "2", "1", "0"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Reverse(test.s)
			if !reflect.DeepEqual(test.s, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, test.s)
			}
		})
	}
}

func TestMap(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	expected := []string{"0", "1", "2", "3", "4", "5"}

	res := Map(s, strconv.Itoa)
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v, actual %+v", expected, res)
	}
}

func TestFilter(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	expected := []int{2, 3}

	res := Filter(s, func(i int) bool {
		return i == 2 || i == 3
	})
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v, actual %+v", expected, res)
	}
}

func TestReduce(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	expected := 20

	res := Reduce(s, 5, func(r, i int) int {
		return r + i
	})
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v, actual %+v", expected, res)
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name     string
		s1, s2   []int
		expected []int
	}{
		{
			name:     "diff",
			s1:       []int{0, 1, 2, 3, 4},
			s2:       []int{2, 3, 5},
			expected: []int{0, 1, 4},
		},
		{
			name:     "no changes",
			s1:       []int{0, 1, 4},
			s2:       []int{2, 3, 5},
			expected: []int{0, 1, 4},
		},
		{
			name:     "empty",
			s1:       []int{},
			s2:       []int{2, 3, 5},
			expected: []int{},
		},
		{
			name:     "nil",
			s1:       nil,
			s2:       []int{2, 3, 5},
			expected: nil,
		},
		{
			name:     "nil 2",
			s1:       []int{2, 3, 5},
			s2:       nil,
			expected: []int{2, 3, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Diff(test.s1, test.s2)
			if !reflect.DeepEqual(s, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, s)
			}
		})
	}
}

func TestSymmetricDiff(t *testing.T) {
	tests := []struct {
		name     string
		s1, s2   []int
		expected []int
	}{
		{
			name:     "diff 1",
			s1:       []int{0, 1, 2, 3, 4},
			s2:       []int{2, 3, 5},
			expected: []int{0, 1, 4, 5},
		},
		{
			name:     "diff 2",
			s1:       []int{0, 1, 4},
			s2:       []int{2, 3, 5},
			expected: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name:     "dif 3",
			s1:       []int{},
			s2:       []int{2, 3, 5},
			expected: []int{2, 3, 5},
		},
		{
			name:     "nil 1",
			s1:       nil,
			s2:       []int{2, 3, 5},
			expected: []int{2, 3, 5},
		},
		{
			name:     "nil 2",
			s1:       []int{2, 3, 5},
			s2:       nil,
			expected: []int{2, 3, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := SymmetricDiff(test.s1, test.s2)
			sort.Ints(s)
			if !reflect.DeepEqual(s, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, s)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name     string
		s1, s2   []int
		expected []int
	}{
		{
			name:     "no empty",
			s1:       []int{2, 3, 5},
			s2:       []int{0, 1, 2, 3, 4},
			expected: []int{2, 3},
		},
		{
			name:     "empty",
			s1:       []int{0, 1, 4},
			s2:       []int{2, 3, 5},
			expected: nil,
		},
		{
			name:     "nil",
			s1:       nil,
			s2:       []int{2, 3, 5},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Intersect(test.s1, test.s2)
			sort.Ints(s)
			if !reflect.DeepEqual(s, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, s)
			}
		})
	}
}

func TestEach(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	var res []int

	Each(s, func(v int) {
		res = append(res, v)
	})
	if !reflect.DeepEqual(s, res) {
		t.Errorf("expected %+v, actual %+v", s, res)
	}
}

func TestFold(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		expected []int
	}{
		{
			name:     "non empty",
			s:        []int{0, 1, 2, 3, 1, 4, 0, 5},
			expected: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name:     "nil",
			s:        nil,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Fold(test.s)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, actual)
			}
		})
	}
}

func TestClone(t *testing.T) {
	s := []int{0, 1, 2, 3, 1, 4, 0, 5}
	res := Clone(s)
	if !reflect.DeepEqual(s, res) {
		t.Errorf("expected %+v, actual %+v", s, res)
	}
}

func ExampleGroup() {
	type data struct {
		Key   int
		Value int
	}

	s := []data{{Key: 4, Value: 5}, {Key: 4, Value: 4}, {Key: 2, Value: 3}, {Key: 1, Value: 2}, {Key: 1, Value: 1}}
	m := Group(s, func(v data) int { return v.Key })

	fmt.Println(m[4])
	fmt.Println(m[2])
	fmt.Println(m[1])
	// Output:
	// [{4 5} {4 4}]
	// [{2 3}]
	// [{1 2} {1 1}]
}

func ExampleGroupOrder() {
	type data struct {
		Key   int
		Value int
	}

	s := []data{{Key: 4, Value: 5}, {Key: 4, Value: 4}, {Key: 2, Value: 3}, {Key: 1, Value: 2}, {Key: 1, Value: 1}}
	m := GroupOrder(s, func(v data) int { return v.Key }, func(s []data, i, j int) bool { return s[i].Value < s[j].Value })

	fmt.Println(m[4])
	fmt.Println(m[2])
	fmt.Println(m[1])
	// Output:
	// [{4 4} {4 5}]
	// [{2 3}]
	// [{1 1} {1 2}]
}

func ExampleMin() {
	fmt.Println(Min(10, 0, 1))
	fmt.Println(Min(1))
	// Output:
	// 0
	// 1
}

func TestMin_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("must panic")
		}
	}()
	_ = Min[int]()
}

func ExampleMax() {
	fmt.Println(Max(10, 0, 1))
	fmt.Println(Max(1))
	// Output:
	// 10
	// 1
}

func TestMax_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("must panic")
		}
	}()
	_ = Max[int]()
}
