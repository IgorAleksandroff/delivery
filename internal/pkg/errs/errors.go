package errs

import (
	"errors"
	"fmt"
	"strings"
)

var ErrObjectNotFound = errors.New("object not found")
var ErrValueIsInvalid = errors.New("value is invalid")
var ErrValueIsRequired = errors.New("value is required")
var ErrVersionIsInvalid = errors.New("version is invalid")
var ErrValueIsOutOfRange = errors.New("value is out of range")

type ObjectNotFoundError struct {
	ParamName string
	ID        any
}

func NewObjectNotFoundError(paramName string, ID any) *ObjectNotFoundError {
	return &ObjectNotFoundError{
		ParamName: paramName,
		ID:        ID,
	}
}

func (e *ObjectNotFoundError) Error() string {
	return fmt.Sprintf("%s: %s", ErrObjectNotFound, e.ID)
}

func (e *ObjectNotFoundError) Unwrap() error {
	return ErrObjectNotFound
}

type ValueIsInvalidError struct {
	ParamName string
}

func NewValueIsInvalidError(paramName string) *ValueIsInvalidError {
	return &ValueIsInvalidError{
		ParamName: paramName,
	}
}

func (e *ValueIsInvalidError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValueIsInvalid, e.ParamName)
}

func (e *ValueIsInvalidError) Unwrap() error {
	return ErrValueIsInvalid
}

type ValueIsRequiredError struct {
	ParamName string
}

func NewValueIsRequiredError(paramName string) *ValueIsRequiredError {
	return &ValueIsRequiredError{
		ParamName: paramName,
	}
}

func (e *ValueIsRequiredError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValueIsRequired, e.ParamName)
}

func (e *ValueIsRequiredError) Unwrap() error {
	return ErrValueIsRequired
}

type VersionIsInvalidError struct {
	ParamName string
	Cause     error
}

func NewVersionIsInvalidError(paramName string, cause error) *VersionIsInvalidError {
	return &VersionIsInvalidError{
		ParamName: paramName,
		Cause:     cause,
	}
}

func NewVersionIsInvalidErrorWithCause(paramName string) *VersionIsInvalidError {
	return &VersionIsInvalidError{
		ParamName: paramName,
	}
}

func (e *VersionIsInvalidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", ErrVersionIsInvalid, e.ParamName, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrVersionIsInvalid, e.ParamName)
}

func (e *VersionIsInvalidError) Unwrap() error {
	return ErrVersionIsInvalid
}

type ValueIsOutOfRangeError struct {
	ParamName string
	Value     any
	Min       any
	Max       any
}

func NewValueIsOutOfRangeError(paramName string, value any, min any, max any) *ValueIsOutOfRangeError {
	return &ValueIsOutOfRangeError{
		ParamName: paramName,
		Value:     value,
		Min:       min,
		Max:       max,
	}
}

func (e *ValueIsOutOfRangeError) Error() string {
	return fmt.Sprintf("%s: %s is %v, min value is %v, max value is %v",
		ErrValueIsInvalid, sanitize(e.Value), e.ParamName, e.Min, e.Max)
}

func (e *ValueIsOutOfRangeError) Unwrap() error {
	return ErrValueIsOutOfRange
}

func sanitize(input interface{}) string {
	str := fmt.Sprintf("%v", input)
	return strings.ReplaceAll(str, "\n", " ")
}
