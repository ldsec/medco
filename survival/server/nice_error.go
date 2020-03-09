package survivalserver

import (
	"errors"
	"fmt"
	"runtime"
)

// NiceError returns nil if err is nil, otherwise return a new error with original error string together with a stack trace if available
func NiceError(err error) error {
	if err == nil {
		return nil
	}
	str := err.Error()
	programCounter := make([]uintptr, errLogTrace)
	callersFound := runtime.Callers(0, programCounter)
	for i := 0; i < callersFound; i++ {
		caller := runtime.FuncForPC(programCounter[i])
		file, line := caller.FileLine(programCounter[i])
		str += fmt.Sprintf("in function %s, in file %s, at line %d\n", caller.Name(), file, line)
	}
	return errors.New(str)

}
