package service

import (
	"github.com/evstrART/taskFlow"
	"github.com/evstrART/taskFlow/pkg/repository"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	ChangePassword(userId int, input taskFlow.ChangePasswordInput) error
}

type Task interface {
	CreateTask(userID, projectId int, input taskFlow.Task) (int, error)
	GetAllTasks(projectId int) ([]taskFlow.Task, error)
	GetTask(projectId int, taskId int) (taskFlow.Task, error)
	DeleteTask(userID, projectId int, taskId int) error
	UpdateTask(userID, projectId int, taskId int, input taskFlow.UpdateTaskInput) error
	GetAllTasksForUser(userId int) ([]taskFlow.Task, error)
}

type User interface {
	GetUser(userId int) (taskFlow.User, error)
	GetUsers() ([]taskFlow.User, error)
	UpdateUser(userId int, input taskFlow.UpdateUserInput) error
	CheckOldPassword(userId int, oldPassword string) (bool, error)
}

type Project interface {
	Create(ownerID int, project taskFlow.Project) (int, error)
	GetAllProjects() ([]taskFlow.Project, error)
	GetProjectById(id int) (taskFlow.Project, error)
	DeleteProject(userID, id int) error
	UpdateProject(userId, id int, input taskFlow.UpdateProjectInput) error
	AddMembers(projectId int, input taskFlow.AddMemberRequest) error
	GetMembers(projectId int) ([]taskFlow.User, error)
}

type Comment interface {
	AddComment(taskId, userId int, input taskFlow.CommentInput) (int, error)
	GetComments(taskId int) ([]taskFlow.Comment, error)
	GetCommentById(taskId, id int) (taskFlow.Comment, error)
	GetAllCommentsForUser(userId int) ([]taskFlow.Comment, error)
	DeleteComment(commentId, userId int) error
	UpdateComment(commentId, userId int, input taskFlow.CommentInput) error
}
type Tag interface {
	CreateTag(taskId int, tag taskFlow.Tag) (int, error)
	GetTags(taskId int) ([]taskFlow.Tag, error)
	AddTag(taskId int, tag taskFlow.Tag) (int, error)
	DeleteTag(tagId int) error
	ChangeTag(taskId, newTag int) error
	UpdateTag(tagId int, input taskFlow.TagInput) error
	GetAllTags() ([]taskFlow.Tag, error)
}
type Log interface {
}
type Service struct {
	AutorisationService
	Project
	Task
	User
	Comment
	Tag
	Log
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AutorisationService: NewAuthService(repos.AutorisationService),
		Project:             NewProjectService(repos.Project),
		Task:                NewTaskService(repos.Task, repos.Project),
		Comment:             NewCommentService(repos.Comment),
		Tag:                 NewTagService(repos.Tag),
		User:                NewUserService(repos.User),
	}
}
