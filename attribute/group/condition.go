package group

import (
	"github.com/gogo/protobuf/types"
	"github.com/gunsluo/go-example/attribute/acpb"
	"github.com/gunsluo/go-example/attribute/attr"
	"github.com/pkg/errors"
)

type Attribute struct {
	Name     string      `json:"name,omitempty"`
	Type     string      `json:"type,omitempty"`
	Required bool        `json:"required,omitempty"`
	Value    interface{} `json:"default,omitempty"`
}

type ConditionOption struct {
	Expression string       `json:"expression,omitempty"`
	Attributes []*Attribute `json:"attributes,omitempty"`
}

type Condition struct {
	Name    string           `json:"name,omitempty"`
	Type    string           `json:"type,omitempty"`
	Options *ConditionOption `json:"options,omitempty"`
}

type Conditions []*Condition

type PredefinedPolicy struct {
	Name        string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string     `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Resources   []string   `protobuf:"bytes,3,rep,name=resources,proto3" json:"resources,omitempty"`
	Actions     []string   `protobuf:"bytes,4,rep,name=actions,proto3" json:"actions,omitempty"`
	Conditions  Conditions `protobuf:"bytes,5,opt,name=conditions,proto3" json:"conditions,omitempty"`
}

func ConvertCondition(c *acpb.Condition) (*Condition, error) {
	nc := &Condition{
		Name: c.Name,
		Type: c.Type,
	}

	if c.Options != nil {
		nc.Options = &ConditionOption{
			Expression: c.Options.Expression,
		}

		attributes, err := convertAttributes(c.Options.Attributes)
		if err != nil {
			return nil, err
		}
		nc.Options.Attributes = attributes
	}

	return nc, nil
}

func ConvertConditions(cs []*acpb.Condition) (Conditions, error) {
	var ncs Conditions
	for _, c := range cs {
		nc, err := ConvertCondition(c)
		if err != nil {
			return nil, err
		}
		ncs = append(ncs, nc)
	}

	return ncs, nil
}

func ConvertPrettyCondition(c *Condition) (*acpb.Condition, error) {
	condition := &acpb.Condition{
		Name: c.Name,
		Type: c.Type,
	}

	if c.Options != nil {
		var attributes []*acpb.Attribute
		for _, a := range c.Options.Attributes {
			av, err := convertInterfaceValue(a.Value)
			if err != nil {
				return nil, err
			}

			attributes = append(attributes, &acpb.Attribute{
				Name:     a.Name,
				Required: a.Required,
				Value:    av,
			})
		}

		condition.Options = &acpb.ConditionOption{
			Expression: c.Options.Expression,
			Attributes: attributes,
		}
	}

	return condition, nil
}

func ConvertPrettyConditions(cs Conditions) ([]*acpb.Condition, error) {
	var conditions []*acpb.Condition
	for _, c := range cs {
		nc, err := ConvertPrettyCondition(c)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, nc)
	}

	return conditions, nil
}

func (c *Condition) ConvertCondition(values map[string]interface{}) (attr.Condition, error) {
	fn, ok := attr.ConditionFactories[c.Type]
	if !ok {
		return nil, errors.Errorf("Condition %s not found", c.Type)
	}

	nc := fn()
	if c.Options != nil {
		nvs := make(map[string]interface{})
		for _, a := range c.Options.Attributes {
			if v, ok := values[a.Name]; ok {
				nvs[a.Name] = v
			} else {
				if a.Required {
					return nil, errors.Errorf("attribute %s is required", a.Name)
				}
				if a.Value != nil {
					nvs[a.Name] = a.Value
				}
			}
		}

		nc.Values(c.Options.Expression, values)
	}

	return nc, nil
}

func (cs Conditions) ConvertConditions(all map[string]map[string]interface{}) (attr.Conditions, error) {
	ncs := attr.Conditions{}
	for _, c := range cs {
		var values map[string]interface{}
		if v, ok := all[c.Name]; ok {
			values = v
		} else {
			values = make(map[string]interface{})
		}

		nc, err := c.ConvertCondition(values)
		if err != nil {
			return nil, err
		}
		ncs[c.Name] = nc
	}

	return ncs, nil
}

func ConvertAttributes(attributesTable map[string]*acpb.PolicyDTO_Attributes) (map[string]map[string]interface{}, error) {
	all := make(map[string]map[string]interface{})
	for k, v := range attributesTable {
		values := make(map[string]interface{})
		for name, av := range v.Values {
			val, _, err := convertAttributeValue(av)
			if err != nil {
				return nil, err
			}
			values[name] = val
		}

		all[k] = values
	}

	return all, nil
}

func convertAttributes(attributes []*acpb.Attribute) ([]*Attribute, error) {
	var attrs []*Attribute
	for _, a := range attributes {
		na := &Attribute{
			Name:     a.Name,
			Required: a.Required,
		}

		if a.Value == nil {
			continue
		}

		value, typ, err := convertAttributeValue(a.Value)
		if err != nil {
			return nil, err
		}

		na.Value = value
		na.Type = typ
		attrs = append(attrs, na)
	}

	return attrs, nil
}

func convertAttributeValue(v *acpb.AttributeValue) (interface{}, string, error) {
	var value interface{}
	var typ string

	switch v.Type {
	case acpb.ATTRIBUTE_TYPE_STRING:
		if v.Value != nil {
			val := &acpb.StringAttributeValue{}
			err := types.UnmarshalAny(v.Value, val)
			if err != nil {
				return value, typ, err
			}
			value = val.Value
		}
		typ = "string"
	case acpb.ATTRIBUTE_TYPE_NUMBER:
		if v.Value != nil {
			val := &acpb.NumberAttributeValue{}
			err := types.UnmarshalAny(v.Value, val)
			if err != nil {
				return value, typ, err
			}
			value = val.Value
		}
		typ = "number"
	case acpb.ATTRIBUTE_TYPE_BOOLEAN:
		if v.Value != nil {
			val := &acpb.BooleanAttributeValue{}
			err := types.UnmarshalAny(v.Value, val)
			if err != nil {
				return value, typ, err
			}
			value = val.Value
		}
		typ = "boolean"
	default:
		return value, typ, errors.New("not supported type")
	}

	return value, typ, nil
}

func convertInterfaceValue(v interface{}) (*acpb.AttributeValue, error) {
	value := &acpb.AttributeValue{}
	switch val := v.(type) {
	case string:
		any, err := types.MarshalAny(&acpb.StringAttributeValue{Value: val})
		if err != nil {
			return nil, err
		}

		value.Type = acpb.ATTRIBUTE_TYPE_STRING
		value.Value = any
	case int:
		any, err := types.MarshalAny(&acpb.NumberAttributeValue{Value: int64(val)})
		if err != nil {
			return nil, err
		}

		value.Type = acpb.ATTRIBUTE_TYPE_NUMBER
		value.Value = any
	case int32:
		any, err := types.MarshalAny(&acpb.NumberAttributeValue{Value: int64(val)})
		if err != nil {
			return nil, err
		}

		value.Type = acpb.ATTRIBUTE_TYPE_NUMBER
		value.Value = any
	case int64:
		any, err := types.MarshalAny(&acpb.NumberAttributeValue{Value: val})
		if err != nil {
			return nil, err
		}

		value.Type = acpb.ATTRIBUTE_TYPE_NUMBER
		value.Value = any
	case bool:
		any, err := types.MarshalAny(&acpb.BooleanAttributeValue{Value: val})
		if err != nil {
			return nil, err
		}

		value.Type = acpb.ATTRIBUTE_TYPE_BOOLEAN
		value.Value = any
	default:
		return value, errors.Errorf("%v not supported type", v)
	}

	return value, nil
}

//AttributeValues map[string]*PolicyDTO_Attributes `protobuf:"bytes,8,rep,name=attribute_values,json=attributeValues,proto3" json:"attribute_values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
