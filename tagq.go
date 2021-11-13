package tagq

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Common query errors
var (
	ErrValueIsNil                = fmt.Errorf("value is nil")
	ErrValueCannotBeCasted       = fmt.Errorf("value cannot be casted")
	ErrStrOfNil                  = fmt.Errorf("Str(): %w", ErrValueIsNil)
	ErrIntOfNil                  = fmt.Errorf("Int(): %w", ErrValueIsNil)
	ErrIntOfIncompatibleType     = fmt.Errorf("Int(): %w to int", ErrValueCannotBeCasted)
	ErrFloat64OfNil              = fmt.Errorf("Float64(): %w", ErrValueIsNil)
	ErrFloat64OfIncompatibleType = fmt.Errorf("Float64(): %w to float64", ErrValueCannotBeCasted)
	ErrTimeOfNil                 = fmt.Errorf("Time(): %w", ErrValueIsNil)
	ErrTimeOfIncompatibleType    = fmt.Errorf("Time(): %w to time.Time", ErrValueCannotBeCasted)
	ErrQueryNilValue             = fmt.Errorf("Q(): %w", ErrValueIsNil)
	ErrQueryEmptySlice           = fmt.Errorf("Q(): %w (empty slice)", ErrValueIsNil)
	ErrQueryEmptyMap             = fmt.Errorf("Q(): %w (empty map)", ErrValueIsNil)
	ErrNonIntegerIndexForSlice   = fmt.Errorf("index query must be an integer (or last; random; rand) while querying a slice")
	ErrIndexOutOfBoundsForSlice  = fmt.Errorf("index query is out of bounds while querying a slice")
	ErrInvalidKeyTypeForMapQuery = fmt.Errorf("invalid key type for map query")
	ErrQueryEmptyStruct          = fmt.Errorf("Q(): %w (empty struct)", ErrValueIsNil)
	ErrStructFieldNotFound       = fmt.Errorf("Q(): %w (struct field not found)", ErrValueIsNil)
	ErrUnsupportedTypeForQuery   = fmt.Errorf("Q(): unsupported type for query")
)

// Value is a queryable interface. It traverses strucrs, maps and slices using reflection.
// It will also check the struct tags defined in the Tags() function.
// The tags can be overriden by altering the global slice DefaultTags or by calling the SetTags() function.
type Value interface {
	// Returns the string representation of the value.
	// If the value is already a string, the value is returned as is.
	Str() string
	// Returns the int representation of the value.
	// If the value is already an int, the value is returned as is.
	Int() int
	// Returns the float64 representation of the value.
	// If the value is already a float64, the value is returned as is.
	Float64() float64
	// Returns the float64 representation of the value.
	// If the value is already a float64, the value is returned as is.
	F64() float64
	// Returns the time representation of the value.
	// If the value is already a time.Time, the value is returned as is.
	Time() time.Time
	// Returns the raw value.
	Interface() interface{}
	// Q performs a query to fetch a value from a struct, map or slice.
	Q(query ...string) Value
	// Return the last error
	Err() error
	// Tags returns the struct tags that will be searched against the query.
	Tags() []string
	// SetTags sets the struct tags that will be searched against the query.
	SetTags(tags []string) Value
}

type value struct {
	v       interface{}
	lastErr error
	tags    []string
}

func (v *value) Str() string {
	if v.v == nil {
		v.lastErr = ErrStrOfNil
		return ""
	}
	if str, ok := v.v.(string); ok {
		v.lastErr = nil
		return str
	}
	if x, ok := v.v.(fmt.Stringer); ok {
		v.lastErr = nil
		return x.String()
	}
	v.lastErr = nil
	return fmt.Sprintf("%v", v.v)
}

func (v *value) Int() int {
	if v.v == nil {
		v.lastErr = ErrIntOfNil
		return 0
	}
	switch v.v.(type) {
	case int:
		v.lastErr = nil
		return v.v.(int)
	case int8:
		v.lastErr = nil
		return int(v.v.(int8))
	case int16:
		v.lastErr = nil
		return int(v.v.(int16))
	case int32:
		v.lastErr = nil
		return int(v.v.(int32))
	case int64:
		v.lastErr = nil
		return int(v.v.(int64))
	case uint:
		v.lastErr = nil
		return int(v.v.(uint))
	case uint8:
		v.lastErr = nil
		return int(v.v.(uint8))
	case uint16:
		v.lastErr = nil
		return int(v.v.(uint16))
	case uint32:
		v.lastErr = nil
		return int(v.v.(uint32))
	case uint64:
		v.lastErr = nil
		return int(v.v.(uint64))
	case float32:
		v.lastErr = nil
		return int(v.v.(float32))
	case float64:
		v.lastErr = nil
		return int(v.v.(float64))
	}
	v.lastErr = nil
	return 0
}

