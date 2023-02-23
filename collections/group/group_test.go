package group

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		m map[string][]string
	}
	tests := []struct {
		name string
		args args
		want Group[string, string]
	}{
		{
			name: "with members",
			args: args{m: map[string][]string{"foo": {"1"}, "bar": {"2"}}},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2"}},
		},
		{
			name: "nil",
			args: args{m: nil},
			want: Group[string, string]{},
		},
		{
			name: "empty",
			args: args{m: map[string][]string{}},
			want: Group[string, string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
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
		want Group[string, T]
	}{
		{
			name: "with members",
			args: args{s: []T{{Name: "foo"}, {Name: "bar"}}, f: func(v T) string { return v.Name }},
			want: Group[string, T]{"foo": {{Name: "foo"}}, "bar": {{Name: "bar"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice(tt.args.s, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	group := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2"}}
	if got := group.Len(); got != len(group) {
		t.Errorf("Len() = %v, want %v", got, len(group))
	}
}

func TestGroupCount(t *testing.T) {
	group := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}
	if got := group.Count(); got != 3 {
		t.Errorf("Len() = %v, want %v", got, 3)
	}
}

func ExampleGroup_Empty() {
	fmt.Println(Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}.Empty())
	fmt.Println(Group[string, string]{}.Empty())

	// Output:
	// false
	// true
}

func TestGroupAdd(t *testing.T) {
	group := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}
	group.Add("foo", "4")
	group.Add("boo", "5")
	expected := Group[string, string]{"foo": []string{"1", "4"}, "bar": []string{"2", "3"}, "boo": []string{"5"}}
	if !reflect.DeepEqual(group, expected) {
		t.Errorf("Add() = %v, want %v", group, expected)
	}
}

func TestGroupDelete(t *testing.T) {
	group := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}
	group.Delete("foo")
	group.Delete("boo")
	expected := Group[string, string]{"bar": []string{"2", "3"}}
	if !reflect.DeepEqual(group, expected) {
		t.Errorf("Delete() = %v, want %v", group, expected)
	}
}

func TestGroupDiff(t *testing.T) {
	type args struct {
		s1 Group[string, string]
		s2 Group[string, string]
	}
	tests := []struct {
		name string
		args args
		want Group[string, string]
	}{
		{
			name: "with members",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}},
				s2: Group[string, string]{"bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}},
		},
		{
			name: "first empty",
			args: args{
				s1: Group[string, string]{},
				s2: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{},
		},
		{
			name: "second empty",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
				s2: Group[string, string]{},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
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

func TestGroupUnion(t *testing.T) {
	type args struct {
		s1 Group[string, string]
		s2 Group[string, string]
	}
	tests := []struct {
		name string
		args args
		want Group[string, string]
	}{
		{
			name: "with members",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}},
				s2: Group[string, string]{"bar": []string{"3", "4", "5"}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3", "3", "4", "5"}, "fur": []string{"5"}},
		},
		{
			name: "first empty",
			args: args{
				s1: Group[string, string]{},
				s2: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": nil, "fur": []string{"5"}},
		},
		{
			name: "second empty",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
				s2: Group[string, string]{},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
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

func TestDiff(t *testing.T) {
	type args struct {
		s1 Group[string, string]
		s2 Group[string, string]
	}
	tests := []struct {
		name string
		args args
		want Group[string, string]
	}{
		{
			name: "with members",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}},
				s2: Group[string, string]{"bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}},
		},
		{
			name: "first empty",
			args: args{
				s1: Group[string, string]{},
				s2: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{},
		},
		{
			name: "second empty",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
				s2: Group[string, string]{},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
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

func TestUnion(t *testing.T) {
	type args struct {
		s1 Group[string, string]
		s2 Group[string, string]
	}
	tests := []struct {
		name string
		args args
		want Group[string, string]
	}{
		{
			name: "with members",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}},
				s2: Group[string, string]{"bar": []string{"3", "4", "5"}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3", "3", "4", "5"}, "fur": []string{"5"}},
		},
		{
			name: "first empty",
			args: args{
				s1: Group[string, string]{},
				s2: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": nil, "fur": []string{"5"}},
		},
		{
			name: "second empty",
			args: args{
				s1: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
				s2: Group[string, string]{},
			},
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{}, "fur": []string{"5"}},
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

func ExampleGroup_EachGroup() {
	Group[string, string]{"bar": []string{"2", "3"}}.EachGroup(func(g string, items []string) { fmt.Println(g, items) })
	// Output:
	// bar [2 3]
}

func ExampleGroup_EachItem() {
	Group[string, string]{"bar": []string{"2", "3"}}.EachItem(func(g string, i string) { fmt.Println(g, i) })
	// Output:
	// bar 2
	// bar 3
}

func ExampleGroup_Groups() {
	groups := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}.Groups()
	sort.Strings(groups)
	fmt.Println(groups)
	// Output:
	// [bar foo]
}

func ExampleGroup_Group() {
	items := Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}.Group("bar")
	fmt.Println(items)
	// Output:
	// [2 3]
}

func ExampleGroup_Count() {
	fmt.Println(Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}.Count())
	// Output:
	// 3
}

func ExampleGroup_Len() {
	fmt.Println(Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3"}}.Len())
	// Output:
	// 2
}

func TestGroupClone(t *testing.T) {
	tests := []struct {
		name string
		want Group[string, string]
	}{
		{
			name: "with members",
			want: Group[string, string]{"foo": []string{"1"}, "bar": []string{"2", "3", "3", "4", "5"}, "fur": []string{"5"}},
		},
		{
			name: "empty",
			want: Group[string, string]{"foo": []string{"1"}, "bar": nil, "fur": []string{"5"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.want.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
