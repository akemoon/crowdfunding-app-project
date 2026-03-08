package validation

// TODO: move to golib

type Error struct {
	fields map[string]string
}

func (e *Error) Error() string {
	return "validation failed"
}

func (e *Error) Add(field, msg string) {
	if e.fields == nil {
		e.fields = make(map[string]string)
	}
	e.fields[field] = msg
}

func (e *Error) Fields() map[string]string {
	return e.fields
}

func (e *Error) HasErrors() bool {
	return len(e.fields) > 0
}