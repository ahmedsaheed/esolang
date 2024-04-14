package object

import (
	"fmt"
)

type CheckFunc func(name string, args []Object) error

func CheckTypings(name string, args []Object, checks ...CheckFunc) error {
	for _, check := range checks {
		if err := check(name, args); err != nil {
			return err
		}
	}
	return nil
}

func ExactArgsLength(n int) CheckFunc {
	return func(name string, args []Object) error {
		if len(args) != n {
			return fmt.Errorf(
				"TypeError: %s() takes exactly %d argument (%d given)",
				name, n, len(args),
			)
		}
		return nil
	}
}

func MinimumArgs(n int) CheckFunc {
	return func(name string, args []Object) error {
		if len(args) < n {
			return fmt.Errorf(
				"TypeError: %s() takes a minimum %d arguments (%d given)",
				name, n, len(args),
			)
		}
		return nil
	}
}

func RangeOfArgs(n, m int) CheckFunc {
	return func(name string, args []Object) error {
		if len(args) < n || len(args) > m {
			return fmt.Errorf(
				"TypeError: %s() takes at least %d arguments at most %d (%d given)",
				name, n, m, len(args),
			)
		}
		return nil
	}
}

func WithTypes(types ...ObjectType) CheckFunc {
	return func(name string, args []Object) error {
		for i, t := range types {
			if i < len(args) && args[i].Type() != t {
				return fmt.Errorf(
					"TypeError: %s() expected argument #%d to be `%s` got `%s`",
					name, (i + 1), t, args[i].Type(),
				)
			}
		}
		return nil
	}
}
