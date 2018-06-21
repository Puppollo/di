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
	params := make([]reflect.Value, 0, paramCount)
	if paramCount > 0 {
		params = append(params, service.config)
		for i := 1; i < paramCount; i++ {
			paramName, err := c.buildParamName(service.builder.Type().In(i), name, service.deps)
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

func (c *DI) buildParamName(in reflect.Type, skip string, deps map[reflect.Type]string) (string, error) {
	c.logger.Debug("look for param name")
	for name, service := range c.storage {
		if name == skip {
			continue
		}
		for i, n := range deps {
			c.logger.Debug("deps check", n)
			c.logger.Debug(in, i)
			if in == i {
				c.logger.Debug("return name deps", n)
				return n, nil
			}
		}
		c.logger.Debug("check", name)
		if !service.out.Implements(in) {
			continue
		}
		return name, nil
	}
	return "", ErrParamBuilderNotFound
}
