package multipart

import (
	"io"
	mime "mime/multipart"
	"net/textproto"
	"strings"
	"testing"
)

func TestWrite_formFields(t *testing.T) {
	form := Write([]Part{
		{Name: "username", Value: "alice"},
		{Name: "message", Value: "hello world"},
	})
	defer form.Close()

	if form.ContentType() == "" {
		t.Fatal("ContentType is empty")
	}
	if !strings.HasPrefix(form.ContentType(), "multipart/form-data; boundary=") {
		t.Fatalf("unexpected Content-Type: %q", form.ContentType())
	}

	body, err := io.ReadAll(form)
	if err != nil {
		t.Fatalf("unexpected error reading form: %v", err)
	}

	mediaType, params := parseContentType(t, form.ContentType())
	if mediaType != "multipart/form-data" {
		t.Fatalf("got media type %q, want %q", mediaType, "multipart/form-data")
	}

	mr := mime.NewReader(strings.NewReader(string(body)), params["boundary"])
	fields := map[string]string{}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("error reading part: %v", err)
		}
		val, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("error reading part body: %v", err)
		}
		fields[part.FormName()] = string(val)
	}

	if fields["username"] != "alice" {
		t.Errorf("username: got %q, want %q", fields["username"], "alice")
	}
	if fields["message"] != "hello world" {
		t.Errorf("message: got %q, want %q", fields["message"], "hello world")
	}
}

func TestWrite_emptyParts(t *testing.T) {
	form := Write(nil)
	defer form.Close()

	if form.ContentType() == "" {
		t.Fatal("ContentType is empty")
	}

	body, err := io.ReadAll(form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, params := parseContentType(t, form.ContentType())
	mr := mime.NewReader(strings.NewReader(string(body)), params["boundary"])
	_, err = mr.NextPart()
	if err != io.EOF {
		t.Fatalf("expected EOF for empty form, got %v", err)
	}
}

func TestWrite_singleField(t *testing.T) {
	form := Write([]Part{{Name: "key", Value: "value"}})
	defer form.Close()

	body, _ := io.ReadAll(form)
	_, params := parseContentType(t, form.ContentType())
	mr := mime.NewReader(strings.NewReader(string(body)), params["boundary"])

	part, err := mr.NextPart()
	if err != nil {
		t.Fatalf("error reading part: %v", err)
	}
	val, _ := io.ReadAll(part)
	if string(val) != "value" {
		t.Errorf("got %q, want %q", val, "value")
	}
	if _, err = mr.NextPart(); err != io.EOF {
		t.Fatalf("expected EOF after single part, got %v", err)
	}
}

// parseContentType splits a Content-Type value into media type and params.
func parseContentType(t *testing.T, ct string) (string, map[string]string) {
	t.Helper()
	// textproto header parsing: mime package expects header-style input
	h := textproto.MIMEHeader{"Content-Type": {ct}}
	_ = h
	// parse manually: "multipart/form-data; boundary=xxx"
	parts := strings.SplitN(ct, ";", 2)
	mediaType := strings.TrimSpace(parts[0])
	params := map[string]string{}
	if len(parts) == 2 {
		for _, kv := range strings.Split(parts[1], ";") {
			kv = strings.TrimSpace(kv)
			if k, v, ok := strings.Cut(kv, "="); ok {
				params[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), `"`)
			}
		}
	}
	return mediaType, params
}
