package dao

import (
	"github.com/jinzhu/gorm"

	"dispatch/model/bean"
)

type WorkerDaoImpl struct {

}

func (d *WorkerDaoImpl) Create(db *gorm.DB, worker *bean.Worker) error {
	if err := db.Create(worker).Error; err != nil {
		return err
	}
	return nil
}

func (d *WorkerDaoImpl) Query(db *gorm.DB, id int64) (*bean.Worker, error) {
	worker := &bean.Worker{}
	if err := db.Where("id = ?", id).First(worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (d *WorkerDaoImpl) List(db *gorm.DB, param *bean.Worker) ([]*bean.Worker, error) {
	workers := make([]*bean.Worker, 0)
	if err := db.Model(&bean.Worker{}).Where(param).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}