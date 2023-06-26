package tesla

import (
	"fmt"
	"strconv"
	"strings"
)

type position struct {
	lat  float64
	long float64
}

func parsePosition(pos string) (position, error) {

	parts := strings.Split(pos, " ")
	if len(parts) != 2 {
		return position{}, fmt.Errorf("Invalid position '%s' (split fails)", pos)
	}

	latitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return position{}, fmt.Errorf("Invalid position '%s' (parse lat) : %w", pos, err)
	}

	longitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return position{}, fmt.Errorf("Invalid position '%s' (parse long) : %w", pos, err)
	}

	return position{lat: latitude, long: longitude}, nil
}
