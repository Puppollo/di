package di_test

import (
	"github.com/Puppollo/di"

	"fmt"
	"testing"
)

type (
	Inter interface {
		Int() int
	}

	config struct{ v int }
	A      struct{ v int }
	B      struct {
		v     int
		inter Inter
	}

	Some struct{ s int }
)

const (
	expectErrFmt   = "expect err:\n\t%v\ngot:\n\t%v\n"
	expectValueFmt = "expect value:\n\t%v\ngot:\n\t%v\n"
)

func NewA(c config) (*A, error) { return &A{v: c.v}, nil }
func (a *A) Int() int           { return a.v }

func NewB(c config, i Inter) (*B, error)    { return &B{v: c.v, inter: i}, nil }
func NewBAlt(i Inter, c config) (*B, error) { return &B{v: c.v, inter: i}, nil }
func NewBConfigless(i Inter) (*B, error)    { return &B{inter: i}, nil }
func (b *B) Int() int                       { return b.v }

func (b *B) Sum() int { return b.v + b.inter.Int() }

func (s Some) Some() {
	println("some", s.s)
}

func TestDI_AddBuild(t *testing.T) {
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
			title: "nil builder",
			name:  "1",
			err:   di.ErrWrongBuilderType,
			check: func(name string) error { return nil },
		},
		{
			title:   "wrong builder type",
			name:    "2",
			config:  struct{}{},
			builder: 10,
			err:     di.ErrWrongBuilderType,
			check:   func(name string) error { return nil },
		},
		{
			title:   "wrong builder return",
			name:    "3",
			config:  struct{}{},
			builder: func() {},
			err:     di.ErrWrongBuilderReturn,
			check:   func(name string) error { return nil },
		},
		{
			title:   "can't find dep service to create",
			name:    "b",
			config:  config{v: 10},
			builder: NewB,
			check: func(name string) error {
				expectErr := di.ErrParamBuilderNotFound
				var b struct{}
				err := d.Build(name, &b)
				if err == expectErr {
					return nil
				}
				if err == nil {
					return fmt.Errorf(expectErrFmt, expectErr, err)
				}
				return err
			},
		},
		{
			title:   "wrong build value receiver",
			name:    "a",
			config:  config{v: 10},
			builder: NewA,
			check: func(name string) error {
				expectErr := di.ErrWrongCallbackReturn
				var b struct{}
				err := d.Build(name, &b)
				if err == expectErr {
					return nil
				}
				if err == nil {
					return fmt.Errorf(expectErrFmt, expectErr, err)
				}
				return err
			},
		},
		{
			title: "wrong build value, not a pointer",
			check: func(_ string) error {
				expectErr := di.ErrWrongValue
				var a A
				err := d.Build("d", a)
				if err == expectErr {
					return nil
				}
				if err == nil {
					return fmt.Errorf(expectErrFmt, expectErr, err)
				}
				return err
			},
		},
		{
			title: "build a",
			check: func(_ string) error {
				var a A
				err := d.Build("a", &a)
				if err != nil {
					return err
				}
				if a.v != 10 {
					return fmt.Errorf(expectValueFmt, 10, a.v)
				}
				return nil
			},
		},
		{
			title:   "build b",
			name:    "b",
			config:  config{v: 5},
			builder: NewB,
			check: func(name string) error {
				var b B
				err := d.Build(name, &b)
				if err != nil {
					return err
				}
				if b.Int() != 5 {
					return fmt.Errorf(expectValueFmt, 5, b.Int())
				}
				if b.Sum() != 15 {
					return fmt.Errorf(expectValueFmt, 15, b.Int())
				}
				return nil
			},
		},
		{
			title:   "build b from other b",
			name:    "bb",
			config:  config{v: 7},
			builder: NewB,
			deps:    di.Deps{new(Inter): "b"},
			check: func(name string) error {
				var b B
				err := d.Build(name, &b)
				if err != nil {
					return err
				}
				if b.Int() != 7 {
					return fmt.Errorf("on Int\n"+expectValueFmt, 7, b.Int())
				}
				if b.Sum() != 12 {
					fmt.Printf("%#v\n", b)
					return fmt.Errorf("on Sum\n"+expectValueFmt, 12, b.Sum())
				}
				return nil
			},
		},
		{
			title:   "build b from a",
			name:    "ba",
			config:  config{v: 3},
			builder: NewB,
			deps:    di.Deps{new(Inter): "a"},
			check: func(name string) error {
				var b B
				err := d.Build(name, &b)
				if err != nil {
					return err
				}
				if b.Int() != 3 {
					return fmt.Errorf("on Int\n"+expectValueFmt, 3, b.Int())
				}
				if b.Sum() != 13 {
					fmt.Printf("%#v\n", b)
					return fmt.Errorf("on Sum\n"+expectValueFmt, 13, b.Sum())
				}
				return nil
			},
		},
		{
			title:   "build b from b configless",
			name:    "bb-configless",
			builder: NewBConfigless,
			deps:    di.Deps{new(Inter): "b"},
			check: func(name string) error {
				var b B
				err := d.Build(name, &b)
				if err != nil {
					return err
				}
				if b.Int() != 0 {
					return fmt.Errorf("on Int\n"+expectValueFmt, 0, b.Int())
				}
				if b.Sum() != 5 {
					fmt.Printf("%#v\n", b)
					return fmt.Errorf("on Sum\n"+expectValueFmt, 5, b.Sum())
				}
				return nil
			},
		},
		{
			title:   "build b from other b with alt builder",
			name:    "bb",
			config:  config{v: 70},
			builder: NewBAlt,
			deps:    di.Deps{new(Inter): "b"},
			check: func(name string) error {
				var b B
				err := d.Build(name, &b)
				if err != nil {
					return err
				}
				if b.Int() != 70 {
					return fmt.Errorf("on Int\n"+expectValueFmt, 70, b.Int())
				}
				if b.Sum() != 75 {
					fmt.Printf("%#v\n", b)
					return fmt.Errorf("on Sum\n"+expectValueFmt, 75, b.Sum())
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.name != "" {
				err := d.Add(test.name, test.config, test.builder, test.deps)
				if err != test.err {
					t.Errorf("expect err:\n\t%v\ngot:\n\t%v\n", test.err, err)
				}
				if err != nil {
					return
				}
			}
			err := test.check(test.name)
			if err != nil {
				t.Error("checker error:", err)
			}
		})
	}
}

func TestDI_Some(t *testing.T) {
	container := di.New()
	err := container.Add("some", nil, func() (Some, error) { return Some{203}, nil }, nil)
	if err != nil {
		t.Fatal(err)
	}
	var s Some
	err = container.Build("some", &s)
	if err != nil {
		t.Fatal(err)
	}
	s.Some()
}
