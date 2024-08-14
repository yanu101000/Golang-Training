package service

import (
	"fmt"
	"time"
	"wallet/entity"
)

type RecordService interface {
	CreateRecord(record *entity.Record) (*entity.Record, error)
	GetRecordByID(id int64) (*entity.Record, error)
	UpdateRecord(record *entity.Record) (*entity.Record, error)
	DeleteRecord(id int64) error
	GetRecordsByTimeRange(walletID int64, startTime, endTime string) ([]*entity.Record, error)
	GetLast10Records() ([]*entity.Record, error)
	GetRecords() map[int64]*entity.Record
}

type recordService struct {
	records map[int64]*entity.Record
	nextID  int64
}

func NewRecordService() RecordService {
	return &recordService{
		records: make(map[int64]*entity.Record),
		nextID:  1,
	}
}

func (s *recordService) CreateRecord(record *entity.Record) (*entity.Record, error) {
	record.ID = s.nextID
	s.records[s.nextID] = record
	s.nextID++
	return record, nil
}

func (s *recordService) GetRecordByID(id int64) (*entity.Record, error) {
	record, exists := s.records[id]
	if !exists {
		return nil, fmt.Errorf("record not found")
	}
	return record, nil
}

func (s *recordService) UpdateRecord(record *entity.Record) (*entity.Record, error) {
	_, exists := s.records[record.ID]
	if !exists {
		return nil, fmt.Errorf("record not found")
	}
	s.records[record.ID] = record
	return record, nil
}

func (s *recordService) DeleteRecord(id int64) error {
	_, exists := s.records[id]
	if !exists {
		return fmt.Errorf("record not found")
	}
	delete(s.records, id)
	return nil
}

func (s *recordService) GetRecordsByTimeRange(walletID int64, startTime, endTime string) ([]*entity.Record, error) {
	var result []*entity.Record
	start, _ := time.Parse(time.RFC3339, startTime)
	end, _ := time.Parse(time.RFC3339, endTime)

	for _, record := range s.records {
		if record.WalletID == walletID {
			recordTime, _ := time.Parse(time.RFC3339, record.Timestamp)
			if recordTime.After(start) && recordTime.Before(end) {
				result = append(result, record)
			}
		}
	}

	return result, nil
}

func (s *recordService) GetLast10Records() ([]*entity.Record, error) {
	var result []*entity.Record
	count := 0

	for _, record := range s.records {
		if count < 10 {
			result = append(result, record)
			count++
		} else {
			break
		}
	}

	return result, nil
}

func (s *recordService) GetRecords() map[int64]*entity.Record {
	return s.records
}
