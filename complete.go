package jupyter

type Completion struct {
	class,
	name,
	typ string
}

type CompletionResponse struct {
	partial     int
	completions []Completion
}

func handleCompleteRequest(ir Interpreter, receipt msgReceipt) error {
	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	cursorPos := int(reqcontent["cursor_pos"].(float64))

	// autocomplete the code at the cursor position
	_, matches, _ := ir.CompleteWords(code, cursorPos)

	// prepare the reply
	content := make(map[string]interface{})

	content["ename"] = "ERROR"
	content["evalue"] = "no completions found"
	content["traceback"] = nil
	content["status"] = "error"

	if len(matches) == 0 {
		content["ename"] = "ERROR"
		content["evalue"] = "no completions found"
		content["traceback"] = nil
		content["status"] = "error"
	}
	//else {
	//	partialWord := interp.TailIdentifier(prefix)
	//	content["cursor_start"] = float64(len(prefix) - len(partialWord))
	//	content["cursor_end"] = float64(cursorPos)
	//	content["matches"] = matches
	//	content["status"] = "ok"
	//}

	return receipt.Reply("complete_reply", content)
}
