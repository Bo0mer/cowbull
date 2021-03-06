// This file was generated by counterfeiter
package cowbullfakes

import (
	"net"
	"sync"
	"time"

	"github.com/Bo0mer/cowbull"
)

type FakeConn struct {
	SetReadDeadlineStub        func(t time.Time) error
	setReadDeadlineMutex       sync.RWMutex
	setReadDeadlineArgsForCall []struct {
		t time.Time
	}
	setReadDeadlineReturns struct {
		result1 error
	}
	WriteJSONStub        func(v interface{}) error
	writeJSONMutex       sync.RWMutex
	writeJSONArgsForCall []struct {
		v interface{}
	}
	writeJSONReturns struct {
		result1 error
	}
	ReadJSONStub        func(v interface{}) error
	readJSONMutex       sync.RWMutex
	readJSONArgsForCall []struct {
		v interface{}
	}
	readJSONReturns struct {
		result1 error
	}
	RemoteAddrStub        func() net.Addr
	remoteAddrMutex       sync.RWMutex
	remoteAddrArgsForCall []struct{}
	remoteAddrReturns     struct {
		result1 net.Addr
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConn) SetReadDeadline(t time.Time) error {
	fake.setReadDeadlineMutex.Lock()
	fake.setReadDeadlineArgsForCall = append(fake.setReadDeadlineArgsForCall, struct {
		t time.Time
	}{t})
	fake.recordInvocation("SetReadDeadline", []interface{}{t})
	fake.setReadDeadlineMutex.Unlock()
	if fake.SetReadDeadlineStub != nil {
		return fake.SetReadDeadlineStub(t)
	} else {
		return fake.setReadDeadlineReturns.result1
	}
}

func (fake *FakeConn) SetReadDeadlineCallCount() int {
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	return len(fake.setReadDeadlineArgsForCall)
}

func (fake *FakeConn) SetReadDeadlineArgsForCall(i int) time.Time {
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	return fake.setReadDeadlineArgsForCall[i].t
}

func (fake *FakeConn) SetReadDeadlineReturns(result1 error) {
	fake.SetReadDeadlineStub = nil
	fake.setReadDeadlineReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) WriteJSON(v interface{}) error {
	fake.writeJSONMutex.Lock()
	fake.writeJSONArgsForCall = append(fake.writeJSONArgsForCall, struct {
		v interface{}
	}{v})
	fake.recordInvocation("WriteJSON", []interface{}{v})
	fake.writeJSONMutex.Unlock()
	if fake.WriteJSONStub != nil {
		return fake.WriteJSONStub(v)
	} else {
		return fake.writeJSONReturns.result1
	}
}

func (fake *FakeConn) WriteJSONCallCount() int {
	fake.writeJSONMutex.RLock()
	defer fake.writeJSONMutex.RUnlock()
	return len(fake.writeJSONArgsForCall)
}

func (fake *FakeConn) WriteJSONArgsForCall(i int) interface{} {
	fake.writeJSONMutex.RLock()
	defer fake.writeJSONMutex.RUnlock()
	return fake.writeJSONArgsForCall[i].v
}

func (fake *FakeConn) WriteJSONReturns(result1 error) {
	fake.WriteJSONStub = nil
	fake.writeJSONReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) ReadJSON(v interface{}) error {
	fake.readJSONMutex.Lock()
	fake.readJSONArgsForCall = append(fake.readJSONArgsForCall, struct {
		v interface{}
	}{v})
	fake.recordInvocation("ReadJSON", []interface{}{v})
	fake.readJSONMutex.Unlock()
	if fake.ReadJSONStub != nil {
		return fake.ReadJSONStub(v)
	} else {
		return fake.readJSONReturns.result1
	}
}

func (fake *FakeConn) ReadJSONCallCount() int {
	fake.readJSONMutex.RLock()
	defer fake.readJSONMutex.RUnlock()
	return len(fake.readJSONArgsForCall)
}

func (fake *FakeConn) ReadJSONArgsForCall(i int) interface{} {
	fake.readJSONMutex.RLock()
	defer fake.readJSONMutex.RUnlock()
	return fake.readJSONArgsForCall[i].v
}

func (fake *FakeConn) ReadJSONReturns(result1 error) {
	fake.ReadJSONStub = nil
	fake.readJSONReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) RemoteAddr() net.Addr {
	fake.remoteAddrMutex.Lock()
	fake.remoteAddrArgsForCall = append(fake.remoteAddrArgsForCall, struct{}{})
	fake.recordInvocation("RemoteAddr", []interface{}{})
	fake.remoteAddrMutex.Unlock()
	if fake.RemoteAddrStub != nil {
		return fake.RemoteAddrStub()
	} else {
		return fake.remoteAddrReturns.result1
	}
}

func (fake *FakeConn) RemoteAddrCallCount() int {
	fake.remoteAddrMutex.RLock()
	defer fake.remoteAddrMutex.RUnlock()
	return len(fake.remoteAddrArgsForCall)
}

func (fake *FakeConn) RemoteAddrReturns(result1 net.Addr) {
	fake.RemoteAddrStub = nil
	fake.remoteAddrReturns = struct {
		result1 net.Addr
	}{result1}
}

func (fake *FakeConn) Close() error {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	} else {
		return fake.closeReturns.result1
	}
}

func (fake *FakeConn) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeConn) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	fake.writeJSONMutex.RLock()
	defer fake.writeJSONMutex.RUnlock()
	fake.readJSONMutex.RLock()
	defer fake.readJSONMutex.RUnlock()
	fake.remoteAddrMutex.RLock()
	defer fake.remoteAddrMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeConn) recordInvocation(key string, args []interface{}) {
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

var _ cowbull.Conn = new(FakeConn)
