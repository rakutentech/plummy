package jvm

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
)

var versionLineRegex = regexp.MustCompile("^java version \"([^\"])+\"")
var oldStyleVersionRegex = regexp.MustCompile("^1\\.(\\d+)")
var newStyleVersionRegex = regexp.MustCompile("^\\d+")

type Version struct {
	Major int
}

func makeVersion(majorStr string) (Version, error) {
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major}, nil
}

func parseVersion(s string) (Version, error) {
	matches := oldStyleVersionRegex.FindStringSubmatch(s)
	if len(matches) >= 2 {
		return makeVersion(matches[1])
	}
	major := newStyleVersionRegex.FindString(s)
	if major == "" {
		return Version{}, fmt.Errorf("bad java version format '%s'", s)
	}
	return makeVersion(major)
}

func ExecVersion(path string) (Version, error) {
	stderr, err := captureStderr(path, "-version")
	if err != nil {
		return Version{}, err
	}
	line, err := bufio.NewReader(stderr).ReadString(byte('\n'))
	if err != nil {
		return Version{}, err
	}
	matches := versionLineRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return Version{}, fmt.Errorf("bad java version line format '%s'", line)
	}
	return parseVersion(matches[1])
}

