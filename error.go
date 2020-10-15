package rj

type RJError struct {
	msg string
}

func newError(msg string) *RJError {
	return &RJError{msg: msg}
}

func (e *RJError) Error() string {
	return e.msg
}

func (e *RJError) addMsg(msg string) {
	if msg == "" {
		return
	}

	if e.msg != "" {
		e.msg += ";"
	}
	e.msg += msg
}
