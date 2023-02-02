package main

import (
	"fmt"
	// "errors"
	"github.com/pkg/errors"
)

type CustomError struct {
  Code int
}

func (c *CustomError) Error() string {
  return fmt.Sprintf("Failed with code %d", c.Code)
}

func main() {
	// result, err := caller1()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }

	// fmt.Println("Result: ", result)

	err := &CustomError{Code: 12}
	// lostErr := fmt.Errorf("failed with error: %v", err)
	// there is no way we can get back the `Code` attribute from `lostErr`

	wrappedErr := errors.Wrap(err, "[1] failed with error:")
	twiceWrappedError := errors.Wrap(wrappedErr, "[2] failed with error:")

  // The `errors.Cause` function returns the originally wrapped error, which we can then type assert to its original struct type
	if originalErr, ok := errors.Cause(twiceWrappedError).(*CustomError); ok {
		fmt.Println("the original error coed was : ", originalErr.Code)
	}
}

func caller1() (int, error) {
	err := caller2()
	if err != nil {
		return 0, fmt.Errorf("[caller1] error in calling caller2: %v", err)
	}

	return 1, nil
}

func caller2() error {
	err := caller3()
	if err != nil {
		return fmt.Errorf("[caller2] error in calling caller3: %v", err)
	}

	return nil
}


func caller3() error {
	return errors.New("failed")
}
