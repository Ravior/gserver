package ghttp

import (
	"bytes"
	"net/http"
)

// ResponseWriter is the custom writer for ghttp response.
type ResponseWriter struct {
	Status      int                 // HTTP status.
	writer      http.ResponseWriter // The underlying ResponseWriter.
	buffer      *bytes.Buffer       // The output buffer.
	wroteHeader bool                // Is header wrote or not, avoiding errors: superfluous/multiple response.WriteHeader call.
}

// RawWriter returns the underlying ResponseWriter.
func (w *ResponseWriter) RawWriter() http.ResponseWriter {
	return w.writer
}

// Header implements the interface function of ghttp.ResponseWriter.Header.
func (w *ResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write implements the interface function of ghttp.ResponseWriter.Write.
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	return len(data), nil
}

// WriteHeader implements the interface of ghttp.ResponseWriter.WriteHeader.
func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
}

// OutputBuffer outputs the buffer to client and clears the buffer.
func (w *ResponseWriter) Flush() {
	if w.Status != 0 && !w.wroteHeader {
		w.wroteHeader = true
		w.writer.WriteHeader(w.Status)
	}
	// Default status text output.
	if w.Status != http.StatusOK && w.buffer.Len() == 0 {
		w.buffer.WriteString(http.StatusText(w.Status))
	}
	if w.buffer.Len() > 0 {
		w.writer.Write(w.buffer.Bytes())
		w.buffer.Reset()
	}
}
