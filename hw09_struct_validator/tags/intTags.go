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
func (t *IntTags) String() string {
	return fmt.Sprintf("{min:%v max:%v in:%v}", *t.min, *t.max, t.in)
}

func (t *IntTags) FillField(tag string) error {
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
		t.min = &i
	case "max":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			errorLog.Printf("strconv.Atoi error %s", err)
			return ErrTagInvalidSyntax
		}
		t.max = &i
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
		t.in = arr
	default:
		return ErrUnsupportedTag
	}
	return nil
}

func (t *IntTags) ValidateValue(i reflect.Value) error {
	infoLog.Printf("\tvalidate value '%v' with tags %+v\n", i, *t)
	val := int(i.Int())

	// min
	if (t.min != nil) && (val < *t.min) {
		return ErrInvaildByTag
	}

	// max
	if (t.max != nil) && (val > *t.max) {
		return ErrInvaildByTag
	}

	// in
	if (len(t.in) != 0) && ((val < t.in[0]) || (val > t.in[1])) {
		return ErrInvaildByTag
	}

	return nil
}
