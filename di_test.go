package di_test

import (
	"fmt"
	"github.com/Puppollo/di"
	"testing"
)

type (
	Inter interface {
		Int() int
	}

	Stringer interface {
		String() string
	}

	StringSum interface {
		Sum() string
	}

	IntSum interface {
		Sum() int
	}

	configA struct{ v int }
	configB struct{ v string }

	A struct{ v int }
	B struct{ v string }
)

func NewA(c configA) (*A, error) {
	return &A{v: c.v}, nil
}

func TestDI(t *testing.T) {
	d := di.New()
	d.SetLogger(di.ConsoleLogger())

	tests := []struct {
		title   string
		name    string
		config  interface{}
		builder interface{}
		deps    di.Deps
		err     error
		check   func(name string) error
	}{
		{
			title:   "nil builder",
			name:    "a",
			config:  nil,
			builder: nil,
			deps:    nil,
			err:     di.ErrWrongConfigType,
			check:   func(name string) error { return nil },
		},
		{
			title:   "wrong builder type",
			name:    "b",
			config:  struct{}{},
			builder: 10,
			deps:    nil,
			err:     di.ErrWrongBuilderType,
			check:   func(name string) error { return nil },
		},
		{
			title:   "wrong builder return",
			name:    "c",
			config:  struct{}{},
			builder: func() {},
			deps:    nil,
			err:     di.ErrWrongBuilderReturn,
			check:   func(name string) error { return nil },
		},
		{
			title:   "wrong build",
			name:    "d",
			config:  configA{v: 10},
			builder: NewA,
			deps:    nil,
			err:     nil,
			check: func(name string) error {
				var a A
				err := d.Build(name, &a)
				if err != nil {
					return err
				}
				if a.v != 10 {
					return fmt.Errorf("expect a.v:\n\t%v\ngot:\n\t%v\n", 10, a.v)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := d.Add(test.name, test.config, test.builder, test.deps)
			if err != test.err {
				t.Errorf("expect err:\n\t%v\ngot:\n\t%v\n", test.err, err)
			}
			if err != nil {
				return
			}
			err = test.check(test.name)
			if err != nil {
				t.Error("checker error:", err)
			}
		})
	}
}
