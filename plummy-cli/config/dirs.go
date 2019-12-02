package config

import (
	"log"
	"os"
	"path"
	"runtime"
)

var cacheDir = func() string {
	if envDir, ok := os.LookupEnv("PLUMMY_CACHE_DIR"); ok {
		return envDir
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Panicf("cannot get user cache dir: %v", err)
	}

	var appDirName string
	switch runtime.GOOS {
	case "darwin":
		appDirName = "com.github.rakutentech.plummy"
	default:
		appDirName = "plummy"
	}

	return path.Join(cacheDir, appDirName)
}()

func CacheDir() string {
	return cacheDir
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
