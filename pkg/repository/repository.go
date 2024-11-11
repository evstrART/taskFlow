package repository

import (
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
	GetUser(username, password string) (taskFlow.User, error)
}

type Task interface {
	CreateTask(projectId int, input taskFlow.Task) (int, error)
	GetAllTasks(projectId int) ([]taskFlow.Task, error)
	GetTask(projectId int, id int) (taskFlow.Task, error)
	DeleteTask(projectId int, taskId int) error
	UpdateTask(projectId int, taskId int, input taskFlow.UpdateTaskInput) error
	GetAllTasksForUser(userID int) ([]taskFlow.Task, error)
}

type User interface {
	GetUser(userId int) (taskFlow.User, error)
	GetUsers() ([]taskFlow.User, error)
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

type Project interface {
	Create(ownerID int, project taskFlow.Project) (int, error)
	GetAllProjects() ([]taskFlow.Project, error)
	GetProjectById(id int) (taskFlow.Project, error)
	DeleteProject(id int) error
	UpdateProject(id int, input taskFlow.UpdateProjectInput) error
	AddMembers(projectId, userId int, input taskFlow.AddMemberRequest) error
}

type Repository struct {
	AutorisationService
	Project
	Task
	User
	Comment
	Tag
	Log
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AutorisationService: NewAuthPostgres(db),
		Project:             NewProjectPostgres(db),
		Task:                NewTaskPostgres(db),
		Comment:             NewCommentPostgres(db),
		Tag:                 NewTagPostgres(db),
		User:                NewUserPostgres(db),
	}
}
