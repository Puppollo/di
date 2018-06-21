# di
simple di container for learning purposes
service builder - func(<config struct>, ...<dep interfaces>) (<service instance>, error)

can build 
```go
buildSomeA(config ConfigA) (someA, error) {...}
buildSomeB(config ConfigB, interfaceB) (someB, error) {...}
buildSomeC(config ConfigC, interfaceA, interfaceB) (someC, error) {...}
buildSomeD() (someD, error) {...}
buildSomeE(interfaceA) (someE, error) {...}
buildSomeF(interfaceA, configF) (someF, error) {...}
```
can't
```go
buildSome(config Config, interfaceA, interfaceA) (some, error) {...}
buildSome(config Config, b structB, interfaceA) (some, error) {...}
buildSome(<literal> C, b structB) (some, error) {...}
```

using

```go
container := di.New()
err := container.Add("some", someconfig{}, NewSome)
if err!= nil {
	panic(err)
}
var s some
err := container.Build("some", &s)
if err!= nil {
	panic(err)
}
s.Some()
```