package jupyter

import (
	"errors"
	"fmt"
)

type Interpreter interface {
	CompleteWords(code string, cursorPos int) (prefix string, completions []string, tail string)
	Eval(code string) (values []any, err error)
}

type ReturnValue any

// doEval evaluates the code in the interpreter. This function captures an uncaught panic
// as well as the values of the last statement/expression.
func doEval(ir Interpreter, outerr OutErr, code string) (val []any, err error) {

	// Capture a panic from the evaluation if one occurs and store it in the `err` return parameter.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()

	code = evalSpecialCommands(outerr, code)

	// Evaluate the code.
	results, err := ir.Eval(code)
	//if results != nil {
	//	for _, result := range results {
	//		fmt.Println(result)
	//	}
	//}
	return results, err
}
