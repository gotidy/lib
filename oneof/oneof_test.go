package oneof

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var ErrOps = errors.New("ops")

type Doer interface {
	Value
	Do() error
}

type DoMe struct {
	You string `json:"you,omitempty"`
}

func (d *DoMe) Do() error {
	return ErrOps
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

func TestOneOf_UnmarshalJSON_DoerFactory(t *testing.T) {
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
	if wantErr, gotErr := want.Do(), got.Do(); wantErr != gotErr {
		t.Errorf("want.Do() '%#v' != got.Do() '%#v'", wantErr, gotErr)
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
	if wantErr, gotErr := want.Do(), got.Do(); wantErr != gotErr {
		t.Errorf("want.Do() '%#v' != got.Do() '%#v'", wantErr, gotErr)
	}
}
