package http

type HttpErr interface {
	error
	Status() int
}

type httpErr struct {
	status int
	msg    string
}

func (h httpErr) Error() string {
	return h.msg
}

func (h httpErr) Status() int {
	return h.status
}

var (
	ErrUnknown        HttpErr = httpErr{500, "Unknown error"}
	ErrBadSyntax      HttpErr = httpErr{400, "Body syntax error"}
	ErrBadBody        HttpErr = httpErr{422, "Bad body content"}
	ErrNoNext         HttpErr = httpErr{422, "Finale has no next section"}
	ErrSectionUnknown HttpErr = httpErr{404, "Unknown section"}
	ErrBadSecret      HttpErr = httpErr{401, "Not Authorized"}
)
