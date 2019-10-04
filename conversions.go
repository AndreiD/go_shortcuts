package utils

import (
	"log"
	"io"
	"reflect"
	"strconv"
)

// GetFloat transforms an interface to float64
func GetFloat(unk interface{}) float64 {
	var floatType = reflect.TypeOf(float64(0))
	var stringType = reflect.TypeOf("")

	switch i := unk.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int64:
		return float64(i)
	case int32:
		return float64(i)
	case int:
		return float64(i)
	case uint64:
		return float64(i)
	case uint32:
		return float64(i)
	case uint:
		return float64(i)
	case string:
		output, err := strconv.ParseFloat(i, 64)
		if err != nil {
			log.Error(err)
			return 0
		}
		return output
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float()
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			output, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Error(err)
				return 0
			}
			return output
		} else {
			log.Error("can't convert %v to float64", v.Type())
			return 0
		}
	}
}
