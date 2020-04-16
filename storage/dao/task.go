package dao

import (
	"github.com/jinzhu/gorm"

	"dispatch/model/bean"
)

type TaskDaoImpl struct {

}

func (d *TaskDaoImpl) Create(db *gorm.DB, task *bean.Task) error {
	if err := db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (d *TaskDaoImpl) Query(db *gorm.DB, id int64) (*bean.Task, error) {
	task := &bean.Task{}
	if err := db.Where("id = ?", id).First(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (d *TaskDaoImpl) QueryByState(db *gorm.DB, state string) ([]*bean.Task, error) {
	tasks := make([]*bean.Task, 0)
	if err := db.Where("task_state = ?", state).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
