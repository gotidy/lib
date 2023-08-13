package slice

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/gotidy/lib/ptr"
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

func TestMapNotNil(t *testing.T) {
	t.Parallel()
	s := []*int{ptr.Of(0), nil, ptr.Of(2), ptr.Of(3), nil, ptr.Of(5)}
	expected := []string{"0", "2", "3", "5"}

	res := MapNotNil(s, func(i *int) string { return strconv.Itoa(*i) })
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v, actual %+v", expected, res)
	}
}

func TestMapFilter(t *testing.T) {
	t.Parallel()
	s := []*int{ptr.Of(0), nil, ptr.Of(2), ptr.Of(3), nil, ptr.Of(5)}
	expected := []string{"0", "2", "3", "5"}

	res := MapFilter(s, func(i *int) (string, bool) {
		if i == nil {
			return "", false
		}
		return strconv.Itoa(*i), true
	})
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %+v, actual %+v", expected, res)
	}
}

func TestMapIndexed(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	expected := []string{"0", "1", "2", "3", "4", "5"}

	var indexes, expectedIndexes []int
	res := MapIndexed(s, func(i, v int) string {
		indexes = append(indexes, i)
		expectedIndexes = append(expectedIndexes, len(expectedIndexes))
		return strconv.Itoa(v)
	})
	if !reflect.DeepEqual(indexes, expectedIndexes) {
		t.Errorf("expected indexes %+v, actual %+v", expectedIndexes, indexes)
	}
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

func TestNew(t *testing.T) {
	expected := []*int{new(int), new(int), new(int), new(int), new(int)}
	result := New[int](len(expected))
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, actual %+v", expected, result)
	}
}

func TestNewInit(t *testing.T) {
	expected := []*int{ptr.Of(0), ptr.Of(1), ptr.Of(2), ptr.Of(3), ptr.Of(4)}
	result := NewInit(len(expected), func(i int, t *int) { *t = i })
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, actual %+v", expected, result)
	}
}

func TestNewInitFilter(t *testing.T) {
	t.Parallel()
	s := []*int{ptr.Of(0), ptr.Of(1), ptr.Of(2), ptr.Of(3), ptr.Of(4)}
	expected := []*int{ptr.Of(0), ptr.Of(2), ptr.Of(4)}
	result := NewInitFilter(len(s), func(i int, t *int) bool { *t = i; return i%2 == 0 })
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, actual %+v", expected, result)
	}
}
func TestNewFrom(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	expected := []*int{ptr.Of(src[0]), ptr.Of(src[1]), ptr.Of(src[2]), ptr.Of(src[3]), ptr.Of(src[4])}
	result := NewFrom(src, func(dst *int, src int) { *dst = src })
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, actual %+v", expected, result)
	}
}

func TestNewFromFilter(t *testing.T) {
	t.Parallel()
	src := []int{0, 1, 2, 3, 4, 5}
	expected := []*int{ptr.Of(src[0]), ptr.Of(src[2]), ptr.Of(src[4])}
	result := NewFromFilter(src, func(dst *int, src int) bool { *dst = src; return src%2 == 0 })
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, actual %+v", expected, result)
	}
}
func BenchmarkNew(b *testing.B) {
	type t struct{ I, J, K, L, M int }

	for size := 2; size <= 1024; size *= size {
		b.Run("New_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New[t](size)
			}
		})

		// b.Run("New_Unsafe_"+strconv.Itoa(size), func(b *testing.B) {
		// 	for i := 0; i < b.N; i++ {
		// 		NewUnsafe[t](size)
		// 	}
		// })

		b.Run("Classic_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := make([]*t, size)
				for i := range s {
					s[i] = &t{}
				}
			}
		})
	}
}

func BenchmarkNewInit(b *testing.B) {
	type t struct{ i, j, k, l, m int }

	for size := 2; size <= 1024; size *= size {
		b.Run("New_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				NewInit(size, func(i int, v *t) { *v = t{1, 2, 3, 4, 5} })
			}
		})

		b.Run("Classic_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := make([]*t, size)
				for i := range s {
					s[i] = &t{1, 2, 3, 4, 5}
				}
			}
		})
	}
}

func ExampleBatch() {
	const size = 4
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var results [][]int
	err := Batch(s, size, func(s []int) error {
		results = append(results, s)
		return nil
	})

	fmt.Println(err)
	fmt.Println(results)

	// with error
	err = Batch(s, size, func(s []int) error {
		return errors.New("oops")
	})

	fmt.Println(err)

	// empty
	results = nil
	_ = Batch(nil, size, func(s []int) error {
		results = append(results, s)
		return nil
	})

	fmt.Println(results)

	// Output:
	// <nil>
	// [[1 2 3 4] [5 6 7 8] [9]]
	// oops
	// []
}

func ExampleMerge() {
	fmt.Println(Merge([]int{1, 2, 3, 4}, []int{4, 5, 6}, nil, []int{5, 7}))
	// Output:
	// [1 2 3 4 4 5 6 5 7]
}

func ExampleUnion() {
	fmt.Println(Union([]int{1, 2, 3, 4}, []int{4, 5, 6}, []int{5, 7}))
	fmt.Println(Union([]int{4, 5, 6}, []int{1, 2, 3, 4}, []int{5, 7}))
	fmt.Println(Union(nil, []int{4, 5, 6}, nil))
	fmt.Println(Union([]string(nil)))

	// Output:
	// [1 2 3 4 5 6 7]
	// [4 5 6 1 2 3 7]
	// [4 5 6]
	// []
}

func TestAppendNotNil(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		got, want []*int
	}{
		{
			name: "nil in middle",
			got:  AppendNotNil([]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}, ptr.Of(4), nil, ptr.Of(6), ptr.Of(7), ptr.Of(8), nil, ptr.Of(10)),
			want: []*int{ptr.Of(1), ptr.Of(2), ptr.Of(3), ptr.Of(4), ptr.Of(6), ptr.Of(7), ptr.Of(8), ptr.Of(10)},
		},
		{
			name: "nil in start and end",
			got:  AppendNotNil([]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}, nil, nil, ptr.Of(6), ptr.Of(7), ptr.Of(8), nil),
			want: []*int{ptr.Of(1), ptr.Of(2), ptr.Of(3), ptr.Of(6), ptr.Of(7), ptr.Of(8)},
		},
		{
			name: "elements is empty",
			got:  AppendNotNil([]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}),
			want: []*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)},
		},
		{
			name: "one nil element",
			got:  AppendNotNil([]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}, nil),
			want: []*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)},
		},
		{
			name: "slice is empty",
			got:  AppendNotNil(nil, ptr.Of(4), nil, ptr.Of(6)),
			want: []*int{ptr.Of(4), ptr.Of(6)},
		},
		{
			name: "all empty",
			got:  AppendNotNil[int](nil),
			want: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("want %+v, got %+v", test.want, test.got)
			}
		})
	}
}
