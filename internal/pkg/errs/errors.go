package errs

import "fmt"

type ObjectNotFoundError struct {
	msg string
}

func NewObjectNotFoundError(msg string) ObjectNotFoundError {
	return ObjectNotFoundError{msg: msg}
}

func (e ObjectNotFoundError) Error() string {
	return fmt.Sprintf("object not found %s", e.msg)
}

type ValueIsInvalidError struct {
	msg string
}

func NewValueIsInvalidError(msg string) ValueIsInvalidError {
	return ValueIsInvalidError{msg: msg}
}

func (v ValueIsInvalidError) Error() string {
	return fmt.Sprintf("value is invalid %s", v.msg)
}

type ValueIsRequiredError struct {
	msg string
}

func NewValueIsRequiredError(msg string) ValueIsRequiredError {
	return ValueIsRequiredError{msg: msg}
}

func (v ValueIsRequiredError) Error() string {
	return fmt.Sprintf("value is required %s", v.msg)
}

type VersionIsInvalidError struct {
	msg string
}

func NewVersionIsInvalidError(msg string) VersionIsInvalidError {
	return VersionIsInvalidError{msg: msg}
}

func (v VersionIsInvalidError) Error() string {
	return fmt.Sprintf("version is invalid %s", v.msg)
}
