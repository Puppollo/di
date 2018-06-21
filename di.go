package di

import (
	"reflect"
	"sync"
)

var errorInterface = reflect.TypeOf(new(error)).Elem()

type (
	Deps map[interface{}]string

	service struct {
		config  reflect.Value           // service config
		builder reflect.Value           // service builder func
		deps    map[reflect.Type]string // dependencies
		out     reflect.Type            // builder(...) (out, error)
		build   *reflect.Value          // builded service
	}

	DI struct {
		logger  Logger
		storage map[string]*service
		mutex   *sync.Mutex
	}
)

func New() *DI {
	return &DI{
		storage: make(map[string]*service),
		mutex:   &sync.Mutex{},
		logger:  NullLogger(),
	}
}

func (c *DI) SetLogger(l Logger) {
	c.logger = l
}

func (c *DI) Build(name string, v interface{}) error {
	return value(v, func() (interface{}, error) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		b, err := c.build(name)
		if err != nil {
			return nil, err
		}
		return b.Interface(), nil
	})
}

// name - service name in container
// config - should be a struct, service config values
// builder - func that create an instance `func(<config>, ...<dependencies>) (<service>, error)`
// deps - map that points to specific dependency implementation by name
func (c *DI) Add(name string, config interface{}, builder interface{}, deps map[interface{}]string) error {
	c.logger.Debug("add", name)
	cfgValue := reflect.ValueOf(config)
	if cfgValue.Kind() != reflect.Struct {
		return ErrWrongConfigType
	}
	bldValue := reflect.ValueOf(builder)
	if bldValue.Kind() != reflect.Func {
		return ErrWrongBuilderType
	}
	if bldValue.Type().NumOut() != 2 {
		return ErrWrongBuilderReturn
	}
	//if bldValue.Type().Out(0).Kind() != reflect.Struct {
	//	return ErrWrongBuilderReturn
	//}
	if !bldValue.Type().Out(1).Implements(errorInterface) {
		return ErrWrongBuilderReturn
	}
	d := make(map[reflect.Type]string, len(deps))
	for t, n := range deps {
		tt := reflect.TypeOf(t)
		c.logger.Debug("type", tt)
		if tt.Kind() == reflect.Ptr {
			c.logger.Debug("elem", tt)
			tt = tt.Elem()
		}
		d[tt] = n
	}
	c.mutex.Lock()
	c.storage[name] = &service{
		config:  cfgValue,
		builder: bldValue,
		deps:    d,
		out:     bldValue.Type().Out(0),
	}
	c.mutex.Unlock()
	return nil
}
