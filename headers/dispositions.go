package headers

import "fmt"

type ContentDisposition string

// common content dispositions sourced from MDN
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition#syntax
const (
	Inline     ContentDisposition = "inline"
	Attachment ContentDisposition = "attachment"
)

// generates a disposition in the format `attachment; filename="filename"`
//
// the filename directive tells browsers the file should be downloaded with the specified name.
// more info: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition#syntax
func FileAttachment(filename string) ContentDisposition {
	return ContentDisposition(fmt.Sprintf("%s; filename=\"%s\"", Attachment, filename))
}
