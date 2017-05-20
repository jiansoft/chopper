package system

type IDisposable interface {
	Dispose()
	IdentifyId() string
}
