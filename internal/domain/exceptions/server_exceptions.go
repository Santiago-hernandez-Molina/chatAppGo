package exceptions

import "fmt"

type AccesDataException struct{}

func (*AccesDataException) Error() string {
	return fmt.Sprint("Error retreiving data")
}

type NilValueException struct{}

func (*NilValueException) Error() string {
	return fmt.Sprint("Error reading values")
}
