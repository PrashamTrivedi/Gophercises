module taskcmd

go 1.16

replace taskdb v0.0.0-unpublished => ../db

require (
	github.com/spf13/cobra v1.1.3
	taskdb v0.0.0-unpublished
)
