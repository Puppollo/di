# di
simple di container for learning purposes
service builder - func(<config struct>, ...<dep interfaces>) (<service instance>, error)

can build 
```go
buildSomeA(config ConfigA) (someA, error) {...}
buildSomeB(config ConfigB, interfaceB) (someB, error) {...}
buildSomeC(config ConfigC, interfaceA, interfaceB) (someC, error) {...}
```
can't
```go
buildSome(config Config, interfaceA, interfaceA) (some, error) {...}
```