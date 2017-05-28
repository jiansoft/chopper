package core

import "reflect"

type IExecutor interface {
	ExecuteTasks(t []Task)
	ExecuteTask(t Task)
}

type Task struct {
	PaddingFunc interface{}
	FuncParams  []interface{}
}

type defaultExecutor struct {
}

func NewDefaultExecutor() defaultExecutor { return defaultExecutor{} }

func (d defaultExecutor) ExecuteTasks(tasks []Task) {
	for _,task := range tasks {
		task.Execute()
	}
}

func (t Task) Execute() {
	execFunc := reflect.ValueOf(t.PaddingFunc)
	params := make([]reflect.Value, len(t.FuncParams))
	for k, param := range t.FuncParams {
		params[k] = reflect.ValueOf(param)
	}
	func(in []reflect.Value) { _ = execFunc.Call(in) }(params)
}
