package repository

import (
	"github.com/evstrART/taskFlow"
	"github.com/jmoiron/sqlx"
)

type AutorisationService interface {
	CreateUser(user taskFlow.User) (int, error)
	GetUser(username, password string) (taskFlow.User, error)
	ChangePassword(userId int, newPasswordHex string) error
	UserExistsForReset(input taskFlow.ResetPasswordInput) (bool, error)
	GetUserByNameAndEmail(username, email string) (taskFlow.User, error)
}

type Task interface {
	CreateTask(userID, projectId int, input taskFlow.Task) (int, error)
	GetAllTasks(projectId int) ([]taskFlow.Task, error)
	GetTask(projectId int, id int) (taskFlow.Task, error)
	DeleteTask(userID, projectId int, taskId int) error
	UpdateTask(userID, projectId int, taskId int, input taskFlow.UpdateTaskInput) error
	GetAllTasksForUser(userID int) ([]taskFlow.Task, error)
	CompleteTask(taskId, userId int) error
}

type User interface {
	GetUser(userId int) (taskFlow.User, error)
	GetUsers() ([]taskFlow.User, error)
	UpdateUser(userId int, input taskFlow.UpdateUserInput) error
	DeleteUser(userId int) error
	GetUserByNameAndEmail(username, email string) (taskFlow.User, error)
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
type Admin interface {
	SelectAdminId(userId int) ([]int, error)
	SelectActivityLogs() ([]taskFlow.ActivityLog, error)
	SelectTable(tableName string) ([]map[string]interface{}, error)
	InsertUser(user taskFlow.User) error
	InsertProject(project taskFlow.Project) error
	InsertTask(task taskFlow.Task) error
	InsertComment(comment taskFlow.Comment) error
	InsertProjectMember(member taskFlow.ProjectMember) error
	InsertActivityLog(log taskFlow.ActivityLog) error
	InsertTag(tag taskFlow.Tag) error
	InsertTaskTag(taskTag taskFlow.TaskTag) error
	GetCompletedTasksByProject() ([]taskFlow.ProjectStats, error)
	GetCreatedTasksByUser() ([]taskFlow.UserStats, error)
}

type Project interface {
	Create(ownerID int, project taskFlow.Project) (int, error)
	GetAllProjects() ([]taskFlow.Project, error)
	GetProjectById(id int) (taskFlow.Project, error)
	DeleteProject(userID, id int) error
	UpdateProject(userId, id int, input taskFlow.UpdateProjectInput) error
	AddMembers(projectId int, input taskFlow.AddMemberRequest) error
	GetMembers(projectId int) ([]taskFlow.User, error)
	DeleteMember(projectId, memberId int) error
	CompleteProject(projectId, userID int) error
	GetAllProjectsForUser(userID int) ([]taskFlow.Project, error)
}

type Repository struct {
	AutorisationService
	Project
	Task
	User
	Comment
	Tag
	Admin
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AutorisationService: NewAuthPostgres(db),
		Project:             NewProjectPostgres(db),
		Task:                NewTaskPostgres(db),
		Comment:             NewCommentPostgres(db),
		Tag:                 NewTagPostgres(db),
		User:                NewUserPostgres(db),
		Admin:               NewAdminPostgres(db),
	}
}
