package extra

// ExtraValueRead can read abitrary value
type ExtraValueReader interface {
	ValueRead(name string) (val interface{}, err error)
}
