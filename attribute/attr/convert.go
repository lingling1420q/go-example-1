package attr

/*

import (
	"github.com/gogo/protobuf/types"
	"github.com/gunsluo/go-example/attribute/acpb"
	"github.com/pkg/errors"
)

func ConvertCondition(c *acpb.Condition) (Condition, error) {
	fn, ok := ConditionFactories[c.Type]
	if !ok {
		return nil, errors.Errorf("Condition %s not found", c.Type)
	}

	nc := fn()
	if c.Options != nil {
		values, err := getAttributeValues(c.Options)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		nc.Values(c.Options.Expression, values)
	}

	return nc, nil
}

func ConvertConditions(cs []*acpb.Condition) (Conditions, error) {
	ncs := Conditions{}
	for _, c := range cs {
		nc, err := ConvertCondition(c)
		if err != nil {
			return nil, err
		}
		ncs[c.Name] = nc
	}

	return ncs, nil
}

func getAttributeValues(options *acpb.ConditionOption) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	if options == nil {
		return values, nil
	}

	for _, a := range options.Attributes {
		switch a.Type {
		case acpb.ATTRIBUTE_TYPE_STRING:
			v := &acpb.StringAttributeValue{}
			err := types.UnmarshalAny(a.Default, v)
			if err != nil {
				return nil, err
			}
			values[a.Name] = v.Value
		case acpb.ATTRIBUTE_TYPE_NUMBER:
			v := &acpb.NumberAttributeValue{}
			err := types.UnmarshalAny(a.Default, v)
			if err != nil {
				return nil, err
			}
			values[a.Name] = v.Value
		case acpb.ATTRIBUTE_TYPE_BOOLEAN:
			v := &acpb.BooleanAttributeValue{}
			err := types.UnmarshalAny(a.Default, v)
			if err != nil {
				return nil, err
			}
			values[a.Name] = v.Value
		}
	}

	return values, nil
}
*/

/*
import "github.com/pkg/errors"

type Key struct {
	Name      string `json:"name"`
	Required  bool   `json:"required"`
	Condition string `json:"condition"`
}

type Value interface{}

type Attributes map[Key]Value

func (as Attributes) ConvertConditions() (Conditions, error) {
	cs := Conditions{}
	for k, v := range as {
		fn, ok := ConditionFactories[k.Condition]
		if !ok {
			return nil, errors.Errorf("CCondition %s not found", k.Condition)
		}

		c := fn()
		c.Value(v)

		cs[k.Name] = c
	}

	return cs, nil
}
*/

/*
import "github.com/pkg/errors"

type Attribute struct {
	Name      string
	Value     interface{}
	Required  bool
	Condition string
}

type Attributes []Attribute

func (a *Attribute) ConvertCondition() (Condition, error) {
	fn, ok := ConditionFactories[a.Condition]
	if !ok {
		return nil, errors.Errorf("CCondition %s not found", a.Condition)
	}

	c := fn()
	c.Value(a.Value)

	return c, nil
}

func (as Attributes) ConvertConditions() (Conditions, error) {
	cs := Conditions{}
	for _, a := range as {
		c, e := a.ConvertCondition()
		if e != nil {
			return nil, e
		}
		cs[a.Name] = c
	}

	return cs, nil
}
*/
