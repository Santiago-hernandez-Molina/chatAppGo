package tasks

import (
	"time"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type UserTask struct {
	userRepo ports.UserRepo
}

func (task *UserTask) DeleteAccountTask(email string) error {
	wait := time.Minute * 10
	time.Sleep(wait)
	err := task.userRepo.DeleteUserByEmailAndStatus(email, false)
	if err != nil {
		return err
	}
	return nil
}

var _ ports.UserTask = (*UserTask)(nil)

func NewUserTasks(userRepo ports.UserRepo) *UserTask {
	return &UserTask{
		userRepo: userRepo,
	}
}
