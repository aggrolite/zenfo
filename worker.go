package zenfo

// Worker provides interface for plugin-like crawlers per source
type Worker interface {
	Init(*Client) error
	Desc() string
	Events() ([]*Event, error)
}
