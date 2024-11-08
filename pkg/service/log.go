package service

import "github.com/evstrART/taskFlow/pkg/repository"

type LogService struct {
	repo repository.Log
}

func NewLogService(repo repository.Log) *LogService {
	return &LogService{repo: repo}
}
