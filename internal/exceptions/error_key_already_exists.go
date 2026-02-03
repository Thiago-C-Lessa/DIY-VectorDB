package exceptions

import "fmt"

type ErrorKeyAlreadyExists struct {
	Key string
}

func (e *ErrorKeyAlreadyExists) Error() string {
	return fmt.Sprintf("The key : %s already exists", e.Key)
}
