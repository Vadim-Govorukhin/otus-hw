package tags

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IntTags struct {
	min *int // To distinguish between value 0 and zero value
	max *int
	in  []int
}

// Pretty print
func (T *IntTags) String() string {
	return fmt.Sprintf("{min:%v max:%v in:%v}", *T.min, *T.max, T.in)
}

func (T *IntTags) FillField(tag string) error {
	m := strings.Split(tag, ":")
	if len(m) > 2 {
		return ErrTagInvalidSyntax
	}

	switch m[0] {
	case "min":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("strconv.Atoi error %s", err)
			return ErrTagInvalidSyntax
		}
		T.min = &i
	case "max":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("strconv.Atoi error %s", err)
			return ErrTagInvalidSyntax
		}
		T.max = &i
	case "in":
		var arr []int
		for _, s := range strings.Split(m[1], ",") {
			num, err := strconv.Atoi(s)
			if err != nil {
				errorLog.Printf("strconv.Atoi error %s", err)
				return ErrTagInvalidSyntax
			}
			arr = append(arr, num)
		}
		T.in = arr
	default:
		return ErrUnsupportedTag
	}
	return nil
}

func (T *IntTags) ValidateValue(i reflect.Value) error {
	infoLog.Printf("\tvalidate value '%v' with tags %+v\n", i, *T)
	val := int(i.Int())

	// min
	if (T.min != nil) && (val < *T.min) {
		return ErrInvaildByTag
	}

	// max
	if (T.max != nil) && (val > *T.max) {
		return ErrInvaildByTag
	}

	// in
	if (len(T.in) != 0) && ((val < T.in[0]) || (val > T.in[1])) {
		return ErrInvaildByTag
	}

	return nil
}
