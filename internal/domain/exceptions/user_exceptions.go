package exceptions

import "fmt"

type DuplicatedUser struct{}

func (*DuplicatedUser) Error() string {
	return fmt.Sprint("Duplicated user")
}
