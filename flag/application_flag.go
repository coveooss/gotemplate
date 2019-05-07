package flag

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type aliasType byte

const (
	undefinedAlias aliasType = iota
	realAlias
	regularShortcut
	negativeShortcut
)

// Clause extends the functionalities of base object from kingpin
type Clause struct {
	*kingpin.FlagClause
	extraArg       *kingpin.FlagModel // Flag that are part of the application but have been added directly through kingpin
	aliases        map[string]aliasType
	noAutoShortcut bool
}

type action = kingpin.Action
type hintAction = kingpin.HintAction
type fc = Clause

func (f *fc) Action(action action) *fc           { f.FlagClause.Action(action); return f }
func (f *fc) PreAction(action action) *fc        { f.FlagClause.PreAction(action); return f }
func (f *fc) HintAction(action hintAction) *fc   { f.FlagClause.HintAction(action); return f }
func (f *fc) HintOptions(options ...string) *fc  { f.FlagClause.HintOptions(options...); return f }
func (f *fc) Envar(name string) *fc              { f.FlagClause.Envar(name); return f }
func (f *fc) NoEnvar() *fc                       { f.FlagClause.NoEnvar(); return f }
func (f *fc) PlaceHolder(placeholder string) *fc { f.FlagClause.PlaceHolder(placeholder); return f }
func (f *fc) Hidden() *fc                        { f.FlagClause.Hidden(); return f }
func (f *fc) Required() *fc                      { f.FlagClause.Required(); return f }
func (f *fc) Short(name rune) *fc                { f.FlagClause.Short(name); return f }

// NoAutoShortcut disables the auto shortcut generation for this field if app AutoShortcut is enabled
func (f *Clause) NoAutoShortcut() *Clause { f.noAutoShortcut = true; return f }

// Default allows providing a default by its real type instead of by string
func (f *Clause) Default(values ...interface{}) *Clause {
	init := make([]string, len(values))
	for i := range values {
		init[i] = fmt.Sprint(values[i])
	}
	f.FlagClause.Default(init...)
	return f
}

// IsSwitch determines if the current flag object is a switch
func (f *Clause) IsSwitch() bool {
	if f.extraArg != nil {
		return f.extraArg.IsBoolFlag()
	}
	return f.Model().IsBoolFlag()
}

// Name returns the name of the flag
func (f *Clause) Name() string {
	if f.extraArg != nil {
		return f.extraArg.Name
	}
	return f.Model().Name
}

// Alias adds a collection of alias to the current object
func (f *Clause) Alias(aliases ...string) *Clause {
	for _, alias := range aliases {
		f.addShortcut(alias, realAlias)
	}
	return f
}

func (f *Clause) addShortcut(name string, t aliasType) {
	if f.aliases == nil {
		f.aliases = make(map[string]aliasType)
	}
	f.aliases[name] = t
}

func (f *Clause) substituteAlias(name, arg string) (result string) {
	if f.Name() == name {
		return arg
	}

	if f.IsSwitch() {
		if strings.HasPrefix(name, "no-") {
			name = name[3:]
		}
		if f.aliases[name] == negativeShortcut {
			return strings.Replace(arg, name, "no-"+f.Name(), 1)
		}
	}
	return strings.Replace(arg, name, f.Name(), 1)
}
