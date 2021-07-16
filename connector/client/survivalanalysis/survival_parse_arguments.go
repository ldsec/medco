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

func parseCensoringFrom(censoringArgument string) (censoring string, err error) {
	censoring = strings.TrimSpace(strings.ToLower(censoringArgument))

	switch censoring {
	case "observation":
		censoring = "observations"
	case "observations":
		break

	case "visits":
	case "visit":
	case "encounter":
		censoring = "encounters"
	case "encounters":
		break
	default:
		err = fmt.Errorf(`argument "%s" not recognized, it must be one of "observations", "observation","encounter", "encounters", "visit" and "visits" (case non-sensitive)`, censoringArgument)
	}
	return
}
