package core

import "fmt"

type SdkError struct {
	general  string
	specific error
}

func (s SdkError) Error() string {
	return fmt.Sprintf("%s: %s", s.general, s.specific.Error())
}
