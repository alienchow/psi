package region

import "strings"

type Region string

const (
	North   Region = `North`
	South          = `South`
	East           = `East`
	West           = `West`
	Central        = `Central`
	Overall        = `Overall`
	Invalid        = `Invalid`
)

var All = []Region{
	North,
	South,
	East,
	West,
	Central,
	Overall,
	Invalid,
}

func GetByArg(arg string) Region {
	switch strings.ToUpper(arg) {
	case "N", "NORTH":
		return North
	case "S", "SOUTH":
		return South
	case "E", "EAST":
		return East
	case "W", "WEST":
		return West
	case "C", "CENTRAL":
		return Central
	case "O", "OVERALL":
		return Overall
	}
	return Invalid
}
