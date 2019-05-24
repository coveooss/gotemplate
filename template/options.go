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
	StrictErrorCheck
)

// Set options to true
func (os OptionsSet) Set(options ...Options) OptionsSet { return os.set(true, options) }

// Unset options
func (os OptionsSet) Unset(options ...Options) OptionsSet { return os.set(false, options) }

func (os OptionsSet) set(value bool, options []Options) OptionsSet {
	for i := range options {
		os[options[i]] = value
	}
	return os
}

// DefaultOptions returns a OptionsSet with the first options turned on by default
func DefaultOptions() OptionsSet {
	os := make(OptionsSet)
	for i := Options(0); i < OptionOnByDefaultCount; i++ {
		os[i] = true
	}
	return os
}

//go:generate stringer -type=Options -output generated_options.go
