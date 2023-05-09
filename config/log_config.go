package config

import "github.com/fatih/color"

const (
	LOG_FILE_ROOT_PATH = "/var/log/frecovery/"

	LOG_CALLER_ENABLED = true

	// 颜色配置
	LOG_INFO_COLOR  = color.FgCyan
	LOG_ERROR_COLOR = color.FgRed
	LOG_DEBUG_COLOR = color.FgWhite
	LOG_TRACE_COLOR = color.FgWhite
	LOG_WARN_COLOR  = color.FgYellow
	LOG_FATAL_COLOR = color.FgRed
)
