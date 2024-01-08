package headers

type ContentEncoding string

// content encodings sourced from MDN
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
const (
	GzipCompression ContentEncoding = "gzip"
	Compress        ContentEncoding = "compress"
	Deflate         ContentEncoding = "deflate"
	Brotli          ContentEncoding = "br"
)
