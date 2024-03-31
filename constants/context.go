package constants

type ContextKey string

const (
	CancelFuncContextKey         ContextKey = "uberfxutils:cancel"
	CancelWillBeCalledContextKey ContextKey = "uberfxutils:cancelFnWillBeCalled"
)
