// This file was generated by counterfeiter
package cowbullfakes

import (
	"sync"

	"github.com/Bo0mer/cowbull"
)

type FakeMessenger struct {
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct{}
	iDReturns     struct {
		result1 string
	}
	OnMessageStub        func(kind string, action func(data string))
	onMessageMutex       sync.RWMutex
	onMessageArgsForCall []struct {
		kind   string
		action func(data string)
	}
	SendMessageStub        func(kind string, data string) error
	sendMessageMutex       sync.RWMutex
	sendMessageArgsForCall []struct {
		kind string
		data string
	}
	sendMessageReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMessenger) ID() string {
	fake.iDMutex.Lock()
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct{}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	} else {
		return fake.iDReturns.result1
	}
}

func (fake *FakeMessenger) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeMessenger) IDReturns(result1 string) {
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeMessenger) OnMessage(kind string, action func(data string)) {
	fake.onMessageMutex.Lock()
	fake.onMessageArgsForCall = append(fake.onMessageArgsForCall, struct {
		kind   string
		action func(data string)
	}{kind, action})
	fake.recordInvocation("OnMessage", []interface{}{kind, action})
	fake.onMessageMutex.Unlock()
	if fake.OnMessageStub != nil {
		fake.OnMessageStub(kind, action)
	}
}

func (fake *FakeMessenger) OnMessageCallCount() int {
	fake.onMessageMutex.RLock()
	defer fake.onMessageMutex.RUnlock()
	return len(fake.onMessageArgsForCall)
}

func (fake *FakeMessenger) OnMessageArgsForCall(i int) (string, func(data string)) {
	fake.onMessageMutex.RLock()
	defer fake.onMessageMutex.RUnlock()
	return fake.onMessageArgsForCall[i].kind, fake.onMessageArgsForCall[i].action
}

func (fake *FakeMessenger) SendMessage(kind string, data string) error {
	fake.sendMessageMutex.Lock()
	fake.sendMessageArgsForCall = append(fake.sendMessageArgsForCall, struct {
		kind string
		data string
	}{kind, data})
	fake.recordInvocation("SendMessage", []interface{}{kind, data})
	fake.sendMessageMutex.Unlock()
	if fake.SendMessageStub != nil {
		return fake.SendMessageStub(kind, data)
	} else {
		return fake.sendMessageReturns.result1
	}
}

func (fake *FakeMessenger) SendMessageCallCount() int {
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	return len(fake.sendMessageArgsForCall)
}

func (fake *FakeMessenger) SendMessageArgsForCall(i int) (string, string) {
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	return fake.sendMessageArgsForCall[i].kind, fake.sendMessageArgsForCall[i].data
}

func (fake *FakeMessenger) SendMessageReturns(result1 error) {
	fake.SendMessageStub = nil
	fake.sendMessageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMessenger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.onMessageMutex.RLock()
	defer fake.onMessageMutex.RUnlock()
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeMessenger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ cowbull.Messenger = new(FakeMessenger)