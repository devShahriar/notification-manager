package utils

import (
	"database/sql"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Int interface {
	int | int8 | int16 | int32 | int64
}

// WPBInt32ToInt returns int value from wrapperspb int32 or nil if the given value is nil
func WPBInt32ToInt(value *wrapperspb.Int32Value) *int {
	if value == nil {
		return nil
	}
	return ToP[int](int(value.Value))
}

// ToString returns the string value from wrapperspb
func ToString(value *wrapperspb.StringValue) *string {
	if value == nil {
		return nil
	}
	return &value.Value
}

// SQLNullString returns empty sql.NullString if the string is an empty string otherwise assigning to sql.NullString
func SQLNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// PtrToString if the val is nil otherwise the string value
func PtrToString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

// ToBool returns false if the given var is nil other wise the value
func ToBool(value *bool) bool {
	return value != nil && *value
}

// PtrToInt64 returns int64 value of pointer int64
func PtrToInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func WPBToFloat64P(value *wrapperspb.DoubleValue) *float64 {
	if value == nil {
		return nil
	}
	return &value.Value
}

// |««««««««  to wrapper spb value »»»»»»»»»»|

// Float32WrapperSPBValue returns wrapped spb value of the given float32
func Float32WrapperSPBValue(value *float32) *wrapperspb.FloatValue {
	if value == nil {
		return nil
	}
	return wrapperspb.Float(*value)
}

// StringWrapperSPBValue returns the value wrapped in wrapperspb
func StringWrapperSPBValue(value *string) *wrapperspb.StringValue {
	if value == nil {
		return nil
	}
	return wrapperspb.String(*value)
}

// Uint64WrapperSPBValue returns the wrappedSPB uint64
func Uint64WrapperSPBValue(value *uint64) *wrapperspb.UInt64Value {
	if value == nil {
		return nil
	}
	return wrapperspb.UInt64(*value)
}

// WrapperSPBToBoolP returns the bool pointer from wrapper spb
func WrapperSPBToBoolP(value *wrapperspb.BoolValue) *bool {
	if value == nil {
		return nil
	}
	return &value.Value
}

// WrapperSpbToInt64P returns int64 value from wrapperspb int64 or nil if the given value is nil
func WrapperSpbToInt64P(value *wrapperspb.Int64Value) *int64 {
	if value == nil {
		return nil
	}
	return &value.Value
}

// Int64PToWrapperSpb returns wrapped spb value of the given int64
func Int64PToWrapperSpb(value *int64) *wrapperspb.Int64Value {
	if value == nil {
		return nil
	}
	return wrapperspb.Int64(*value)
}

// BoolPToWrapperSpb returns wrapped spb value of the given int64
func BoolPToWrapperSpb(value *bool) *wrapperspb.BoolValue {
	if value == nil {
		return nil
	}
	return wrapperspb.Bool(*value)
}

// IntPToWPB32 is generic function to convert any variations of int
func IntPToWPB32[T Int](value *T) *wrapperspb.Int32Value {
	if value == nil {
		return nil
	}
	return wrapperspb.Int32(int32(*value))
}

// IntPToWPB64 is generic function to convert any variations of int
func IntPToWPB64[T Int](value *T) *wrapperspb.Int64Value {
	if value == nil {
		return nil
	}
	return wrapperspb.Int64(int64(*value))
}

// ToFloat32 returns float32 value from wrapperspb float or nil if the given value is nil
func ToFloat32(value *wrapperspb.FloatValue) *float32 {
	if value == nil {
		return nil
	}
	return &value.Value
}

// Float64PtoWPB64 returns float64 wpb double value
func Float64PtoWPB64(value *float64) *wrapperspb.DoubleValue {
	if value == nil {
		return nil
	}
	return wrapperspb.Double(*value)
}

// |«««««««« to pointer functions »»»»»»»»»»|

// ToFloat32P returns pointer of given float32 value
func ToFloat32P(value float32) *float32 {
	return &value
}

// ToBoolP returns pointer of bool
func ToBoolP(value bool) *bool {
	return &value
}

// ToInt64P returns pointer of int64 value
func ToInt64P(value int64) *int64 {
	return &value
}

// ToStringP returns nil if the string is empty otherwise pointer of the string
func ToStringP(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

// ToTimeP returns pointer of time
func ToTimeP(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// ToP returns pointer of a value
func ToP[T comparable](v T) *T {
	return &v
}

// ToValue is the genetic function which returns the zero value if the given pointer is nil otherwise the actual value
func ToValue[T comparable](v *T) T {
	if v == nil {
		var noop T
		return noop
	}
	return *v
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
