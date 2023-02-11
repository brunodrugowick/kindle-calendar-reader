package mocks

// Returns holds the values to return for a mocked call
type Returns struct {
	fields []interface{}
}

// Get returns the given return value. It should be used in your mocks.
func (r *Returns) Get(id int) interface{} {
	return r.fields[id]
}

// Error is a special Get that returns a nil error in case a nil error is expected.
func (r *Returns) Error(id int) error {
	got := r.Get(id)
	if got != nil {
		return got.(error)
	}
	return error(nil)
}

// Mocker is an embeddable struct to put on your mocks. It has methods to set return values and other things
// that for now I'm not using (calls and with).
type Mocker struct {
	calls   map[string]int
	with    map[string][]interface{}
	returns map[string]Returns
}

// On sets the return values for a mocked method
func (m *Mocker) On(method string, returns ...interface{}) Mocker {
	if m.with == nil || m.calls == nil || m.returns == nil {
		m.calls = make(map[string]int)
		m.with = make(map[string][]interface{})
		m.returns = make(map[string]Returns)
	}

	m.calls[method] += 1
	m.returns[method] = Returns{fields: returns}
	return *m
}

// called gets the return values for the mocked thing to use
func (m *Mocker) called(method string, calledWith ...interface{}) Returns {
	m.calls[method] += 1
	m.with[method] = calledWith
	return m.returns[method]
}
