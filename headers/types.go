package headers

type ContentType string

// common content types sourced from MDN's list of common MIME types
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
const (
	Text   ContentType = "text/plain"
	HTML   ContentType = "text/html"
	CSS    ContentType = "text/css"
	JS     ContentType = "text/javascript"
	JSON   ContentType = "application/json"
	CSV    ContentType = "text/csv"
	Binary ContentType = "application/octet-stream"
	Gzip   ContentType = "application/gzip"
	GIF    ContentType = "image/gif"
	JPEG   ContentType = "image/jpeg"
	PNG    ContentType = "image/png"
	Icon   ContentType = "image/vnd.microsoft.icon"
	WOFF   ContentType = "font/woff"
	WOFF2  ContentType = "font/woff2"
	OTF    ContentType = "font/otf"
	TTF    ContentType = "font/ttf"
)

// create a new content type
//
// useful if your required content type is not defined in this package
func New(value string) ContentType {
	return ContentType(value)
}
