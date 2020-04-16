package dispatch

import (
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

type taskGraph struct {
	graph map[string]*task
	todo  mapset.Set
}

type task struct {
	out      map[string]bool
	in       map[string]bool
	outCount int
	inCount  int
	done     bool
}

func NewTaskGraph() *taskGraph {
	graph := &taskGraph{
		graph: make(map[string]*task),
		todo:  mapset.NewSet(),
	}
	return graph
}

func NewTask(depens []string) *task {
	task := &task{
		out:      make(map[string]bool),
		in:       make(map[string]bool),
		outCount: 0,
		inCount:  0,
		done:     false,
	}
	for _, dep := range depens {
		task.out[dep] = false
		task.outCount++
	}
	return task
}

func (t *taskGraph) AddTask(name string, depens []string) bool {
	if _, ok := t.graph[name]; ok {
		return false
	}
	task := NewTask(depens)
	t.graph[name] = task
	return true
}

func (t *taskGraph) InitGraph()  {
	for name, task := range t.graph {
		for item := range task.out {
			t.graph[item].in[name] = false
			t.graph[item].inCount++
		}
		if task.outCount == 0 {
			t.todo.Add(name)
		}
	}
}

func (t *taskGraph) MarkTaskDone(taskName string) bool {
	if !t.todo.Contains(taskName) {
		return false
	}
	t.todo.Remove(taskName)
	task := t.graph[taskName]
	task.done = true
	for name, _ := range task.in {
		ta := t.graph[name]
		ta.out[taskName] = true
		ta.outCount--
		if ta.outCount == 0 {
			t.todo.Add(name)
		}
	}
	for name, _ := range task.out {
		ta := t.graph[name]
		ta.in[taskName] = true
		ta.inCount--
	}
	return true
}

func (t *taskGraph) AddDoneTasks(tasks []string) {
	dep := 0
	t.doAddDoneTasks(tasks, &dep)
}

func (t *taskGraph) doAddDoneTasks(tasks []string, dep *int)  {
	if len(tasks) == 0 {
		return
	}
	*dep++
	residue := make([]string, 0)
	for _, taskName := range tasks {
		if !t.MarkTaskDone(taskName) {
			residue = append(residue, taskName)
		}
	}
	if *dep > 100 {
		return
	}
	t.doAddDoneTasks(residue, dep)
}

func (t *taskGraph) GetTodoTasks() []string {
	var todo []string
	for taskName := range t.todo.Iter() {
		todo = append(todo, taskName.(string))
	}
	return todo
}

func (t *taskGraph) PrintGraph() {
	fmt.Println("-----------------------------------")
	for k, node := range t.graph {
		fmt.Println("任务名：", k)
		if node.done {
			fmt.Println("是否完成：", "YES")
		} else {
			fmt.Println("是否完成：", "NO")
		}
		fmt.Println("（当前）依赖这些任务：")
		for taskName, v := range node.out {
			if !v {
				fmt.Print(" ", taskName, " ")
			}
		}
		fmt.Println()
		fmt.Println("（当前）被这些任务依赖：")
		for taskName, v := range node.in {
			if !v {
				fmt.Print(" ", taskName, " ")
			}
		}
		fmt.Println()
	}
	fmt.Println("-----------------------------------")
}