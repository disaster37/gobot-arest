package extra

import "sync"

type extraTestBareAdaptor struct{}

func (t *extraTestBareAdaptor) Connect() (err error)  { return }
func (t *extraTestBareAdaptor) Finalize() (err error) { return }
func (t *extraTestBareAdaptor) Name() string          { return "" }
func (t *extraTestBareAdaptor) SetName(n string)      {}

type extraTestExtraReader struct {
	extraTestBareAdaptor
}

func (t *extraTestExtraReader) ValueRead(name string) (val interface{}, err error)   { return }
func (t *extraTestExtraReader) ValuesRead() (vals map[string]interface{}, err error) { return }
func (t *extraTestExtraReader) FunctionCall(name string, parameters string) (val int, err error) {
	return
}

type extraTestAdaptor struct {
	name                    string
	port                    string
	mtx                     sync.Mutex
	testAdaptorValueRead    func(name string) (val interface{}, err error)
	testAdaptorValuesRead   func() (vals map[string]interface{}, err error)
	testAdaptorFunctionCall func(name string, parameters string) (val int, err error)
}

func (t *extraTestAdaptor) TestAdaptorValueRead(f func(name string) (val interface{}, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorValueRead = f
}
func (t *extraTestAdaptor) TestAdaptorValuesRead(f func() (vals map[string]interface{}, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorValuesRead = f
}
func (t *extraTestAdaptor) TestAdaptorFunctionCall(f func(name string, parameters string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorFunctionCall = f
}

func (t *extraTestAdaptor) ValueRead(name string) (val interface{}, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorValueRead(name)
}
func (t *extraTestAdaptor) ValuesRead() (vals map[string]interface{}, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorValuesRead()
}
func (t *extraTestAdaptor) FunctionCall(name string, parameters string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorFunctionCall(name, parameters)
}
func (t *extraTestAdaptor) Connect() (err error)  { return }
func (t *extraTestAdaptor) Finalize() (err error) { return }
func (t *extraTestAdaptor) Name() string          { return t.name }
func (t *extraTestAdaptor) SetName(n string)      { t.name = n }
func (t *extraTestAdaptor) Port() string          { return t.port }

func newExtraTestAdaptor() *extraTestAdaptor {
	return &extraTestAdaptor{
		port: "/dev/null",
		testAdaptorFunctionCall: func(name string, parameters string) (val int, err error) {
			return 0, nil
		},
		testAdaptorValueRead: func(name string) (val interface{}, err error) {
			return 99, nil
		},
		testAdaptorValuesRead: func() (vals map[string]interface{}, err error) {
			return map[string]interface{}{
				"test": 99,
			}, nil
		},
	}
}
