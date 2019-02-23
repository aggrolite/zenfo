package zenfo

// Worker provides interface for plugin-like crawlers per source
type Worker interface {
	Name() string
	Init(*Client, chan string) error
	Desc() string
	Events() ([]*Event, error)
}
