package exceptions

import "fmt"

type (
	DuplicatedUser struct{}
	UserNotFound   struct{}
)

func (*DuplicatedUser) Error() string {
	return fmt.Sprint("Duplicated user")
}

func (*UserNotFound) Error() string {
	return fmt.Sprint("User not found")
}
