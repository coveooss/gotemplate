package template

// OptionsSet represents the map of enabled options
type OptionsSet map[Options]bool

// Options defines the type that hold the various options & libraries that should be included
type Options int

// Options values
const (
	Razor Options = iota
	Extension
	Math
	Sprig
	Data
	Logging
	Runtime
	Utils
	Net
	OS
	OptionOnByDefaultCount // Trigger of options that are on by default
	Overwrite
	OutputStdout
	RenderingDisabled
	AcceptNoValue
)

// DefaultOptions returns a OptionsSet with the first options turned on by default
func DefaultOptions() OptionsSet {
	os := make(OptionsSet)
	for i := Options(0); i < OptionOnByDefaultCount; i++ {
		os[i] = true
	}
	return os
}

//go:generate stringer -type=Options -output generated_options.go
