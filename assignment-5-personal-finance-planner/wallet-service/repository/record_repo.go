package repository

import (
	"time"
	"wallet/entity"

	"gorm.io/gorm"
)

type RecordRepository interface {
	CreateRecord(record *entity.Record) (*entity.Record, error)
	GetRecordByID(id int64) (*entity.Record, error)
	UpdateRecord(record *entity.Record) (*entity.Record, error)
	DeleteRecord(id int64) error
	GetRecordsByTimeRange(walletID int64, startTime, endTime time.Time) ([]*entity.Record, error)
	GetLast10Records() ([]*entity.Record, error)
}

type recordRepository struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) RecordRepository {
	return &recordRepository{db}
}

func (r *recordRepository) CreateRecord(record *entity.Record) (*entity.Record, error) {
	if err := r.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (r *recordRepository) GetRecordByID(id int64) (*entity.Record, error) {
	var record entity.Record
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *recordRepository) UpdateRecord(record *entity.Record) (*entity.Record, error) {
	if err := r.db.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (r *recordRepository) DeleteRecord(id int64) error {
	if err := r.db.Delete(&entity.Record{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *recordRepository) GetRecordsByTimeRange(walletID int64, startTime, endTime time.Time) ([]*entity.Record, error) {
	var records []*entity.Record
	if err := r.db.Where("wallet_id = ? AND timestamp BETWEEN ? AND ?", walletID, startTime, endTime).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *recordRepository) GetLast10Records() ([]*entity.Record, error) {
	var records []*entity.Record
	if err := r.db.Order("timestamp desc").Limit(10).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
