package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

func PrintError(err error) {
	if err != nil {
		fmt.Printf("error = [%+v]\n", err)
		panic(err)
	}
}

func PrintWrapError(err error, msg string) {
	if err != nil {
		fmt.Printf("error = [%+v]\n", errors.Wrap(err, msg))
		panic(err)
	}
}
