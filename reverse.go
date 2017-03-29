/*
 Package reverse is a proof of concept package which implements a reverse io.Reader.
*/
package reverse

import "io"

// ReadAtSeeker is minimum interface necessary for creating a Reader
type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

// Reader implements io.Reader and reads bytes in reverse order
type Reader struct {
	// ReadAtSeeker implements io.ReaderAt and io.Seeker (ie: *os.File)
	ras ReadAtSeeker
}

// NewReader is most easily invoked with a *os.File. It returns a new Reader instance
func NewReader(ras ReadAtSeeker) *Reader {
	return &Reader{
		ras: ras,
	}
}

// SeekToEnd is a convenience function to move the offset at the end of the underlying io.Seeker.
// For purposes of reading from some other offset, invoke the underlying io.Seeker. It returns
// the new offset and any error that occurred. The contract is identical to io.Seeker#Seek.
func (r *Reader) SeekToEnd() (int64, error) {
	return r.ras.Seek(0, io.SeekEnd)
}

/*
  Read reads up to len(b) bytes from the underlying ReadAtSeeker, but in reverse order. It returns
  the number of bytes read and any error encountered. At offset 0, Read returns 0, io.EOF.
*/
func (r *Reader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	// This no-op gives us the current offset value
	offset, err := r.ras.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	var m int
	for i := 0; i < len(b); i++ {
		if offset == 0 {
			return m, io.EOF
		}
		// Seek in case someone else is relying on seek too
		offset, err = r.ras.Seek(-1, io.SeekCurrent)
		if err != nil {
			return m, err // Should never happen
		}

		// Just read one byte at a time
		n, err := r.ras.ReadAt(b[i:i+1], offset)
		if err != nil {
			return m + n, err // Should never happen
		}
		m += n
	}
	return m, nil
}
