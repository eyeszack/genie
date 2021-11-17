package geenee

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

func (i *IntSlice) String() string {
	return strings.Trim(strings.Replace(fmt.Sprint(*i), " ", ", ", -1), "[]")
}

func (i *IntSlice) Set(value string) error {
	if strings.Contains(value, ",") {
		intValues := strings.Split(value, ",")
		for _, v := range intValues {
			intValue, err := strconv.Atoi(strings.Trim(v, " "))
			if err != nil {
				return err
			}
			*i = append(*i, intValue)
		}

		return nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*i = append(*i, intValue)

	return nil
}
