package dao

import (
	"github.com/jinzhu/gorm"

	"github.com/navieboy/dispatch/model/bean"
)

type DagInstanceDaoImpl struct {

}

func (d *DagInstanceDaoImpl) Create(db *gorm.DB, dagInstance *bean.DagInstance) error {
	if err := db.Create(dagInstance).Error; err != nil {
		return err
	}
	return nil
}

func (d *DagInstanceDaoImpl) Query(db *gorm.DB, id int64) (*bean.DagInstance, error) {
	dagInstance := &bean.DagInstance{}
	if err := db.Where("id = ?", id).First(dagInstance).Error; err != nil {
		return nil, err
	}
	return dagInstance, nil
}

func (d *DagInstanceDaoImpl) QueryByState(db *gorm.DB, state string) ([]*bean.DagInstance, error) {
	dagInstances := make([]*bean.DagInstance, 0)
	if err := db.Where("dag_state = ?", state).Find(&dagInstances).Error; err != nil {
		return nil, err
	}
	return dagInstances, nil
}

func (d *DagInstanceDaoImpl) QueryByTaskID(db *gorm.DB, taskID int64) ([]*bean.DagInstance, error) {
	dagInstances := make([]*bean.DagInstance, 0)
	if err := db.Where("task_id = ?", taskID).Find(&dagInstances).Error; err != nil {
		return nil, err
	}
	return dagInstances, nil
}
