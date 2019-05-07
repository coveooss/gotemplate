package flag

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestNewApplication(t *testing.T) {
	tests := []struct {
		name string
		args []string
		b    bool
		s    string
		err  error
	}{
		{"No Args", nil, true, "-", nil},
		{"With long flag", []string{"--bool-flag"}, true, "-", nil},
		{"With negative long flag", []string{"--no-bool-flag"}, false, "-", nil},
		{"With alias", []string{"--bool"}, true, "-", nil},
		{"With negative alias", []string{"--no-bool"}, false, "-", nil},
		{"With other alias", []string{"--switch"}, true, "-", nil},
		{"With shortcut", []string{"--bf"}, true, "-", nil},
		{"With negative shortcut", []string{"--nb"}, false, "-", nil},
		{"Duplicated flags", []string{"--bf", "--bool"}, true, "-", fmt.Errorf("flag 'bool-flag' cannot be repeated")},
		{"With string long arg 1", []string{"--sf=test"}, true, "test", nil},
		{"With string long arg 2", []string{"--sf", "test"}, true, "test", nil},
		{"With missing long arg", []string{"--sf"}, true, "-", fmt.Errorf("expected argument for flag '--string-flag'")},
		{"With short arg 1", []string{"-s", "test"}, true, "test", nil},
		{"With short arg 2", []string{"-stest"}, true, "test", nil},
		{"With short arg 3", []string{"-s=test"}, true, "=test", nil},
		{"With missing short", []string{"-s"}, true, "-", fmt.Errorf("expected argument for flag '-s'")},
		{"With extra arg 1", []string{"--b", "-sx", "--unmanaged", "extra"}, true, "x", fmt.Errorf("unknown long flag '--unmanaged'")},
		{"With extra arg 2", []string{"--switch", "-sx", "extra"}, true, "x", fmt.Errorf("unexpected extra")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication(tt.name).
				Author("coveo").
				AutoShortcut().
				DefaultEnvars().
				Writer(os.Stdout).
				ErrorWriter(os.Stderr).
				UsageWriter(os.Stderr).
				UsageTemplate("").
				Terminate(func(int) {}).
				Validate(func(*kingpin.Application) error { return nil })
			b := app.Flag("bool-flag", "").Default(true).Alias("bool", "switch").Bool()
			s := app.Flag("string-flag", "").Short('s').Default("-").String()
			_, err := app.Parse(tt.args)
			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.b, *b)
				assert.Equal(t, tt.s, *s)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

func TestUnmanaged(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		b     bool
		s     string
		extra []string
		err   error
	}{
		{"With extra arg 1", []string{"--unmanaged"}, false, "-", []string{"--unmanaged"}, nil},
		{"With extra arg 2", []string{"-b", "-sx", "extra"}, true, "x", []string{"extra"}, nil},
		{"With extra arg 3", []string{"-xabc", "-b", "extra", "-sx", "--unmanaged"}, true, "x", []string{"-xabc", "extra", "--unmanaged"}, nil},
		{"With all extra", []string{"--", "-b", "test", "--sf", "--unmanaged"}, false, "-", []string{"-b", "test", "--sf", "--unmanaged"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication(tt.name).AllowUnmanaged()
			b := app.Flag("bool-flag", "").Short('b').Alias("bool", "switch").Bool()
			s := app.Flag("string-flag", "").Short('s').Default("-").String()
			_, err := app.Parse(tt.args)
			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.b, *b)
				assert.Equal(t, tt.s, *s)
				assert.Equal(t, tt.extra, app.UnmanagedArgs)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

func TestOthers(t *testing.T) {
	app := NewApplication("My app").AutoShortcut()
	var list []*Clause
	list = append(list, app.Flag("bool-flag", ""))
	list = append(list, app.Flag("string-flag", ""))
	list[0].Bool()
	list[1].String()
	assert.ElementsMatch(t, app.Flags(), list)
	app.Flag("other", "").Alias("bf").String()
	_, err := app.Parse(nil)
	assert.EqualError(t, err, "Flag alias bf on other is already mapped to bool-flag")
	assert.Nil(t, app.getFlag("none"))
	assert.NotNil(t, app.getFlag("bf"))
	assert.NotNil(t, app.getFlag("nbf"))
	assert.NotNil(t, app.getFlag("no-bool-flag"))
	assert.Nil(t, app.getFlag("no-arg"))
	assert.True(t, app.getFlag("help").IsSwitch())
	assert.Equal(t, app.getFlag("help").Name(), "help")

	app.Flag("all", "").
		Action(func(*kingpin.ParseContext) error { return nil }).
		PreAction(func(*kingpin.ParseContext) error { return nil }).
		HintAction(func() []string { return nil }).
		HintOptions().
		Envar("TEST").
		NoEnvar().
		PlaceHolder("").
		Hidden().
		Required().
		Short('A').
		NoAutoShortcut()
}
