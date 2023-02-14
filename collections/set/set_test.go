package set

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		members []string
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{members: []string{"foo", "bar"}},
			want: Set[string]{"foo": struct{}{}, "bar": struct{}{}},
		},
		{
			name: "empty",
			args: args{members: nil},
			want: Set[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Of(tt.args.members...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromMapKeys(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{m: map[string]string{"foo": "1", "bar": "2"}},
			want: Set[string]{"foo": struct{}{}, "bar": struct{}{}},
		},
		{
			name: "nil",
			args: args{m: nil},
			want: Set[string]{},
		},
		{
			name: "empty",
			args: args{m: map[string]string{}},
			want: Set[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OfMapKeys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromMapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromSliceFunc(t *testing.T) {
	type T struct {
		Name string
	}

	type args struct {
		s []T
		f func(v T) string
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s: []T{{Name: "foo"}, {Name: "bar"}}, f: func(v T) string { return v.Name }},
			want: Of("foo", "bar"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromSliceFunc(tt.args.s, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromSliceFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetLen(t *testing.T) {
	members := []string{"a", "b", "c"}
	if got := Of(members...).Len(); got != len(members) {
		t.Errorf("Len() = %v, want %v", got, len(members))
	}
}

func TestSetEmpty(t *testing.T) {
	type args struct {
		members []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with members",
			args: args{members: []string{"foo", "bar"}},
			want: false,
		},
		{
			name: "empty",
			args: args{members: nil},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Of(tt.args.members...).Empty(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAdd(t *testing.T) {
	members := []string{"a", "b", "c"}
	want := Of(members...)
	if got := Of[string]().Add(members...); !reflect.DeepEqual(got, want) {
		t.Errorf("Add() = %v, want %v", got, want)
	}
}

func TestSetDelete(t *testing.T) {
	members := []string{"a", "b", "c"}
	want := Of(members[:2]...)
	if got := Of(members...).Delete("c"); !reflect.DeepEqual(got, want) {
		t.Errorf("Delete() = %v, want %v", got, want)
	}
}

func TestSetDiff(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("b")},
			want: Of("a", "c"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s1.Diff(tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSymmetricDiff(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "intersect",
			args: args{s1: Of("a", "b", "c"), s2: Of("b")},
			want: Of("a", "c"),
		},
		{
			name: "symmetric diff",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("a", "c", "f", "d"),
		},
		{
			name: "diff",
			args: args{s1: Of("a", "b", "c"), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of("a", "b", "c"),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s1.SymmetricDiff(tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SymmetricDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUnion(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("a", "b", "c", "f", "d"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of("a", "b", "c"),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s1.Union(tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIntersect(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("b"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of[string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s1.Intersect(tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDiff(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("b")},
			want: Of("a", "c"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSymmetricDiff(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "intersect",
			args: args{s1: Of("a", "b", "c"), s2: Of("b")},
			want: Of("a", "c"),
		},
		{
			name: "symmetric diff",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("a", "c", "f", "d"),
		},
		{
			name: "diff",
			args: args{s1: Of("a", "b", "c"), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of("a", "b", "c"),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SymmetricDiff(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SymmetricDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("a", "b", "c", "f", "d"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of("a", "b", "c"),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Union(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	type args struct {
		s1 Set[string]
		s2 Set[string]
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "with members",
			args: args{s1: Of("a", "b", "c"), s2: Of("f", "b", "d")},
			want: Of("b"),
		},
		{
			name: "first empty",
			args: args{s1: Of[string](), s2: Of("a", "b", "c")},
			want: Of[string](),
		},
		{
			name: "second empty",
			args: args{s1: Of("a", "b", "c"), s2: Of[string]()},
			want: Of[string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.args.s1, tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSet_String() {
	fmt.Println(Of[int]())
	fmt.Println(Of("a"))
	// Output:
	// []
	// [a]
}

func ExampleSetGoString() {
	fmt.Printf("%#v", Of("a"))
	// Output: ["a"]
}

func ExampleSet_Equal() {
	fmt.Println(Of(1, 2, 3).Equal(Of(1, 2, 3)))
	fmt.Println(Of[int]().Equal(Of[int]()))
	fmt.Println(Of(1, 2, 3).Equal(Of(1, 5, 3)))
	fmt.Println(Of(1, 2, 3).Equal(Of(1, 3)))
	// Output:
	// true
	// true
	// false
	// false
}

func ExampleSet_Each() {
	Of(1).Each(func(m int) { fmt.Println(m) })
	// Output:
	// 1
}

func ExampleSet_Members() {
	m := Of(3, 1, 2).Members()
	sort.Ints(m)
	fmt.Println(m)
	// Output:
	// [1 2 3]
}

func ExampleSet_MarshalJSON() {
	b, _ := json.Marshal(Of("a"))
	fmt.Println(string(b))
	// Output: ["a"]
}

func ExampleSet_UnmarshalJSON() {
	var s Set[int]
	_ = json.Unmarshal([]byte("[2, 1, 3]"), &s)
	fmt.Println(Of(2, 1, 3).Equal(s))
	// Output: true
}

func ExampleSet_UnmarshalJSON_Error() {
	var s Set[int]
	err := json.Unmarshal([]byte("2, 1, 3]"), &s)
	fmt.Println(err != nil)
	// Output: true
}

func ExampleSet_MarshalText() {
	b, _ := Of("a").MarshalText()
	fmt.Println(string(b))
	// Output: ["a"]
}

func ExampleSet_UnmarshalText() {
	var s Set[int]
	_ = (&s).UnmarshalText([]byte("[2, 1, 3]"))
	fmt.Println(Of(2, 1, 3).Equal(s))
	// Output: true
}

func ExampleSet_UnmarshalText_Error() {
	var s Set[int]
	err := (&s).UnmarshalText([]byte("2, 1, 3]"))
	fmt.Println(err != nil)
	// Output: true
}
