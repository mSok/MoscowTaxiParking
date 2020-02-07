package utils

import (
	"strconv"
)

// ParseQueries parse http get params
func ParseQueries(params map[string]string) (int, int, error) {
	var (
		limit  int
		offset int
		err    error
	)
	if params["limit"] == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(params["limit"])
		if err != nil {
			return 0, 0, err
		}
	}
	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			return 0, 0, err
		}
	}
	return limit, offset, err
}
