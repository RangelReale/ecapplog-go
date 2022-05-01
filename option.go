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

func WithFlushOnClose(flushOnClose bool) Option {
	return func(c *Client) {
		c.flushOnClose = flushOnClose
	}
}

type logOptions struct {
	source           string
	originalCategory string
	extraCategories  []string
	color            string
	bgColor          string
}

type LogOption func(*logOptions)

func WithSource(source string) LogOption {
	return func(lo *logOptions) {
		lo.source = source
	}
}

func WithOriginalCategory(originalCategory string) LogOption {
	return func(lo *logOptions) {
		lo.originalCategory = originalCategory
	}
}

func WithExtraCategories(extraCategories []string) LogOption {
	return func(lo *logOptions) {
		lo.extraCategories = extraCategories
	}
}

func WithColor(color string) LogOption {
	return func(lo *logOptions) {
		lo.color = color
	}
}

func WithBgColor(bgColor string) LogOption {
	return func(lo *logOptions) {
		lo.bgColor = bgColor
	}
}
