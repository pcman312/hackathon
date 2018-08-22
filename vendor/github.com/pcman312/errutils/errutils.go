package errutils

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func NewMultiError() *multierror.Error {
	return NewMultiErrorFmt(SingleLineFormat)
}

func NewMultiErrorFmt(format multierror.ErrorFormatFunc) *multierror.Error {
	merr := &multierror.Error{
		ErrorFormat: format,
	}
	return merr
}

func SingleLineFormat(errs []error) string {
	return JoinErrsStr(", ", errs...)
}

// JoinErrs combines the given errors into a single error using the delimiter
// If no errors are provided, this will return nil
func JoinErrs(delim string, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.New(JoinErrsStr(delim, errs...))
}

// JoinErrsStr combines the given errors into a single string using the delimiter
// If no errors are provided, this will return an empty string
func JoinErrsStr(delim string, errs ...error) string {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = err.Error()
	}
	return strings.Join(strs, delim)
}
