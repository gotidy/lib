package oneof

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
)

type Matcher interface {
	GetName() string
	Match(s string) (matched bool, err error)
}

type Mask struct {
	Pattern string `json:"pattern,omitempty"`
}

func (*Mask) GetName() string {
	return "mask"
}

func (m *Mask) Match(s string) (matched bool, err error) {
	matched, err = filepath.Match(m.Pattern, s)
	if err != nil {
		return false, fmt.Errorf("mask matching '%s' with '%s': %w", s, m.Pattern, err)
	}
	return matched, nil
}

type Regexp struct {
	Pattern string `json:"pattern,omitempty"`
}

func (*Regexp) GetName() string {
	return "regexp"
}

func (r *Regexp) Match(s string) (matched bool, err error) {
	matched, err = regexp.MatchString(r.Pattern, s)
	if err != nil {
		return false, fmt.Errorf("regexp matching '%s' with '%s': %w", s, r.Pattern, err)
	}
	return matched, nil
}

type MatcherFactory struct{}

func (f MatcherFactory) New(name string) (Matcher, error) {
	switch {
	case (*Mask)(nil).GetName() == name:
		return new(Mask), nil
	case (*Regexp)(nil).GetName() == name:
		return new(Regexp), nil
	}
	return nil, fmt.Errorf("unknown type name %s", name)
}

type AnyMatcher = OneOf[Matcher, MatcherFactory]

type Matchers struct {
	Matchers []AnyMatcher `json:"matchers,omitempty"`
}

func (m Matchers) Match(s string) (bool, error) {
	for _, matcher := range m.Matchers {
		matched, err := matcher.Get().Match(s)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func ExampleOneOf() {
	matchers := Matchers{
		Matchers: []AnyMatcher{
			{Value: &Mask{Pattern: "hell?"}},
			{Value: &Regexp{Pattern: "g.*ye"}},
		},
	}

	// Marshall.
	b, err := json.Marshal(matchers)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Clear
	matchers = Matchers{}

	// Unmarshal.
	err = json.Unmarshal(b, &matchers)
	if err != nil {
		panic(err)
	}

	// Check matchers.
	matched, err := matchers.Match("hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(matched)

	matched, err = matchers.Match("goodbye")
	if err != nil {
		panic(err)
	}
	fmt.Println(matched)

	matched, err = matchers.Match("unmatched")
	if err != nil {
		panic(err)
	}
	fmt.Println(matched)

	// Marshall unmarshalled.
	b, err = json.Marshal(matchers)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	// {"matchers":[{"mask":{"pattern":"hell?"}},{"regexp":{"pattern":"g.*ye"}}]}
	// true
	// true
	// false
	// {"matchers":[{"mask":{"pattern":"hell?"}},{"regexp":{"pattern":"g.*ye"}}]}
}
