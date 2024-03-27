package helper

type HttpErr struct {
	Code int
	Err  error
}

func (h *HttpErr) Error() string {
	return h.Err.Error()
}

func (ctx *FiberCtx) JsonResponse(respCode int, data map[string]interface{}) error {
	return ctx.Fctx.
		Status(respCode).
		JSON(data)
}
