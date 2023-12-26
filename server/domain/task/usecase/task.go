package usecase

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/db"
	"Dur4nC2/server/domain/models"
	_taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"github.com/gofrs/uuid"
)

type taskUsecase struct {
	repository models.BeaconTaskRepository
}

func NewTaskUsecase() models.BeaconTaskUsecase {
	b := _taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	return &taskUsecase{repository: b}
}

func (u *taskUsecase) List() (*clientpb.BeaconTasks, error) {
	tasks, err := u.repository.List()
	if err != nil {
		return nil, err
	}
	tasks_pb := &clientpb.BeaconTasks{Tasks: []*clientpb.BeaconTask{}}
	for _, task := range tasks {
		tasks_pb.Tasks = append(tasks_pb.Tasks, &clientpb.BeaconTask{
			ID:          task.ID.String(),
			BeaconID:    task.BeaconID.String(),
			CreatedAt:   task.CreatedAt.Unix(),
			State:       task.State,
			SentAt:      task.SentAt.Unix(),
			CompletedAt: task.CompletedAt.Unix(),
			Request:     task.Request,
			Response:    task.Response,
			Description: task.Description,
		})
	}
	return tasks_pb, nil

}
func (u *taskUsecase) ListTasksStateAndBeaconID(id string, state string) (*clientpb.BeaconTasks, error) {
	tasks, err := u.repository.ListTasksStateAndBeaconID(id, state)
	if err != nil {
		return nil, err
	}
	tasks_pb := &clientpb.BeaconTasks{Tasks: []*clientpb.BeaconTask{}}
	for _, task := range tasks {
		tasks_pb.Tasks = append(tasks_pb.Tasks, &clientpb.BeaconTask{
			ID:          task.ID.String(),
			BeaconID:    task.BeaconID.String(),
			CreatedAt:   task.CreatedAt.Unix(),
			State:       task.State,
			SentAt:      task.SentAt.Unix(),
			CompletedAt: task.CompletedAt.Unix(),
			Request:     task.Request,
			Response:    task.Response,
			Description: task.Description,
		})
	}
	return tasks_pb, nil
}
func (u *taskUsecase) Read(task_pb *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	id, _ := uuid.FromString(task_pb.ID)
	task, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}
	task_pb = &clientpb.BeaconTask{
		ID:          task.ID.String(),
		BeaconID:    task.BeaconID.String(),
		CreatedAt:   task.CreatedAt.Unix(),
		State:       task.State,
		SentAt:      task.SentAt.Unix(),
		CompletedAt: task.CompletedAt.Unix(),
		Request:     task.Request,
		Response:    task.Response,
		Description: task.Description,
	}
	return task_pb, nil
}
func (u *taskUsecase) Cancel(task_pb *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	id, _ := uuid.FromString(task_pb.ID)
	err := u.repository.Cancel(id)
	if err != nil {
		return nil, err
	}
	canceledTask, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}

	return &clientpb.BeaconTask{ID: canceledTask.BeaconID.String(), State: canceledTask.State}, nil

}
