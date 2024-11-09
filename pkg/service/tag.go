package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(taskId int, tag taskFlow.Tag) (int, error) {
	return s.repo.CreateTag(taskId, tag)
}
func (s *TagService) GetTags(taskId int) ([]taskFlow.Tag, error) {
	return s.repo.GetTags(taskId)
}
func (s *TagService) AddTag(taskId int, tag taskFlow.Tag) (int, error) {
	return s.repo.AddTag(taskId, tag)
}
func (s *TagService) DeleteTag(tagId int) error {
	return s.repo.DeleteTag(tagId)
}
func (s *TagService) ChangeTag(taskId, newTag int) error {
	return s.repo.ChangeTag(taskId, newTag)
}
func (s *TagService) UpdateTag(tagId int, input taskFlow.TagInput) error {
	return s.repo.UpdateTag(tagId, input)
}

func (s *TagService) GetAllTags() ([]taskFlow.Tag, error) {
	return s.repo.GetAllTags()
}
