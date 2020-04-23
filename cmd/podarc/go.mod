module podarc/cmd/podarc

go 1.14

require podarc/providers v0.0.0

replace podarc/providers => ../../internal/providers

require podarc/interfaces v0.0.0

replace podarc/interfaces => ../../internal/interfaces

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	podarc/utils v0.0.0
)

replace podarc/utils => ../../internal/utils
