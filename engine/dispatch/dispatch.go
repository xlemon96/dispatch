package dispatch

import (
	"log"
	"strconv"
	"sync"

	"dispatch/constant"
	"dispatch/model/running"
	"dispatch/storage"
)

type dispatch struct {
	sync.Mutex
	logger     *log.Logger
	taskDao    storage.Storage
	taskGraphs map[int64]*taskGraph //key为taskID
	taskDags   map[int64]*dag       //key为taskID
	todoDags   chan *dagBag
}

type dag struct {
	task *running.Task
	dags map[int64]*running.DAGInstance //key为dagID
}

type dagBag struct {
	task         *running.Task
	dagInstances []*running.DAGInstance
}

func NewDispatch(taskDao storage.Storage, logger *log.Logger) *dispatch {
	return &dispatch{
		logger:     logger,
		taskDao:    taskDao,
		taskGraphs: make(map[int64]*taskGraph),
		taskDags:   make(map[int64]*dag),
		//todo,缓冲区大小
		todoDags: make(chan *dagBag, 100),
	}
}

func NewDagBag(task *running.Task, dagInstances []*running.DAGInstance) *dagBag {
	return &dagBag{
		task:         task,
		dagInstances: dagInstances,
	}
}

func (d *dispatch) Start() error {
	if err := d.init(); err != nil {
		return err
	}
	return nil
}

func (d *dispatch) GetTodoDags() chan *dagBag {
	return d.todoDags
}

func (d *dispatch) init() error {
	d.Lock()
	defer d.Unlock()
	tasks, err := d.taskDao.DescribeTasks(constant.TaskStateRunning)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return nil
	}
	for _, task := range tasks {
		err := d.initTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dispatch) initTask(task *running.Task) error {
	if _, ok := d.taskDags[task.Id]; ok {
		return nil
	}
	dags, err := d.taskDao.DescribeDAGInstances(task.Id)
	if err != nil {
		return err
	}
	var graph *taskGraph
	if _, ok := d.taskGraphs[task.Id]; !ok {
		graph = NewTaskGraph()
		for _, dag := range dags {
			graph.AddTask(int64ToString(dag.Id), dag.Depends)
		}
		doneDag := make([]string, 0)
		for _, dag := range dags {
			if dag.State == constant.DAGStateFailed {
				return nil
			}
			if dag.State == constant.DAGStateFailed || dag.State == constant.DAGStateSucceed {
				doneDag = append(doneDag, int64ToString(dag.Id))
			}
		}
		graph.InitGraph()
		graph.AddDoneTasks(doneDag)
		//graph.PrintGraph()
		//todo,判断task是否已经完成
	}
	d.taskGraphs[task.Id] = graph
	taskDag := &dag{
		task: task,
		dags: make(map[int64]*running.DAGInstance),
	}
	for _, dag := range dags {
		taskDag.dags[dag.Id] = dag
	}
	d.taskDags[task.Id] = taskDag
	d.sendTask(graph, task.Id)
	return nil
}

func (d *dispatch) sendTask(graph *taskGraph, taskID int64) {
	taskDag := d.taskDags[taskID]
	dagNames := graph.GetTodoTasks()
	dags := make([]*running.DAGInstance, 0, len(dagNames))
	for _, dagName := range dagNames {
		if _, ok := taskDag.dags[stringToInt64(dagName)]; ok {
			dags = append(dags, taskDag.dags[stringToInt64(dagName)])
		}
	}
	if len(dags) == 0 {
		return
	}
	dagBag := &dagBag{
		task:         taskDag.task,
		dagInstances: dags,
	}
	d.todoDags <- dagBag
}

func (d *dispatch) reschedule() error {
	dagInstances, err := d.taskDao.DescribeDAGInstancesByState(constant.DAGStateRescheduling)
	if err != nil {
		return err
	}
	if len(dagInstances) == 0 {
		return nil
	}
	dagBags := make(map[int64]*dagBag)
	for _, dagInstance := range dagInstances {
		//？？？
		if d.taskDags[dagInstance.TaskId].dags[dagInstance.Id].State == constant.DAGStateRescheduling {
			continue
		}
		if _, ok := dagBags[dagInstance.TaskId]; ok {
			dagBags[dagInstance.TaskId].dagInstances = append(dagBags[dagInstance.TaskId].dagInstances, dagInstance)
		} else {
			task, err := d.taskDao.DescribeTask(dagInstance.TaskId)
			if err != nil {
				return err
			}
			dagBags[dagInstance.TaskId] = &dagBag{
				task:         task,
				dagInstances: []*running.DAGInstance{dagInstance},
			}
		}
	}
	for _, dagBag := range dagBags {
		taskId := dagBag.task.Id
		for _, dagInstance := range dagBag.dagInstances {
			d.taskDags[taskId].dags[dagInstance.Id].State = constant.DAGStateRescheduling
		}
		d.todoDags <- dagBag
	}
	return nil
}

func (d *dagBag) GetDagInstances() []*running.DAGInstance {
	return d.dagInstances
}

func (d *dagBag) GetTask() *running.Task {
	return d.task
}

func int64ToString(content int64) string {
	return strconv.FormatInt(content, 10)
}

func stringToInt64(content string) int64 {
	if number, err := strconv.ParseInt(content, 10, 64); err != nil {
		return -1
	} else {
		return number
	}
}
