package main

import (
	"errors"
	"fmt"
	"github.com/cosmos72/gomacro/ast2"
	"github.com/cosmos72/gomacro/base"
	basereflect "github.com/cosmos72/gomacro/base/reflect"
	interp "github.com/cosmos72/gomacro/fast"
	"github.com/cosmos72/gomacro/xreflect"
)

type Interpreter interface {
	CompleteWords(code string, cursorPos int) (prefix string, completions []string, tail string)
}

// doEval evaluates the code in the interpreter. This function captures an uncaught panic
// as well as the values of the last statement/expression.
func doEval(ir *interp.Interp, outerr OutErr, code string) (val []interface{}, typ []xreflect.Type, err error) {

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

	// Prepare and perform the multiline evaluation.
	compiler := ir.Comp

	// Don't show the gomacro prompt.
	compiler.Options &^= base.OptShowPrompt

	// Don't swallow panics as they are recovered above and handled with a Jupyter `error` message instead.
	compiler.Options &^= base.OptTrapPanic

	// Reset the error line so that error messages correspond to the lines from the cell.
	compiler.Line = 0

	// Parse the input code (and don't perform gomacro's macroexpansion).
	// These may panic but this will be recovered by the deferred recover() above so that the error
	// may be returned instead.
	nodes := compiler.ParseBytes([]byte(code))
	srcAst := ast2.AnyToAst(nodes, "doEval")

	// If there is no srcAst then we must be evaluating nothing. The result must be nil then.
	if srcAst == nil {
		return nil, nil, nil
	}

	// Compile the ast.
	compiledSrc := ir.CompileAst(srcAst)

	// Evaluate the code.
	results, types := ir.RunExpr(compiledSrc)

	// Convert results from xreflect.Value to interface{}
	values := make([]interface{}, len(results))
	for i, result := range results {
		values[i] = basereflect.ValueInterface(result)
	}

	return values, types, nil
}
