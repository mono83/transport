// Package multipart provides a multipart/form-data request body writer.
package multipart

import (
	"io"
	mime "mime/multipart"
)

// Part is a single text field in a multipart form.
type Part struct {
	Name  string
	Value string
}

// FilePart is a file upload field in a multipart form.
// Body is always closed by [Write] after its contents are copied.
type FilePart struct {
	Name     string        // form field name
	FileName string        // file name sent to the server
	Body     io.ReadCloser // file contents
}

// Form is the encoded multipart body returned by [Write].
// It implements [io.ReadCloser] and carries the Content-Type header value
// (including the boundary) via [Form.ContentType].
type Form struct {
	io.ReadCloser
	contentType string
}

// ContentType returns the Content-Type header value for this form,
// including the generated multipart boundary. Pass this value as the
// Content-Type request header alongside the body.
func (f *Form) ContentType() string { return f.contentType }

// Write encodes parts and files as a multipart/form-data body.
// Encoding runs in a goroutine; errors are propagated to the reader.
// All FilePart bodies are closed after their contents are copied.
// Use [Form.ContentType] as the Content-Type header value — it includes the
// generated boundary and must not be replaced with a static value.
func Write(parts []Part, files ...FilePart) *Form {
	pr, pw := io.Pipe()
	mw := mime.NewWriter(pw)

	go func() {
		var encErr error

		for _, p := range parts {
			var w io.Writer
			if w, encErr = mw.CreateFormField(p.Name); encErr != nil {
				break
			}
			if _, encErr = io.WriteString(w, p.Value); encErr != nil {
				break
			}
		}

		if encErr == nil {
			for i, f := range files {
				var w io.Writer
				if w, encErr = mw.CreateFormFile(f.Name, f.FileName); encErr != nil {
					for _, remaining := range files[i:] {
						remaining.Body.Close()
					}
					break
				}
				if _, encErr = io.Copy(w, f.Body); encErr != nil {
					f.Body.Close()
					for _, remaining := range files[i+1:] {
						remaining.Body.Close()
					}
					break
				}
				f.Body.Close()
			}
		}

		if encErr != nil {
			mw.Close()
			pw.CloseWithError(encErr)
			return
		}

		pw.CloseWithError(mw.Close())
	}()

	return &Form{
		ReadCloser:  pr,
		contentType: mw.FormDataContentType(),
	}
}
