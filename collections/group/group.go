// Realize `Group` type.
package group

// Group.
type Group[G comparable, T any] map[G][]T

// New creates a new Group.
func New[G comparable, T any]() Group[G, T] {
	return make(Group[G, T])
}

// Slice creates a new Group from slice.
func Slice[G comparable, T any](items []T, group func(T) G) Group[G, T] {
	result := make(Group[G, T])
	for _, v := range items {
		result.Add(group(v), v)
	}
	return result
}

// Map convert the map to the Group.
func Map[G comparable, T any](m map[G][]T) Group[G, T] {
	if m == nil {
		return make(Group[G, T])
	}
	return Group[G, T](m)
}

// Len of group.
func (g Group[G, T]) Len() int { return len(g) }

// Count of all elements.
func (g Group[G, T]) Count() int {
	var count int
	for _, s := range g {
		count += len(s)
	}
	return count
}

// Empty checks that the group is empty.
func (g Group[G, T]) Empty() bool {
	for _, s := range g {
		if len(s) != 0 {
			return false
		}
	}
	return true
}

// EachGroup iterates through groups.
func (g Group[G, T]) EachGroup(f func(group G, items []T)) {
	for k, v := range g {
		f(k, v)
	}
}

// EachItem iterates through all items.
func (g Group[G, T]) EachItem(f func(group G, item T)) {
	for group, items := range g {
		for _, item := range items {
			f(group, item)
		}
	}
}

// Groups returns the slice of groups keys.
func (g Group[G, T]) Groups() []G {
	result := make([]G, 0, len(g))
	for m := range g {
		result = append(result, m)
	}
	return result
}

// Groups returns items of the group.
func (g Group[G, T]) Group(group G) []T {
	return g[group]
}

// Add item to the group.
func (g Group[G, T]) Add(group G, v ...T) Group[G, T] {
	g[group] = append(g[group], v...)
	return g
}

// Delete specified groups.
func (g Group[G, T]) Delete(groups ...G) Group[G, T] {
	for _, group := range groups {
		delete(g, group)
	}
	return g
}

// Diff removes groups existed in the passed group.
func (g Group[G, T]) Diff(group Group[G, T]) Group[G, T] {
	for key := range group {
		delete(g, key)
	}
	return g
}

// Union groups.
func (g Group[G, T]) Union(group Group[G, T]) Group[G, T] {
	for key, items := range group {
		g[key] = append(g[key], items...)
	}
	return g
}

// Clone group.
func (g Group[G, T]) Clone() Group[G, T] {
	result := make(Group[G, T], len(g))
	for k, v := range g {
		result[k] = append(([]T(nil)), v...)
	}
	return result
}

// Has group.
func (g Group[G, T]) Has(group G) bool {
	_, exists := g[group]
	return exists
}

// Diff returns groups with items of g1 and without groups contained in g2.
func Diff[G comparable, T any](g1, g2 Group[G, T]) Group[G, T] {
	result := New[G, T]()
	for k, v := range g1 {
		if !g2.Has(k) {
			result[k] = v
		}
	}
	return result
}

// Union groups items.
func Union[G comparable, T any](g1, g2 Group[G, T]) Group[G, T] {
	result := New[G, T]()
	for k, v := range g1 {
		result[k] = v
	}
	for k, v := range g2 {
		result[k] = append(result[k], v...)
	}
	return result
}
