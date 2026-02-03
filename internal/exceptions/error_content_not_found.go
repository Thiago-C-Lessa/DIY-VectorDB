package exceptions

import "fmt"

type ErrorContentNotFound struct {
	Key string
}

func (e *ErrorContentNotFound) Error() string {
	return fmt.Sprintf("The key : %s does not have content stored", e.Key)
}
