package geenee

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

func (i *IntSlice) String() string {
	return strings.Trim(strings.Replace(fmt.Sprint(*i), " ", ", ", -1), "[]")
}

func (i *IntSlice) Set(value string) error {
	temp := []int{}
	if strings.Contains(value, ",") {
		intValues := strings.Split(value, ",")
		for _, v := range intValues {
			intValue, err := strconv.Atoi(strings.Trim(v, " "))
			if err != nil {
				return err
			}
			temp = append(temp, intValue)
		}

		*i = append(*i, temp...)

		return nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*i = append(*i, intValue)

	return nil
}

type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Trim(strings.Replace(fmt.Sprint(*s), " ", ", ", -1), "[]")
}

func (s *StringSlice) Set(value string) error {
	sr := strings.NewReader(value)
	cr := csv.NewReader(sr)
	values, err := cr.Read()
	if err != nil {
		return err
	}

	for _, v := range values {
		*s = append(*s, strings.Trim(v, " "))
	}

	return nil
}
