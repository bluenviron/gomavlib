package frame

// ReadWriter is a Frame Reader and Writer.
type ReadWriter struct {
	*Reader
	*Writer
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(readerConf ReaderConf, writerConf WriterConf) (*ReadWriter, error) {
	r, err := NewReader(readerConf)
	if err != nil {
		return nil, err
	}

	w, err := NewWriter(writerConf)
	if err != nil {
		return nil, err
	}

	return &ReadWriter{
		Reader: r,
		Writer: w,
	}, nil
}