func (v *value) Float64() float64 {
	if v.v == nil {
		v.lastErr = ErrIntOfNil
		return 0
	}
	switch v.v.(type) {
	case int:
		v.lastErr = nil
		return float64(v.v.(int))
	case int8:
		v.lastErr = nil
		return float64(v.v.(int8))
	case int16:
		v.lastErr = nil
		return float64(v.v.(int16))
	case int32:
		v.lastErr = nil
		return float64(v.v.(int32))
	case int64:
		v.lastErr = nil
		return float64(v.v.(int64))
	case uint:
		v.lastErr = nil
		return float64(v.v.(uint))
	case uint8:
		v.lastErr = nil
		return float64(v.v.(uint8))
	case uint16:
		v.lastErr = nil
		return float64(v.v.(uint16))
	case uint32:
		v.lastErr = nil
		return float64(v.v.(uint32))
	case uint64:
		v.lastErr = nil
		return float64(v.v.(uint64))
	case float32:
		v.lastErr = nil
		return float64(v.v.(float32))
	case float64:
		v.lastErr = nil
		return v.v.(float64)
	}
	v.lastErr = nil
	return 0
}

func (v *value) F64() float64 {
	return v.Float64()
}

func (v *value) Time() time.Time {
	if v.v == nil {
		v.lastErr = ErrTimeOfNil
		return time.Time{}
	}
	switch v.v.(type) {
	case time.Time:
		v.lastErr = nil
		return v.v.(time.Time)
	case *time.Time:
		v.lastErr = nil
		return *v.v.(*time.Time)
	case string:
		return v.timeFromStr(v.v.(string))
	case *string:
		return v.timeFromStr(*v.v.(*string))
	case int64:
		return v.timeFromInt64(v.v.(int64))
	case *int64:
		return v.timeFromInt64(*v.v.(*int64))
	case int32:
		return v.timeFromInt64(int64(v.v.(int32)))
	case *int32:
		return v.timeFromInt64(int64(*v.v.(*int32)))
	case int:
		return v.timeFromInt64(int64(v.v.(int)))
	case *int:
		return v.timeFromInt64(int64(*v.v.(*int)))
	case float64:
		return v.timeFromInt64(int64(v.v.(float64)))
	case *float64:
		return v.timeFromInt64(int64(*v.v.(*float64)))
	case float32:
		return v.timeFromInt64(int64(v.v.(float32)))
	case *float32:
		return v.timeFromInt64(int64(*v.v.(*float32)))
	}
	v.lastErr = ErrTimeOfIncompatibleType
	return time.Time{}
}

func (v *value) Interface() interface{} {
	return v.v
}

func (v *value) timeFromStr(vs string) time.Time {
	if tx, err := time.Parse(time.RFC3339Nano, vs); err == nil {
		v.lastErr = nil
		return tx
	}
	if tx, err := time.Parse(time.RFC3339, vs); err == nil {
		v.lastErr = nil
		return tx
	}
	// Datetime format
	if tx, err := time.Parse("2006-01-02 15:04:05", vs); err == nil {
		v.lastErr = nil
		return tx
	}
	v.lastErr = ErrTimeOfIncompatibleType
	return time.Time{}
}

func (v *value) timeFromInt64(i int64) time.Time {
	if i > 9999999999 {
		// Epoch with millisecond precision
		return time.Unix(i/1000, i%1000*1000000)
	}
	return time.Unix(i, 0)
}

func (v *value) Q(query ...string) Value {
	if len(query) == 0 {
		return v
	}
	if v.v == nil {
		v.lastErr = ErrQueryNilValue
		return v
	}
	return v.qnext(reflect.ValueOf(v.v), query...)
}

func (v *value) new(rv reflect.Value) *value {
	var vi interface{}
	if rv.IsValid() {
		if rv.CanInterface() {
			vi = rv.Interface()
		}
	}
	return &value{
		lastErr: nil,
		v:       vi,
		tags:    v.tags,
	}
}

