// Model for attributes within a SCIM schema
package resource

import (
	"reflect"
	"strings"
)

type AttributeGetter interface {
	GetAttribute(path string) *Attribute
}

type Attribute struct {
	Name            string       `json:"name"`
	Type            string       `json:"type"`
	SubAttributes   []*Attribute `json:"subAttributes"`
	MultiValued     bool         `json:"multiValued"`
	Description     string       `json:"description"`
	Required        bool         `json:"required"`
	CanonicalValues []string     `json:"canonicalValues"`
	CaseExact       bool         `json:"caseExact"`
	Mutability      string       `json:"mutability"`
	Returned        string       `json:"returned"`
	Uniqueness      string       `json:"uniqueness"`
	ReferenceTypes  []string     `json:"referenceTypes"`
	Assist          *Assist      `json:"_assist"`
}

func (a *Attribute) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"name":            a.Name,
		"type":            a.Type,
		"subAttributes":   make([]map[string]interface{}, 0, len(a.SubAttributes)),
		"multiValued":     a.MultiValued,
		"description":     a.Description,
		"required":        a.Required,
		"canonicalValues": a.CanonicalValues,
		"caseExact":       a.CaseExact,
		"mutability":      a.Mutability,
		"returned":        a.Returned,
		"uniqueness":      a.Uniqueness,
		"referenceTypes":  a.ReferenceTypes,
	}
	for _, subAttr := range a.SubAttributes {
		data["subAttributes"] = append(data["subAttributes"].([]map[string]interface{}), subAttr.ToMap())
	}
	return data
}

func (a *Attribute) GetAttribute(path string) *Attribute {
	for _, attr := range a.SubAttributes {
		if strings.ToLower(attr.Name) == strings.ToLower(path) {
			return attr
		}
	}
	return nil
}

func (a *Attribute) IsValueAssigned(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}

	if v.Kind() == reflect.Interface {
		return a.IsValueAssigned(v.Elem())
	}

	if a.MultiValued {
		switch {
		case v.IsNil():
			return false
		case v.Kind() != reflect.Slice && v.Kind() != reflect.Array:
			return false
		default:
			return v.Len() > 0
		}
	} else {
		switch a.Type {
		case Complex:
			switch {
			case v.IsNil():
				return false
			case v.Kind() != reflect.Map:
				return false
			default:
				return v.Len() > 0
			}
		default:
			return true
		}
	}
}

func (a *Attribute) IsUnassigned(object interface{}) bool {
	if a.MultiValued {
		if nil == object {
			return true
		} else if value := reflect.ValueOf(object); value.Kind() == reflect.Slice {
			return value.Len() == 0
		} else {
			return false
		}
	}

	switch a.Type {
	case String, Reference, DateTime, Binary:
		return nil == object || "" == object
	case Integer:
		return nil == object
	case Decimal:
		return nil == object
	case Complex:
		if nil == object {
			return true
		} else if m, ok := object.(map[string]interface{}); ok {
			return len(m) == 0
		} else {
			return false
		}
	default:
		return false
	}
}

func (a *Attribute) Clone() *Attribute {
	cloned := &Attribute{
		Name:            a.Name,
		Type:            a.Type,
		MultiValued:     a.MultiValued,
		Description:     a.Description,
		Required:        a.Required,
		CanonicalValues: a.CanonicalValues,
		CaseExact:       a.CaseExact,
		Mutability:      a.Mutability,
		Returned:        a.Returned,
		Uniqueness:      a.Uniqueness,
		ReferenceTypes:  a.ReferenceTypes,
		Assist:          a.Assist,
		SubAttributes:   make([]*Attribute, 0),
	}
	for _, subAttr := range a.SubAttributes {
		cloned.SubAttributes = append(cloned.SubAttributes, subAttr.Clone())
	}
	return cloned
}
