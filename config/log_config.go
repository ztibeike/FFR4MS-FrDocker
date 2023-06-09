package config

import "github.com/fatih/color"

const (
	LOG_FILE_ROOT_PATH = "/var/log/frecovery/"

	LOG_FILE = "frecovery.log"

	LOG_BANNER = `
	.----------------.  .----------------.  .----------------.  .----------------.  .----------------.  .----------------.  .----------------.  .----------------. 
	| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |
	| |  _________   | || |  _______     | || |  ________    | || |     ____     | || |     ______   | || |  ___  ____   | || |  _________   | || |  _______     | |
	| | |_   ___  |  | || | |_   __ \    | || | |_   ___  \  | || |   /      \   | || |   .' ___  |  | || | |_  ||_  _|  | || | |_   ___  |  | || | |_   __ \    | |
	| |   | |_  \_|  | || |   | |__) |   | || |   | |    \ \ | || |  /  .--.  \  | || |  / .'   \_|  | || |   | |_/ /    | || |   | |_  \_|  | || |   | |__) |   | |
	| |   |  _|      | || |   |  __ /    | || |   | |    | | | || |  | |    | |  | || |  | |         | || |   |  __'.    | || |   |  _|  _   | || |   |  __ /    | |
	| |  _| |_       | || |  _| |  \ \_  | || |  _| |___.' / | || |  \  '--'  /  | || |  \ \.___.'\  | || |  _| |  \ \_  | || |  _| |___/ |  | || |  _| |  \ \_  | |
	| | |_____|      | || | |____| |___| | || | |________.'  | || |   '.____.'   | || |   '._____.'  | || | |____||____| | || | |_________|  | || | |____| |___| | |
	| |              | || |              | || |              | || |              | || |              | || |              | || |              | || |              | |
	| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |
	 '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------' 
	`

	LOG_CALLER_ENABLED = false

	// 颜色配置
	LOG_INFO_COLOR  = color.FgCyan
	LOG_ERROR_COLOR = color.FgRed
	LOG_DEBUG_COLOR = color.FgWhite
	LOG_TRACE_COLOR = color.FgWhite
	LOG_WARN_COLOR  = color.FgYellow
	LOG_FATAL_COLOR = color.FgRed
)
