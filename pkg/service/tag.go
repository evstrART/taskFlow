package service

import "github.com/evstrART/taskFlow/pkg/repository"

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}
