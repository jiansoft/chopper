package core

import "reflect"

type Task struct {
	PaddingFunc interface{}
	FuncParams  []interface{}
}

type defaultExecutor struct {
}

func NewDefaultExecutor() defaultExecutor { return defaultExecutor{} }

func (d defaultExecutor) ExecuteTasks(tasks []Task) {
	for i := range tasks {
		d.ExecuteTask(tasks[i])
	}
}

func (defaultExecutor) ExecuteTask(task Task) {
	task.Execute()
}

func (t Task) Execute() {
	f := reflect.ValueOf(t.PaddingFunc)
	in := make([]reflect.Value, len(t.FuncParams))
	for k, param := range t.FuncParams {
		in[k] = reflect.ValueOf(param)
	}
	func(in []reflect.Value) { _ = f.Call(in) }(in)
}

func (t Task) ExecuteWithParams(params ...interface{}) {
	t.FuncParams = params
	t.Execute()
}
