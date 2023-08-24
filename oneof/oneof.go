package oneof

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gotidy/lib/types"
)

// ErrEmpty returned by the unmarshaler when data is empty '{}'.
var ErrEmpty = errors.New("empty")

// Value is the interface that must be implemented by any OneOf value.
type Value interface {
	GetName() string
}

// Factory describe objects factory.
type Factory[V Value] interface {
	New(name string) (V, error)
}

// OneOf is container of many values.
type OneOf[V Value, R Factory[V]] struct {
	factory R
	Value   V
}

// Get the value.
func (o *OneOf[V, R]) Get() V {
	return o.Value
}

// Set the value.
func (o *OneOf[V, R]) Set(v V) {
	o.Value = v
}

// UnmarshalJSON unmarshal OneOf value.
func (o *OneOf[V, R]) UnmarshalJSON(b []byte) error {
	v := make(map[string]json.RawMessage, 1)
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("unmarshalling sub object: %w", err)
	}
	if l := len(v); l > 1 {
		return fmt.Errorf("expected a one field, but contains fields %d", l)
	}

	var name string
	var raw json.RawMessage
	for k, v := range v {
		name = k
		raw = v
		break
	}

	value, err := o.factory.New(name)
	switch {
	case err != nil && len(v) == 0:
		return ErrEmpty // Drop factory error.
	case err != nil:
		return fmt.Errorf("getting object '%s': %w", name, err)
	case len(v) == 0:
		var zero V
		o.Value = zero
		return nil
	}

	err = json.Unmarshal(raw, value)

	if o.Value = value; err != nil {
		return fmt.Errorf("unmarshalling '%s': %w", name, err)
	}

	return nil
}

// MarshalJSON marshals OneOf value.
func (o OneOf[V, R]) MarshalJSON() ([]byte, error) {
	if types.IsNil(o.Value) {
		return []byte("{}"), nil
	}

	v := map[string]V{o.Value.GetName(): o.Value}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshalling '%s': %w", o.Value.GetName(), err)
	}
	return b, nil
}