func (v *value) qnext(rv reflect.Value, query ...string) *value {
	if len(query) == 0 {
		return v
	}
	t := rv.Type()
	if t.Kind() == reflect.Ptr {
		return v.qnext(rv.Elem(), query...)
	}
	switch t.Kind() {
	case reflect.Slice:
		if rv.Len() == 0 {
			vn := v.new(reflect.ValueOf(nil))
			vn.lastErr = ErrQueryNilValue
			return vn
		}
		if query[0] == "last" {
			vn := v.new(rv.Index(rv.Len() - 1))
			if len(query) == 1 {
				return vn
			}
			return vn.qnext(rv.Index(rv.Len()-1), query[1:]...)
		} else if query[0] == "random" || query[0] == "rand" {
			n := rand.Intn(rv.Len())
			vn := v.new(rv.Index(n))
			if len(query) == 1 {
				return vn
			}
			return vn.qnext(rv.Index(n), query[1:]...)
		}
		vi, err := strconv.Atoi(query[0])
		if err != nil {
			vn := v.new(reflect.ValueOf(nil))
			vn.lastErr = ErrNonIntegerIndexForSlice
			return vn
		}
		if vi < 0 || vi >= rv.Len() {
			vn := v.new(reflect.ValueOf(nil))
			vn.lastErr = ErrIndexOutOfBoundsForSlice
			return vn
		}
		vn := v.new(rv.Index(vi))
		if len(query) == 1 {
			return vn
		}
		return vn.qnext(rv.Index(vi), query[1:]...)
	case reflect.Map:
		if rv.Len() == 0 {
			vn := v.new(reflect.ValueOf(nil))
			vn.lastErr = ErrQueryEmptyMap
			return vn
		}
		// try to get a velue by key like rv.MapIndex(reflect.ValueOf(query[0]))
		keytype := rv.Type().Key()
		switch keytype.Kind() {
		case reflect.String:
			item := rv.MapIndex(reflect.ValueOf(query[0]))
			vn := v.new(item)
			if len(query) == 1 {
				return vn
			}
			return vn.qnext(item, query[1:]...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sconv, err := strconv.Atoi(query[0])
			if err != nil {
				vn := v.new(reflect.ValueOf(nil))
				vn.lastErr = ErrInvalidKeyTypeForMapQuery
				return vn
			}
			var vx reflect.Value
			switch keytype.Kind() {
			case reflect.Int:
				vx = reflect.ValueOf(sconv)
			case reflect.Int8:
				vx = reflect.ValueOf(int8(sconv))
			case reflect.Int16:
				vx = reflect.ValueOf(int16(sconv))
			case reflect.Int32:
				vx = reflect.ValueOf(int32(sconv))
			case reflect.Int64:
				vx = reflect.ValueOf(int64(sconv))
			case reflect.Uint:
				vx = reflect.ValueOf(uint(sconv))
			case reflect.Uint8:
				vx = reflect.ValueOf(uint8(sconv))
			case reflect.Uint16:
				vx = reflect.ValueOf(uint16(sconv))
			case reflect.Uint32:
				vx = reflect.ValueOf(uint32(sconv))
			case reflect.Uint64:
				vx = reflect.ValueOf(uint64(sconv))
			}
			item := rv.MapIndex(vx)
			vn := v.new(item)
			if len(query) == 1 {
				return vn
			}
			return vn.qnext(item, query[1:]...)
		}
	case reflect.Struct:
		if rv.NumField() == 0 {
			vn := v.new(reflect.ValueOf(nil))
			vn.lastErr = ErrQueryEmptyStruct
			return vn
		}
		rvt := rv.Type()
		if len(query[0]) > 0 && strings.ToUpper(query[0][:1]) == query[0][:1] {
			// try to get a field by name like rv.FieldByName(query[0])
			field := rv.FieldByName(query[0])
			if field.IsValid() {
				vn := v.new(rv.FieldByName(query[0]))
				if len(query) == 1 {
					return vn
				}
				return vn.qnext(rv.FieldByName(query[0]), query[1:]...)
			}
		}
		for fieldi := 0; fieldi < rv.NumField(); fieldi++ {
			field := rv.Field(fieldi)
			sfield := rvt.Field(fieldi)
			for _, tlookup := range v.tags {
				if sfield.Tag.Get(tlookup) == query[0] {
					vn := v.new(field)
					if len(query) == 1 {
						return vn
					}
					return vn.qnext(field, query[1:]...)
				}
			}
		}
		// no field found
		vn := v.new(reflect.ValueOf(nil))
		vn.lastErr = ErrStructFieldNotFound
		return vn
		//TODO: case reflect.Func
	}
	vn := v.new(reflect.ValueOf(nil))
	vn.lastErr = ErrUnsupportedTypeForQuery
	return nil
}

func (v *value) Err() error {
	return v.lastErr
}

func (v *value) Tags() []string {
	return v.tags
}

func (v *value) SetTags(tags []string) Value {
	v.tags = make([]string, len(tags))
	copy(v.tags, tags)
	return v
}

// Q works similar to jquery, but using reflection.
// It will also check the struct tags defined in the Tags() function.
// The tags can be overriden by altering the global slice DefaultTags or by calling the SetTags() function.
func Q(v interface{}, query ...string) Value {
	x := &value{
		v:       v,
		lastErr: nil,
		tags:    copyDefaultTags(),
	}
	return x.Q(query...)
}

// DefaultTags is the default slice of tags copied to a new Q instance.
var DefaultTags = []string{"json", "xml", "param", "query", "header"}

func copyDefaultTags() []string {
	tags := make([]string, len(DefaultTags))
	copy(tags, DefaultTags)
	return tags
}
