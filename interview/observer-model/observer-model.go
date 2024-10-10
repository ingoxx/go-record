package main

import "fmt"

// 观察者模式：它允许你定义一个订阅机制，可以在对象事件发生时通知多个“观察”该对象的其他对象。
// 在Go语言中，可以使用通道（channel）或者通过维护一个观察者列表来实现这一模式。
//下面是一个简单的观察者模式示例，我们将创建一个Subject（主题）和一些Observer（观察者），每当主题发生变化时，所有注册的观察者都会得到通知

// Subject 主题接口，定义了添加、删除观察者以及通知观察者的方法
type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyObservers(message string)
}

// ConcreteSubject 具体的主题实现
type ConcreteSubject struct {
	observers []Observer
}

func (s *ConcreteSubject) RegisterObserver(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *ConcreteSubject) RemoveObserver(observer Observer) {
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *ConcreteSubject) NotifyObservers(message string) {
	for _, observer := range s.observers {
		observer.Update(message)
	}
}

// Observer 观察者接口，定义了更新方法
type Observer interface {
	Update(message string)
}

// ConcreteObserver 具体的观察者实现
type ConcreteObserver struct {
	name string
}

func (o *ConcreteObserver) Update(message string) {
	fmt.Printf("%s received message: %s\n", o.name, message)
}

func main() {
	subject := &ConcreteSubject{}

	observer1 := &ConcreteObserver{name: "Observer 1"}
	observer2 := &ConcreteObserver{name: "Observer 2"}

	subject.RegisterObserver(observer1)
	subject.RegisterObserver(observer2)

	subject.NotifyObservers("Hello, observers!")

	// 移除一个观察者
	subject.RemoveObserver(observer1)

	subject.NotifyObservers("Another update...")
}
