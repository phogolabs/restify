// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"net/http"
	"sync"
)

type ResponseWriter struct {
	HeaderStub        func() http.Header
	headerMutex       sync.RWMutex
	headerArgsForCall []struct {
	}
	headerReturns struct {
		result1 http.Header
	}
	headerReturnsOnCall map[int]struct {
		result1 http.Header
	}
	WriteStub        func([]byte) (int, error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		arg1 []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	WriteHeaderStub        func(int)
	writeHeaderMutex       sync.RWMutex
	writeHeaderArgsForCall []struct {
		arg1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ResponseWriter) Header() http.Header {
	fake.headerMutex.Lock()
	ret, specificReturn := fake.headerReturnsOnCall[len(fake.headerArgsForCall)]
	fake.headerArgsForCall = append(fake.headerArgsForCall, struct {
	}{})
	fake.recordInvocation("Header", []interface{}{})
	fake.headerMutex.Unlock()
	if fake.HeaderStub != nil {
		return fake.HeaderStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.headerReturns
	return fakeReturns.result1
}

func (fake *ResponseWriter) HeaderCallCount() int {
	fake.headerMutex.RLock()
	defer fake.headerMutex.RUnlock()
	return len(fake.headerArgsForCall)
}

func (fake *ResponseWriter) HeaderCalls(stub func() http.Header) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = stub
}

func (fake *ResponseWriter) HeaderReturns(result1 http.Header) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = nil
	fake.headerReturns = struct {
		result1 http.Header
	}{result1}
}

func (fake *ResponseWriter) HeaderReturnsOnCall(i int, result1 http.Header) {
	fake.headerMutex.Lock()
	defer fake.headerMutex.Unlock()
	fake.HeaderStub = nil
	if fake.headerReturnsOnCall == nil {
		fake.headerReturnsOnCall = make(map[int]struct {
			result1 http.Header
		})
	}
	fake.headerReturnsOnCall[i] = struct {
		result1 http.Header
	}{result1}
}

func (fake *ResponseWriter) Write(arg1 []byte) (int, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("Write", []interface{}{arg1Copy})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.writeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ResponseWriter) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *ResponseWriter) WriteCalls(stub func([]byte) (int, error)) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = stub
}

func (fake *ResponseWriter) WriteArgsForCall(i int) []byte {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	argsForCall := fake.writeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ResponseWriter) WriteReturns(result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *ResponseWriter) WriteReturnsOnCall(i int, result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = nil
	if fake.writeReturnsOnCall == nil {
		fake.writeReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.writeReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *ResponseWriter) WriteHeader(arg1 int) {
	fake.writeHeaderMutex.Lock()
	fake.writeHeaderArgsForCall = append(fake.writeHeaderArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("WriteHeader", []interface{}{arg1})
	fake.writeHeaderMutex.Unlock()
	if fake.WriteHeaderStub != nil {
		fake.WriteHeaderStub(arg1)
	}
}

func (fake *ResponseWriter) WriteHeaderCallCount() int {
	fake.writeHeaderMutex.RLock()
	defer fake.writeHeaderMutex.RUnlock()
	return len(fake.writeHeaderArgsForCall)
}

func (fake *ResponseWriter) WriteHeaderCalls(stub func(int)) {
	fake.writeHeaderMutex.Lock()
	defer fake.writeHeaderMutex.Unlock()
	fake.WriteHeaderStub = stub
}

func (fake *ResponseWriter) WriteHeaderArgsForCall(i int) int {
	fake.writeHeaderMutex.RLock()
	defer fake.writeHeaderMutex.RUnlock()
	argsForCall := fake.writeHeaderArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ResponseWriter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.headerMutex.RLock()
	defer fake.headerMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	fake.writeHeaderMutex.RLock()
	defer fake.writeHeaderMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ResponseWriter) recordInvocation(key string, args []interface{}) {
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

var _ http.ResponseWriter = new(ResponseWriter)
