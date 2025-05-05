package example

import "fmt"

type ServiceA struct{}

func NewServiceA() ServiceA {
	return ServiceA{}
}

type IServiceB interface {
	Echo()
}

type ServiceB struct{}

func (s ServiceB) Echo() {
	fmt.Println("Hello from ServiceB")
}

func NewServiceB() IServiceB {
	return &ServiceB{}
}

func main() {
	var b IServiceB = NewServiceB() // should not trigger linter
	b.Echo()

	a := NewServiceA() // should trigger linter
	_ = a

}
