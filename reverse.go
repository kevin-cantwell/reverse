/*
 Package reverse is a proof of concept package which implements a reverse io.Reader.

 In order to efficiently read a byte stream in reverse order, we must have access to
 a io.Seeker. Almost certainly this will be an instance of *os.File, but this package only
 requires interfaces. The interface declared here is a ReadAtSeeker (io.ReaderAt + io.Seeker).
 The *os.File type implements ReadAt and Seek, which is handy, but not many other types do.

 This package could also have been designed using a more basic combination of just Seek and
 Read. However, ReadAt is useful in that it doesn't change the internal seek position when
 invoked whereas Read does.

 Before reading in reverse, the seek position must be changed from 0, otherwise Read will
 return immediately with an io.EOF error (since we're reading tail-to-head). This can be
 done before instantiating a *Reader by invoking Seek, or after instantiating the reader
 by invoking the convenience method SeekToEnd.

 All calls to Read will read bytes in reverse order, updating the underlying seek offset.
 Bytes can be read in a forwards direction by invoking ReadForward, which is similar to
 a normal Read on the underlying ReadAtSeeker.
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

// SeekToEnd is a convenience function to move the offset to the end of the underlying ReadAtSeeker.
// For purposes of reading from some other offset, invoke the underlying ReadAtSeeker. It returns
// the new offset and any error that occurred. The contract is identical to io.Seeker#Seek(0, io.SeekEnd).
func (r *Reader) SeekToEnd() (int64, error) {
	return r.ras.Seek(0, io.SeekEnd)
}

// SeekToBeginning is a convenience function to move the offset to the beginning of the underlying
// ReadAtSeeker. It returns the new offset and any error that occurred. The contract is identical
// to io.Seeker#Seek(0, io.SeekStart).
func (r *Reader) SeekToStart() (int64, error) {
	return r.ras.Seek(0, io.SeekStart)
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

/*
  ReadForward is the opposite of Read. Bytes are read in the forward direction.
  It returns the number of bytes read and any error encounted. Read and ReadForward
  may be used in combination to reset the underlying seek position to where it started.
*/
func (r *Reader) ReadForward(b []byte) (int, error) {
	offset, err := r.ras.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Make sure to update the seek position, even on errors

	n, err := r.ras.ReadAt(b, offset)
	if n != 0 {
		r.ras.Seek(int64(n), io.SeekCurrent)
	}
	return n, err
}
