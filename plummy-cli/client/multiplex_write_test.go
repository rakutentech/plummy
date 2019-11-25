package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
	"unsafe"
)

func readRand(n int) []byte {
	b := make([]byte, n)
	count, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if count != n {
		panic(fmt.Sprintf("Read %d instead of %d", count, n))
	}
	return b
}

func makeEncoded(parts ...interface{}) string {
	l := 0
	for _, part := range parts {
		switch v := part.(type) {
		case string:
			l += len(v)
		case []byte:
			l += len(v)
		}
	}
	dest := make([]byte, 0, l)
	for _, part := range parts {
		switch v := part.(type) {
		case string:
			dest = append(dest, v...)
		case []byte:
			dest = append(dest, v...)
		}
	}
	return bytesToStr(dest)
}

func bytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func TestMultiplex(t *testing.T) {
	type test struct {
		Name  string
		Files []FileInput
		Want  string
	}

	small1 := FileInput{
		Name:     "My Name",
		Metadata: nil,
		Reader:   strings.NewReader("Hello World\nFoo\n123\n"),
	}
	small1Encoded := "\x07\x00\x00\x00My Name\x00\x00\x00\x00\x14\x00\x00\x00Hello World\nFoo\n123\n"
	small2 := FileInput{
		Name:     "Long long long long long filename",
		Metadata: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Reader:   strings.NewReader("Hello World\nFoo\n123\n"),
	}
	small2Encoded := "\x21\x00\x00\x00Long long long long long filename\x06\x00\x00\x00\xaa\xbb\xcc\xdd\xee\xff\x14\x00\x00\x00Hello World\nFoo\n123\n"
	small3 := FileInput{
		Name:     "Small3",
		Metadata: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Reader:   strings.NewReader("Hello World\nFoo\n123\n"),
	}
	small3Encoded := "\x06\x00\x00\x00Small3\x06\x00\x00\x00\xaa\xbb\xcc\xdd\xee\xff\x14\x00\x00\x00Hello World\nFoo\n123\n"

	large1Data := readRand(128 * 1024)
	large2Data := readRand(16 * 1024 * 1024)

	large1 := FileInput{
		Name:     "Large1",
		Metadata: nil,
		Reader:   bytes.NewReader(large1Data),
	}
	large2 := FileInput{
		Name:     "Large2",
		Metadata: nil,
		Reader:   bytes.NewReader(large2Data),
	}

	tests := []test{
		{
			Name:  "No files",
			Files: []FileInput{},
			Want:  "",
		},
		{
			Name:  "One file, no metadata",
			Files: []FileInput{small1},
			Want:  small1Encoded,
		},
		{
			Name:  "One file, with metadata",
			Files: []FileInput{small2},
			Want:  small2Encoded,
		},
		{
			Name: "Two small files",
			Files: []FileInput{
				small1,
				small2,
			},
			Want: small1Encoded + small2Encoded,
		},
		{
			Name: "Small1 Large1 Small2 files",
			Files: []FileInput{
				small1,
				large1,
				small2,
			},
			Want: makeEncoded(
				small1Encoded,
				"\x06\x00\x00\x00Large1", "\x00\x00\x00\x00", "\x00\x00\x02\x00", large1Data,
				small2Encoded,
			),
		},
		{
			Name: "Small1 Large1 Small2 Large1 Small2 files",
			Files: []FileInput{
				small1,
				large1,
				small2,
				large2,
				small3,
			},
			Want: makeEncoded(
				small1Encoded,
				"\x06\x00\x00\x00Large1", "\x00\x00\x00\x00", "\x00\x00\x02\x00", large1Data,
				small2Encoded,
				"\x06\x00\x00\x00Large2", "\x00\x00\x00\x00", "\x00\x00\x00\x01", large2Data,
				small3Encoded,
			),
		},
	}

	for _, tc := range tests {
		for _, f := range tc.Files {
			// Reset reader to start for each file
			_, err := f.Reader.(io.Seeker).Seek(0, io.SeekStart)
			if err != nil {
				t.Errorf("Case [%s] Got error on reader.Seek(): %v", tc.Name, err)
			}
		}
		want := []byte(tc.Want)
		reader := newMultiplexEncoderReader(tc.Files)
		got, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("Case [%s] Got error: %v", tc.Name, err)
		}
		if !bytes.Equal(want, got) {
			t.Errorf("Case [%s] failed:\nexpected:\n%s\ngot:\n%s", tc.Name, hexDump(want), hexDump(got))
		}
	}
}

func hexDump(b []byte) string {
	if len(b) > 1024 {
		b = b[:1024]
	}
	return hex.Dump(b)
}