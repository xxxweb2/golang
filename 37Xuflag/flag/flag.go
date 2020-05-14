package xuFlag

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2).
	PanicOnError
)

type Value interface {
	String() string
	Set(string) error
}

type Flag struct {
	Name     string
	Usage    string
	Value    Value
	DefValue string
}

type FlagSet struct {
	Usage         func()
	name          string
	parsed        bool
	actual        map[string]*Flag
	formal        map[string]*Flag
	args          []string
	errorHandling ErrorHandling
	output        io.Writer
}

var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVar(p, name, value, usage)

	return p
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.Var(newIntValue(value, p), name, usage)
}

func (f *FlagSet) Var(value Value, name string, usage string) {
	flag := &Flag{
		name, usage, value, value.String(),
	}

	//判断是否有重复的
	_, alreadythere := f.formal[name]
	if alreadythere {
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		fmt.Fprintln(f.Output(), msg)
		panic(msg)
	}

	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}

	f.formal[name] = flag
}

func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}

	return f.output
}

func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = intValue(v)

	return err
}

var errRange = errors.New("value out of range")
var errParse = errors.New("parse error")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}

	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}
	return err

}

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val

	return (*intValue)(p)
}

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
	f.Usage = f.defaultUsage

	return f
}

func (f *FlagSet) defaultUsage() {
}

func Parse() {
	CommandLine.Parse(os.Args[1:])
}

func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments

	for {
		seen, err := f.parseOne()

		if seen {
			continue
		}
		if err == nil {
			break
		}

		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}

	return nil
}

func (f *FlagSet) parseOne() (bool, error) {
	//没有参数结束
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 {
			f.args = f.args[1:]
			return false, nil
		}
	}

	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, f.failf("bad flag syntax: %s", s)
	}

	f.args = f.args[1:]
	hasValue := false
	value := ""

	for i := 1; i < len(name); i++ {
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	m := f.formal
	flag, alreadythree := m[name]
	if !alreadythree {
		if name == "help" || name == "h" { // special case for nice help message.
			f.usage()
			return false, ErrHelp
		}

		return false, f.failf("flag provided but not defined: -%s", name)
	}

	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		if !hasValue && len(f.args) > 0 {
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, f.failf("flag needs an argument: -%s", name)
		}
		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag

	return true, nil

}

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

var ErrHelp = errors.New("flag: help requested")

func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Output(), err)
	f.usage()
	return err
}
func (f *FlagSet) usage() {
	if f.Usage == nil {
		f.defaultUsage()
	} else {
		f.Usage()
	}
}
