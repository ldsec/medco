package survivalclient

import (
	"fmt"
	"strings"
)

func parseStartsEndsWhen(whenArgument string) (when string, err error) {
	when = strings.TrimSpace(strings.ToLower(whenArgument))

	switch when {
	case "first":
		when = "earliest"
	case "earliest":
		break

	case "last":
		when = "last"
	case "latest":
		break

	default:
		err = fmt.Errorf(`argument "%s" not recognized, it must be one of "first", "last", "earliest" and "latest" (case non-sensitive)`, whenArgument)
		break
	}
	return
}
