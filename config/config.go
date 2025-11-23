package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed version
var version string

//go:embed name
var name string

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"

	ModeMaster = "master"
	ModeClient = "client"
)

func GetVersion() string {
	return strings.TrimSpace(version)
}

func GetName() string {
	return strings.TrimSpace(name)
}

func GetLogLevel() LogLevel {
	if IsDebug() {
		return Debug
	}
	logLevel := os.Getenv("SUI_LOG_LEVEL")
	if logLevel == "" {
		return Info
	}
	return LogLevel(logLevel)
}

func IsDebug() bool {
	return os.Getenv("SUI_DEBUG") == "true"
}

func GetDBFolderPath() string {
	dbFolderPath := os.Getenv("SUI_DB_FOLDER")
	if dbFolderPath == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			// Cross-platform fallback path
			if runtime.GOOS == "windows" {
				return "C:\\Program Files\\s-ui\\db"
			}
			return "/usr/local/s-ui/db"
		}
		dbFolderPath = filepath.Join(dir, "db")
	}
	return dbFolderPath
}

func GetDBPath() string {
	return fmt.Sprintf("%s/%s.db", GetDBFolderPath(), GetName())
}

func GetClusterMode() string {
	mode := os.Getenv("SUI_MODE")
	if mode == "" {
		return ModeMaster
	}
	return mode
}

func GetClusterToken() string {
	return os.Getenv("SUI_CLUSTER_TOKEN")
}

func GetMasterEndpoint() string {
	return os.Getenv("SUI_MASTER_ENDPOINT")
}

func GetNodeID() string {
	if id := os.Getenv("SUI_NODE_ID"); id != "" {
		return id
	}

	host, err := os.Hostname()
	if err != nil {
		return ""
	}
	return host
}

func GetNodeAddress() string {
	return os.Getenv("SUI_NODE_ADDRESS")
}
