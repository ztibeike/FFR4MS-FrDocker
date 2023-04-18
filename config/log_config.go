package config

import "github.com/fatih/color"

const (
	LOG_FILE_ROOT_PATH = "/var/log/frecovery/"

	// 颜色配置
	LOG_INFO_COLOR  = color.FgHiBlue
	LOG_ERROR_COLOR = color.FgRed
	LOG_DEBUG_COLOR = color.FgHiGreen
	LOG_FATAL_COLOR = color.FgHiRed
	LOG_TRACE_COLOR = color.FgCyan
	LOG_WARN_COLOR  = color.FgYellow
)
