package exceptions

import "fmt"

type DuplicatedUser struct{}

func (*DuplicatedUser) Error() string {
	return fmt.Sprint("Duplicated user")
}

type UserNotFound struct{}

func (*UserNotFound) Error() string {
	return fmt.Sprint("User not found")
}
