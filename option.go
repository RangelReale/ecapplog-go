package ecapplog

type Option func(*Client)

func WithAddress(address string) Option {
	return func(c *Client) {
		c.address = address
	}
}

func WithAppName(appname string) Option {
	return func(c *Client) {
		c.appname = appname
	}
}

func WithBufferSize(bufferSize int) Option {
	return func(c *Client) {
		c.bufferSize = bufferSize
	}
}

type logOptions struct {
	source          string
	extraCategories []string
}

type LogOption func(*logOptions)

func WithSource(source string) LogOption {
	return func(lo *logOptions) {
		lo.source = source
	}
}

func WithExtraCategories(extraCategories []string) LogOption {
	return func(lo *logOptions) {
		lo.extraCategories = extraCategories
	}
}
