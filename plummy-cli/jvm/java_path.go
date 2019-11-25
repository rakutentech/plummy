package jvm

import (
	"os"
	"path"
	"runtime"
)

func isFile(dir string) bool {
	fi, _ := os.Stat(dir)
	return fi != nil && !fi.IsDir()
}

func javaExecFilename() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}

func getJavaExec(dir string) string {
	execFilename := javaExecFilename()
	javaBin := path.Join(dir, "bin", execFilename)
	if isFile(javaBin) {
		return javaBin
	}
	// Extra check for macOS packages
	if runtime.GOOS == "darwin" {
		javaBin = path.Join(dir, "Contents/Home/bin/java")
		if isFile(javaBin) {
			return javaBin
		}
	}
	return ""
}