package di

type stringErr string

func (e stringErr) Error() string {
	return string(e)
}

const (
	ErrWrongConfigType      = stringErr("wrong config type, should be struct")
	ErrWrongBuilderType     = stringErr("wrong builder type, should be func")
	ErrWrongBuilderReturn   = stringErr("wrong builder should return (<service>, error)")
	ErrParamBuilderNotFound = stringErr("param builder not found")
	ErrNotFound             = stringErr("not found")
	ErrWrongValue           = stringErr("value should be a pointer")
	ErrWrongCallbackReturn  = stringErr("can't assign callback return to value")
)
