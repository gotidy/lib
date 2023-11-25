package slice

import (
	"errors"
	"fmt"
	"math/rand"
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

func TestFoldFunc(t *testing.T) {
	t.Parallel()
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
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := FoldFunc(test.s, strconv.Itoa)
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

func ExampleCount() {
	fmt.Println(Count([]string{"fish", "crab", "", "octopus", "", "squid", "", ""}, func(v string) bool { return v == "" }))

	// Output:
	// 4
}

func ExampleMergeSorted() {
	fmt.Println(MergeSorted([]int{1, 3, 4, 7}, []int{2, 3, 6, 8}, func(v1, v2 int) bool { return v1 < v2 }, 0))
	fmt.Println(MergeSorted([]int{1, 3, 4, 7}, []int{2, 3, 6, 8}, func(v1, v2 int) bool { return v1 < v2 }, 3))

	// Output:
	// [1 2 3 3 4 6 7 8]
	// [1 2 3]
}

func ExampleMergeSortedTo() {
	s1, s2 := []int{1, 3, 4, 7}, []int{2, 3, 6, 8}
	fmt.Println(MergeSortedTo(s1, s1, s2, func(v1, v2 int) bool { return v1 < v2 }, 0))
	s1, s2 = []int{1, 3, 4, 7}, []int{2, 3, 6, 8}
	fmt.Println(MergeSortedTo(s1, s1, s2, func(v1, v2 int) bool { return v1 < v2 }, 3))
	s1, s2 = []int{1, 3, 4, 7}, []int{2, 3, 6, 8}
	fmt.Println(MergeSortedTo(s1, s1, s2, func(v1, v2 int) bool { return v1 < v2 }, 20))
	s1, s2 = nil, []int{2, 3, 6, 8}
	fmt.Println(MergeSortedTo(s1, s1, s2, func(v1, v2 int) bool { return v1 < v2 }, 0))
	s1, s2 = []int{2, 3, 6, 8}, nil
	fmt.Println(MergeSortedTo(s1, s1, s2, func(v1, v2 int) bool { return v1 < v2 }, 0))

	// Output:
	// [1 2 3 3 4 6 7 8]
	// [1 2 3]
	// [1 2 3 3 4 6 7 8]
	// [2 3 6 8]
	// [2 3 6 8]
}

func BenchmarkMergeSorted(b *testing.B) {
	const count = 10
	sizes := []int{10, 100, 1000, 10000}

	bench := func(b *testing.B, size, limit int, limitName string) {
		var resultMerge, resultQuicksort []int
		b.Run("merge_____"+limitName+"_"+strconv.Itoa(size), func(b *testing.B) {
			rand.Seed(1) //nolint
			var ss [][]int
			for i := 0; i < count; i++ {
				s := make([]int, size)
				for j := 0; j < size; j++ {
					s[j] = j*3 + rand.Intn(3) //nolint
				}
				ss = append(ss, s)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resultMerge = ss[0]
				for i := 1; i < len(ss); i++ {
					resultMerge = MergeSortedTo(resultMerge, resultMerge, ss[i], func(v1, v2 int) bool { return v1 < v2 }, limit)
				}
			}
		})

		b.Run("quicksort_"+limitName+"_"+strconv.Itoa(size), func(b *testing.B) {
			rand.Seed(1) //nolint
			var ss [][]int
			for i := 0; i < count; i++ {
				s := make([]int, size)
				for j := 0; j < size; j++ {
					s[j] = j*3 + rand.Intn(3) //nolint
				}
				ss = append(ss, s)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resultQuicksort = nil
				for i := 0; i < len(ss); i++ {
					resultQuicksort = append(resultQuicksort, ss[i]...)
				}
				sort.Slice(resultQuicksort, func(i, j int) bool { return resultQuicksort[i] < resultQuicksort[j] })
			}
		})
		if limit > 0 {
			resultQuicksort = resultQuicksort[:limit]
		}

		if !reflect.DeepEqual(resultMerge, resultQuicksort) {
			b.Fatalf("results (size %d, limit %d) is different len(resultMerge)=%d, len(resultQuicksort)=%d", size, limit, len(resultMerge), len(resultQuicksort))
		}
	}

	for _, size := range sizes {
		bench(b, size, size, "limited=size")
		bench(b, size, 0, "limited=_no_")
	}
}

func TestRemoveManyPanic(t *testing.T) {
	t.Parallel()
	defer func() { _ = recover() }()
	_ = RemoveMany([]string{"0", "1", "2", "3", "4", "5"}, 4, 3)
	t.Errorf("expected panic")
}

func ExampleRemoveFunc() {
	fmt.Println(RemoveFunc([]int{2, 3, 4, 7, 0, 5, 100, 15, 30, 31}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(RemoveFunc([]int{1, 3, 4, 7, 0, 5, 100, 15, 30, 32}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(RemoveFunc([]int{1, 1}, func(i int) bool { return i%2 == 1 }))
	fmt.Println(RemoveFunc([]int{1, 1}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(RemoveFunc([]int{1}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(RemoveFunc([]int{1}, func(i int) bool { return i%2 == 1 }))
	fmt.Println(RemoveFunc([]int{}, func(i int) bool { return i%2 == 0 }))
	fmt.Println(RemoveFunc([]int(nil), func(i int) bool { return i%2 == 0 }))

	// Output:
	// [3 7 5 15 31]
	// [1 3 7 5 15]
	// []
	// [1 1]
	// [1]
	// []
	// []
	// []
}

func TestRemove(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		s        []string
		pos      int
		expected []string
	}{
		{
			name:     "Remove from start",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      0,
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			name:     "Remove from middle",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      3,
			expected: []string{"0", "1", "2", "4", "5"},
		},
		{
			name:     "Remove from end",
			s:        []string{"0", "1", "2", "3", "4", "5"},
			pos:      5,
			expected: []string{"0", "1", "2", "3", "4"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if res := Remove(test.s, test.pos); !reflect.DeepEqual(res, test.expected) {
				t.Errorf("expected %+v, actual %+v", test.expected, res)
			}
		})
	}
}

func TestRemovePanic(t *testing.T) {
	t.Parallel()
	defer func() { _ = recover() }()
	_ = Remove([]string{"0", "1", "2", "3", "4", "5"}, 25)
	t.Errorf("expected panic")
}

func ExampleHasDuplicates() {
	fmt.Println(HasDuplicates([]int{1, 3, 4, 7, 0, 5, 100, 15, 30, 31}))
	fmt.Println(HasDuplicates([]int{1, 3, 4, 7, 0, 5, 100, 3, 30, 31}))
	fmt.Println(HasDuplicates([]int{1, 1}))
	fmt.Println(HasDuplicates([]int{1}))
	fmt.Println(HasDuplicates([]int{}))
	fmt.Println(HasDuplicates([]int(nil)))

	// Output:
	// false
	// true
	// true
	// false
	// false
	// false
}

func ExampleHasDuplicatesFunc() {
	fmt.Println(HasDuplicatesFunc([]int{1, 3, 4, 7, 0, 5, 100, 15, 30, 31}, func(i int) int { return i }))
	fmt.Println(HasDuplicatesFunc([]int{1, 3, 4, 7, 0, 5, 100, 3, 30, 31}, func(i int) int { return i }))
	fmt.Println(HasDuplicatesFunc([]int{1, 1}, func(i int) int { return i }))
	fmt.Println(HasDuplicatesFunc([]int{1}, func(i int) int { return i }))
	fmt.Println(HasDuplicatesFunc([]int{}, func(i int) int { return i }))
	fmt.Println(HasDuplicatesFunc([]int(nil), func(i int) int { return i }))

	// Output:
	// false
	// true
	// true
	// false
	// false
	// false
}

func ExampleFitIndex() {
	fmt.Println(FitIndex(-1, []int{1, 1, 1}))
	fmt.Println(FitIndex(11, []int{1, 1, 1}))
	fmt.Println(FitIndex(1, []int{1, 1, 1}))
	fmt.Println(FitIndex(1, []int{}))

	// Output:
	// 0
	// 2
	// 1
	// -1
}
