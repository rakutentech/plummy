package client

import (
	"encoding/binary"
	"fmt"
	"io"
)

func readFiles(r io.Reader) ([]FileData, error) {
	var files []FileData
	for {
		file, err := readFile(r)
		if err == io.EOF {
			// This is not an error: we just have no more files
			break
		}
		if err != nil {
			return nil, err
		}
		// Add file
		files = append(files, file)
	}
	return files, nil
}

func readFile(r io.Reader) (FileData, error) {
	name, err := readChunk(r)
	if err != nil {
		return FileData{}, err
	}
	metadata, err := readChunk(r)
	if err != nil {
		return FileData{}, err
	}
	data, err := readChunk(r)
	if err != nil {
		return FileData{}, err
	}

	return FileData{
		Name:     string(name),
		Metadata: metadata,
		Data:     data,
	}, nil
}

func readChunk(r io.Reader) ([]byte, error) {
	chunkSize, err := readInt(r)
	if err != nil {
		return nil, err
	}
	chunk := make([]byte, chunkSize)
	bytesRead, err := io.ReadFull(r, chunk)
	if err != nil {
		return nil, err
	}
	if bytesRead != chunkSize {
		return nil, fmt.Errorf("read size mismatch (expected: %d, got: %d)", chunkSize, bytesRead)
	}
	return chunk, nil
}

func readInt(r io.Reader) (int, error) {
	var value int32
	if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	if value < 0 {
		value = 0
	}
	return int(value), nil
}
