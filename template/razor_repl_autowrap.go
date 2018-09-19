package template

import (
	"fmt"

	"github.com/coveo/gotemplate/utils"
)

func autoWrap(repl replacement, match string) string {
	matches, _ := utils.MultiMatch(match, repl.re)
	before := String(matches["before"])
	context := String(matches["context"])
	context, strings := context.Protect()
	args := context.SelectContext(1, "(", ")")
	if args == "" {
		log.Warningf("Missing closing parenthesis in %s%s", matches["func"], context.RestoreProtected(strings))
		return match
	}
	after := context[len(args):]
	return fmt.Sprintf("@%s%sjoin(\"%s\", formatList(\"%s%%v%s\", %s)",
		matches["nl"],
		matches["reduce"],
		iif(matches["nl"] != "", "\\n", ""),
		before.Escape(),
		after.RestoreProtected(strings).Escape(),
		args.RestoreProtected(strings)[1:],
	)
}
