module task

go 1.16

replace taskcmd v0.0.0-unpublished => ./cmd
replace taskdb v0.0.0-unpublished => ./db

require (
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	taskcmd v0.0.0-unpublished
	taskdb v0.0.0-unpublished
)
