package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type TaskService struct {
	repo        repository.Task
	projectRepo repository.Project
}

func NewTaskService(repo repository.Task, projectRepo repository.Project) *TaskService {
	return &TaskService{repo: repo, projectRepo: projectRepo}
}

func (s *TaskService) CreateTask(projectId int, input taskFlow.Task) (int, error) {
	_, err := s.projectRepo.GetProjectById(projectId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateTask(projectId, input)
}

func (s *TaskService) GetAllTasks(projectId int) ([]taskFlow.Task, error) {
	_, err := s.projectRepo.GetProjectById(projectId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllTasks(projectId)
}

func (s *TaskService) GetTask(projectId int, taskId int) (taskFlow.Task, error) {
	_, err := s.projectRepo.GetProjectById(projectId)
	if err != nil {
		return taskFlow.Task{}, err
	}
	return s.repo.GetTask(projectId, taskId)
}

func (s *TaskService) DeleteTask(projectId int, taskId int) error {
	return s.repo.DeleteTask(projectId, taskId)
}

func (s *TaskService) UpdateTask(projectId int, taskId int, input taskFlow.UpdateTaskInput) error {
	_, err := s.projectRepo.GetProjectById(projectId)
	if err != nil {
		return err
	}
	return s.repo.UpdateTask(projectId, taskId, input)
}
func (s *TaskService) GetAllTasksForUser(userID int) ([]taskFlow.Task, error) {
	return s.repo.GetAllTasksForUser(userID)
}
