package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type ProjectService struct {
	repo repository.Project
}

func NewProjectService(repo repository.Project) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ownerID int, project taskFlow.Project) (int, error) {
	return s.repo.Create(ownerID, project)
}

func (s *ProjectService) GetAllProjects() ([]taskFlow.Project, error) {
	return s.repo.GetAllProjects()
}

func (s *ProjectService) GetProjectById(id int) (taskFlow.Project, error) {
	return s.repo.GetProjectById(id)
}

func (s *ProjectService) DeleteProject(userID, id int) error {
	return s.repo.DeleteProject(userID, id)
}

func (s *ProjectService) UpdateProject(userId, id int, input taskFlow.UpdateProjectInput) error {
	return s.repo.UpdateProject(userId, id, input)
}
func (s *ProjectService) AddMembers(projectId int, input taskFlow.AddMemberRequest) error {
	return s.repo.AddMembers(projectId, input)
}
func (s *ProjectService) GetMembers(projectId int) ([]taskFlow.User, error) {
	return s.repo.GetMembers(projectId)
}

func (s *ProjectService) DeleteMember(projectId, memberId int) error {
	return s.repo.DeleteMember(projectId, memberId)
}
func (s *ProjectService) CompleteProject(projectId, userID int) error {
	return s.repo.CompleteProject(projectId, userID)
}
