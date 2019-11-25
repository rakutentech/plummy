package client

import (
	"encoding/binary"
	"errors"
	"io"
)

const multiplexContentType = "application/x-plummy-multiplex"

func newMultiplexEncoderReader(files []FileInput) *multiplexEncoderReader {
	return &multiplexEncoderReader{files: files}
}

// multiplexEncoderReader reads multiple FileInputs as if they were encoded as one multiplexed file
type multiplexEncoderReader struct {
	files       []FileInput
	current     int
	wroteHeader bool
}

func (m *multiplexEncoderReader) Read(p []byte) (int, error) {
	readSize := len(p)
	if m.isEOF(p) {
		return eofResult(0, readSize)
	}

	// Get current file
	file := m.files[m.current]

	// Write header for first file if needed
	bytesRead := 0
	if !m.wroteHeader {
		headerWriter := newByteWriter(p)
		err := writeHeader(headerWriter, file)
		if err != nil {
			return headerWriter.BytesWritten, err
		}
		m.wroteHeader = true

		// Adjust bytesRead and p to accommodate remaining space
		bytesRead = headerWriter.BytesWritten
		p = p[headerWriter.BytesWritten:]
	}

	if len(p) == 0 {
		// No remaining read buffer
		return bytesRead, nil
	}

	// Read/write as much content as possible
	n, err := file.Reader.Read(p)
	bytesRead += n

	// Advance to next file if we ran out of space here
	if err == io.EOF {
		// Advance current file and reset header flag
		// (since we need to write header again for the next file
		m.current += 1
		m.wroteHeader = false

		// Return EOF if we reached last file
		if m.isEOF(p) {
			return eofResult(bytesRead, readSize)
		}
	}
	return bytesRead, nil
}

func (m *multiplexEncoderReader) isEOF(p []byte) bool {
	return m.current == len(m.files)
}

func eofResult(bytesRead, readSize int) (int, error) {
	// If read size is 0, we shouldn't return an error
	if readSize == 0 {
		return bytesRead, nil
	}
	return bytesRead, io.EOF
}

func writeHeader(w io.Writer, input FileInput) error {
	err := writeChunk(w, []byte(input.Name))
	if err != nil {
		return err
	}
	err = writeChunk(w, input.Metadata)
	if err != nil {
		return err
	}
	err = writeInt(w, int(input.Reader.Size()))
	if err != nil {
		return err
	}

	return nil
}

func writeChunk(w io.Writer, data []byte) error {
	err := writeInt(w, len(data))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func writeInt(w io.Writer, value int) error {
	if value < 0 {
		value = 0
	}
	return binary.Write(w, binary.LittleEndian, int32(value))
}

type byteWriter struct {
	output       []byte
	BytesWritten int
}

func newByteWriter(output []byte) *byteWriter {
	return &byteWriter{output: output, BytesWritten: 0}
}

func (bw *byteWriter) Write(p []byte) (int, error) {
	if len(p) > len(bw.output) {
		return 0, errors.New("too many bytes to write")
	}
	written := copy(bw.output, p)
	bw.output = bw.output[written:]
	bw.BytesWritten += written
	return written, nil
}
