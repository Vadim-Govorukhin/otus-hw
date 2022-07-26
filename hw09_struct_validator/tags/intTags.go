package tags

import (
	"strconv"
	"strings"
)

type IntTags struct {
	min int
	max int
	in  []int
}

func (T *IntTags) FillField(tag string) error {
	m := strings.Split(tag, ":")
	if len(m) > 2 {
		return ErrUnsupportedTag
	}

	switch m[0] {
	case "min":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			ErrorLog.Printf("parsing error %e", err)
			return err
		}
		T.min = i
	case "max":
		i, err := strconv.Atoi(m[1])
		if err != nil {
			ErrorLog.Printf("parsing error %e", err)
			return err
		}
		T.max = i
	case "in":
		var arr []int
		var err error
		for i, s := range strings.Split(m[1], ",") {
			arr[i], err = strconv.Atoi(s)
			if err != nil {
				ErrorLog.Printf("parsing error %e", err)
				return err
			}
		}
		T.in = arr
	default:
		return ErrUnsupportedTag
	}
	return nil
}

func (T *IntTags) IsValid(i interface{}) (bool, error) {
	_, ok := i.(int)
	if !ok {
		return false, ErrUnsupportedTypeField
	}

	// min

	return true, nil
}
