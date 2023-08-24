package oneof

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type Doer interface {
	Value
	Do() error
}

type DoMe struct {
	You string `json:"you,omitempty"`
}

func (d *DoMe) Do() error {
	return nil
}

func (*DoMe) GetName() string {
	return "me"
}

type DoYou struct {
	Me string `json:"me,omitempty"`
}

func (d *DoYou) Do() error {
	return nil
}

func (*DoYou) GetName() string {
	return "you"
}

func DoerNew(name string) (Doer, error) {
	switch {
	case (*DoMe)(nil).GetName() == name:
		return new(DoMe), nil
	case (*DoYou)(nil).GetName() == name:
		return new(DoYou), nil
	}
	return nil, fmt.Errorf("unknown type name %s", name)
}

type DoerFactory struct{}

func (f DoerFactory) New(name string) (Doer, error) {
	switch {
	case (*DoMe)(nil).GetName() == name:
		return new(DoMe), nil
	case (*DoYou)(nil).GetName() == name:
		return new(DoYou), nil
	}
	return nil, fmt.Errorf("unknown type name %s", name)
}

type AnyDoer = OneOf[Doer, DoerFactory]

func TestOneOf_MarshalJSON(t *testing.T) {
	t.Parallel()

	s := struct {
		Doer AnyDoer `json:"doer,omitempty"`
	}{}

	s.Doer.Set(&DoMe{You: "me"})

	b, err := json.Marshal(s)
	if err != nil {
		t.Error("marshalling OneOf", err)
	}

	want := `{"doer":{"me":{"you":"me"}}}`
	if got := string(b); got != want {
		t.Errorf("want '%s' != got '%s'", want, got)
	}

	s.Doer.Set(&DoYou{Me: "you"})

	b, err = json.Marshal(s)
	if err != nil {
		t.Error("marshalling OneOf", err)
	}

	want = `{"doer":{"you":{"me":"you"}}}`
	if got := string(b); got != want {
		t.Errorf("want '%s' != got '%s'", want, got)
	}
}

func TestOneOf_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	s := struct {
		Doer AnyDoer `json:"doer,omitempty"`
	}{}
	var want Doer

	err := json.Unmarshal([]byte(`{"doer":{"me":{"you":"me"}}}`), &s)
	if err != nil {
		t.Error("unmarshalling OneOf", err)
	}

	want = &DoMe{You: "me"}
	got := s.Doer.Get()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want '%#v' != got '%#v'", want, got)
	}

	err = json.Unmarshal([]byte(`{"doer":{"you":{"me":"you"}}}`), &s)
	if err != nil {
		t.Error("unmarshalling OneOf", err)
	}

	want = &DoYou{Me: "you"}
	got = s.Doer.Get()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want '%#v' != got '%#v'", want, got)
	}
}

func TestOneOf_MarshalJSON_Empty(t *testing.T) {
	t.Parallel()

	s := struct {
		Doer AnyDoer `json:"doer,omitempty"`
	}{}

	b, err := json.Marshal(s)
	if err != nil {
		t.Error("marshalling OneOf", err)
	}

	want := `{"doer":{}}`
	if got := string(b); got != want {
		t.Errorf("want '%s' != got '%s'", want, got)
	}

	err = json.Unmarshal(b, &s)
	if err != nil && !errors.Is(err, ErrEmpty) {
		t.Error("unmarshalling OneOf", err)
	}
	if err == nil {
		t.Error("expected error")
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("want panic on empty value")
			}
		}()
		_ = s.Doer.Get().Do()
	}()
}

type DoerFactoryWithEmpty struct{}

func (f DoerFactoryWithEmpty) New(name string) (Doer, error) {
	switch {
	case (*DoMe)(nil).GetName() == name:
		return new(DoMe), nil
	case (*DoYou)(nil).GetName() == name:
		return new(DoYou), nil
	case "" == name:
		// The value will be nil instead of returning an error when empty JSON is unmarshaled.
		return nil, nil
	}
	return nil, fmt.Errorf("unknown type name %s", name)
}

func TestOneOf_MarshalJSON_EmptyNoErr(t *testing.T) {
	t.Parallel()

	type AnyDoer = OneOf[Doer, DoerFactoryWithEmpty]

	s := struct {
		Doer AnyDoer `json:"doer,omitempty"`
	}{
		Doer: AnyDoer{Value: &DoMe{You: "me"}},
	}

	b := []byte(`{"doer":{}}`)
	err := json.Unmarshal(b, &s)
	if err != nil {
		t.Error("unmarshalling OneOf", err)
	}
}
