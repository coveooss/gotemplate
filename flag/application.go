package flag

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/coveo/gotemplate/v3/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Application allows proper management between managed and non managed arguments provided to kingpin
type Application struct {
	*kingpin.Application
	UnmanagedArgs []string

	allowUnmanaged bool
	autoShortcut   bool
	flags          map[string]*Clause
	runes          map[rune]*Clause
}

type app = Application
type av = kingpin.ApplicationValidator

func (a *app) Author(author string) *app          { a.Application.Author(author); return a }
func (a *app) DefaultEnvars() *app                { a.Application.DefaultEnvars(); return a }
func (a *app) Terminate(terminate func(int)) *app { a.Application.Terminate(terminate); return a }
func (a *app) Writer(w io.Writer) *app            { a.Application.Writer(w); return a }
func (a *app) ErrorWriter(w io.Writer) *app       { a.Application.ErrorWriter(w); return a }
func (a *app) UsageWriter(w io.Writer) *app       { a.Application.UsageWriter(w); return a }
func (a *app) UsageTemplate(template string) *app { a.Application.UsageTemplate(template); return a }
func (a *app) Validate(validator av) *app         { a.Application.Validate(validator); return a }

// AllowUnmanaged allows arguments to contain non managed args
func (app *Application) AllowUnmanaged() *Application {
	app.allowUnmanaged = true
	return app
}

// AutoShortcut automatically generates short cut from arguments by using the first letter of each argument
func (app *Application) AutoShortcut() *Application {
	app.autoShortcut = true
	return app
}

// Flag overwrite the base Flag definition to force providing a value for the given flag
func (app *Application) Flag(name string, description string) *Clause {
	newFlag := Clause{FlagClause: app.Application.Flag(name, description)}
	app.flags[name] = &newFlag
	return &newFlag
}

// Flags returns the list of all flags in the application
func (app *Application) Flags() (result []*Clause) {
	for _, clause := range app.flags {
		result = append(result, clause)
	}
	return
}

// Parse splits the managed by kingpin and unmanaged argument to avoid error
func (app *Application) Parse(args []string) (command string, err error) {
	app.runes = make(map[rune]*Clause)
	var errs errors.Array
	for _, f := range app.Model().Flags {
		flag := app.flags[f.Name]
		if flag == nil {
			// The flag exists, but have been added directly to kingpin
			flag = &Clause{extraArg: f}
			app.flags[f.Name] = flag
		}
		if f.Short != 0 {
			app.runes[f.Short] = flag
		}

		addShortcut := func(name string, t aliasType) {
			if current := app.flags[name]; current == nil {
				app.flags[name] = flag
				flag.addShortcut(name, t)
			} else if current != flag {
				errs = append(errs, errors.Managed(fmt.Sprintf("Flag alias %s on %s is already mapped to %s", name, flag.Name(), current.Name())))
			}
		}

		tryAddShortcut := func(name string) {
			if app.autoShortcut && flag.FlagClause != nil && !flag.noAutoShortcut {
				var shortcut string
				for _, word := range strings.Split(name, "-") {
					shortcut += string(word[0])
				}
				if flag.Model().Short == 0 || len(shortcut) > 1 {
					addShortcut(shortcut, regularShortcut)
				}
				if flag.IsSwitch() {
					addShortcut("n"+shortcut, negativeShortcut)
				}
			}
		}
		tryAddShortcut(f.Name)

		for alias, value := range flag.aliases {
			if value <= realAlias {
				addShortcut(alias, realAlias)
				tryAddShortcut(alias)
			}
		}
	}

	if len(errs) > 0 {
		return "", errs
	}

	managed, unmanaged := app.splitManaged(args)
	command, err = app.Application.Parse(managed)
	app.UnmanagedArgs = unmanaged

	// We reset the switch default values after evaluation because they should not be revaluated on subsequent parse
	for _, flag := range app.flags {
		if flag.FlagClause != nil && flag.IsSwitch() && len(flag.Model().Default) != 0 {
			flag.Default()
		}
	}
	return
}

func (app *Application) getFlag(argName string) *Clause {
	if flag := app.flags[argName]; flag != nil {
		return flag
	}
	if !strings.HasPrefix(argName, "no-") {
		return nil
	}
	if flag := app.flags[argName[3:]]; flag != nil && flag.IsSwitch() {
		return flag
	}
	return nil
}

func (app *Application) splitManaged(args []string) (managed, unmanaged []string) {
	notManaged := func(arg ...string) {
		if app.allowUnmanaged {
			unmanaged = append(unmanaged, arg...)
		} else {
			managed = append(managed, arg...)
		}
	}
Arg:
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			// All arguments after -- are considered as unmanaged
			notManaged(args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			argSplit := strings.SplitN(args[i][2:], "=", 2)
			argName := argSplit[0]
			if flag := app.getFlag(argName); flag != nil {
				managed = append(managed, flag.substituteAlias(argName, arg))
				if !flag.IsSwitch() && len(argSplit) == 1 {
					// This is not a switch (bool flag) and there is no argument with
					// the flag, so the argument must be after and we add it to
					// the managed args if there is.
					i++
					if i < len(args) {
						managed = append(managed, args[i])
					}
				}
			} else {
				notManaged(arg)
			}
		} else if strings.HasPrefix(arg, "-") {
			withArg := false
			for pos, opt := range arg[1:] {
				if flag := app.runes[opt]; flag != nil {
					if !flag.IsSwitch() {
						// This is not a switch (bool flag), so we check if there are characters
						// following the current flag in the same word. If it is not the case,
						// then the argument must be after and we add it to the managed args
						// if there is. If it is the case, then, the argument is included in
						// the current flag and we consider the whole word as a managed argument.
						withArg = pos == len(arg[1:])-1
						break
					}
				} else {
					notManaged(arg)
					continue Arg
				}
			}
			managed = append(managed, arg)
			if withArg {
				// The next argument must be an argument to the current flag
				i++
				if i < len(args) {
					managed = append(managed, args[i])
				}
			}
		} else {
			notManaged(arg)
		}
	}
	return
}

// NewApplication returns an initialized copy of TGFApplication
func NewApplication(description string) *Application {
	return &Application{
		Application: kingpin.New(os.Args[0], description),
		flags:       make(map[string]*Clause),
	}
}
