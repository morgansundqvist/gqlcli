package executioncontext

type ExecutionContext struct {
	Vars map[string]interface{}
}

func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Vars: make(map[string]interface{}),
	}
}

func (ctx *ExecutionContext) Set(key string, value interface{}) {
	ctx.Vars[key] = value
}

func (ctx *ExecutionContext) Get(key string) interface{} {
	return ctx.Vars[key]
}
