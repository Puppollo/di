package di

import "reflect"

func value(value interface{}, cb func() (interface{}, error)) error {
	rvalue := reflect.ValueOf(value)
	if rvalue.Kind() != reflect.Ptr {
		return ErrWrongValue
	}
	rvalue = reflect.Indirect(rvalue)

	v, err := cb()
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	if !rv.Type().AssignableTo(rvalue.Type()) {
		return ErrWrongCallbackReturn
	}
	rvalue.Set(rv)
	rvalue = rvalue.Addr()
	return nil
}

func build(builder reflect.Value, in []reflect.Value) (*reflect.Value, error) {
	out := builder.Call(in)
	if out[1].Interface() != nil {
		err := out[1].Interface().(error)
		if err != nil {
			return nil, err
		}
	}

	return &out[0], nil
}

func (c *DI) build(name string) (*reflect.Value, error) {
	c.logger.Info("build", name)
	service, ok := c.storage[name]
	if !ok {
		return nil, ErrNotFound
	}

	if service.build != nil {
		return service.build, nil
	}

	c.logger.Debug("empty build")
	paramCount := service.builder.Type().NumIn()
	//params := []reflect.Value{service.config}
	var params []reflect.Value
	if paramCount > 0 {
		for i := 0; i < paramCount; i++ {
			if configParam(service.builder.Type().In(i), service.config) {
				params = append(params, *service.config)
				continue
			}
			paramName, err := c.buildParamName(service.builder.Type().In(i), name, service.deps)
			c.logger.Debug("param name:", paramName, ",for argument N:", i)
			if err != nil {
				return nil, err
			}
			param, err := c.build(paramName)
			if err != nil {
				return nil, err
			}
			params = append(params, *param)
		}
	}

	value, err := build(service.builder, params)
	if err != nil {
		return nil, err
	}
	service.build = value
	return value, nil
}

func configParam(in reflect.Type, config *reflect.Value) bool {
	if config == nil {
		return false
	}
	if in == config.Type() {
		return true
	}
	return false
}

func (c *DI) buildParamName(in reflect.Type, skip string, deps map[reflect.Type]string) (string, error) {
	for name, service := range c.storage {
		if name == skip {
			continue
		}
		for i, n := range deps {
			if in == i {
				return n, nil
			}
		}
		if !service.out.Implements(in) {
			continue
		}
		return name, nil
	}
	return "", ErrParamBuilderNotFound
}
