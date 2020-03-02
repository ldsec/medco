package survivalserver

import (
	"errors"
	"fmt"
	"runtime"
)

const trace = 10

func NiceError(err error) error {
	if err == nil {
		return nil
	} else {
		str := err.Error()
		programCounter := make([]uintptr, trace)
		callersFound := runtime.Callers(0, programCounter)
		for i := 0; i < callersFound; i++ {
			caller := runtime.FuncForPC(programCounter[i])
			file, line := caller.FileLine(programCounter[i])
			str += fmt.Sprintf("in function %s, in file %s, at line %d\n", caller.Name(), file, line)
		}
		return errors.New(str)

	}

}
