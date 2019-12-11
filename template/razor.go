package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Add additional functions to the go template context
func (t *Template) applyRazor(content []byte) (result []byte, changed bool) {
	if !t.options[Razor] || !t.IsRazor(string(content)) {
		return content, false
	}
	t.ensureInit()

	for _, r := range replacementsInit[fmt.Sprint(t.delimiters)] {
		printDebugInfo(r, string(content))
		if r.parser == nil {
			content = r.re.ReplaceAll(content, []byte(r.replace))
		} else {
			content = r.re.ReplaceAllFunc(content, func(match []byte) []byte {
				return []byte(r.parser(r, string(match)))
			})
		}
	}
	content = []byte(strings.Replace(string(content), funcCall, "", -1))
	InternalLog.Infof("Generated content\n\n%s\n", color.HiCyanString(String(content).AddLineNumber(0).Str()))
	return content, true
}

var highlight = color.New(color.BgHiBlack, color.FgBlack).SprintFunc()

// This is indented to simplify the following regular expression for patterns that are repeated several times
// Warning: The declaration order is important
var customMetaclass = [][2]string{
	{"function;", `@reduce;(?P<expr>[id]\([sp][expr]*[sp]\))`},
	{"assign;", `(?P<assign>(?:\$[id][ \t,]*){1,2}:=[sp])`},       // Optional assignment
	{"index;", `(?P<index>\[[expr]+\])`},                          // Extended index operator that support picky selection using ',' as element separator
	{"selector;", `(?P<selection>\.[expr]+)`},                     // Optional selector following expression indicating that the expression must include the content after the closing ) (i.e. @function(args).selection)
	{"reduce;", `(?P<reduce>((?P<reduce1>-)|_)?(?P<reduce2>-?))`}, // Optional reduces sign (-) indicating that the generated code must start with {{- (and end with -}} if two dashes are specified @--)
	{"endexpr;", `(?:[sp];)?`},                                    // End expression (spaces + ; or end of line)
	{"[sp]", `[[:blank:]]*`},                                      // Optional spaces
	{"[id]", `[\p{L}\d_]+`},                                       // Go language id
	{"[id_comp]", `[\p{L}\d_\.]+`},                                // Go language id with .
	{"[flexible_id]", `[map_id;][map_id;\.\-]*`},                  // Id with additional character that could be used to create variables in maps
	{"[idSel]", `[\p{L}_][\p{L}\d_\.]*`},                          // Id with optional selection (object.selection.subselection)
	{"map_id;", `\p{L}\d_\+\*%#!~`},                               // Id with additional character that could be used to create variables in maps

	{"assign_op;", `:=|~=|=|\+=|-=|\*=|\/=|÷=|%=|&=|\|=|\^=|<<=|«=|>>=|»=|&\^=`}, // Assignment operator
}

// Expression (any character that is not a new line, a start of razor expression or a semicolumn)
var expressionList = []string{`[^\n]`, `[^@\n]`, `[^@;\n]`, `[^@{\n]`, `[^@;{\n]`, fmt.Sprintf(`[^@;{\n'%s]`, "`"), `[\p{L}_\.]`}

const expressionKey = "[expr]"

// Warning: The declaration order is important
var expressions = [][]interface{}{
	// Literals
	{"Protect email", `(\W|^)[\w.!#$%&'*+/=?^_{|}~-]+@[\w-]{1,61}(?:\.[\w-]{1,61})+`, "", replacementFunc(protectEmail)},
	{"", `\${`, literalReplacement},
	{"", `@@`, literalAt},
	{"", `@{{`, literalStart},
	{"", "(?s)`+.*?`+", "", replacementFunc(protectMultiLineStrings)},
	{"", `@<;`, `{{- $.NEWLINE }}`},
	{"Auto indent", `(?m)^(?P<before>.*)@reduce;(?:autoIndent|aindent|aIndent)\(`, "@<-spaceIndent(`${before}`, "},
	{"Auto wrap", `(?m)^(?P<before>.*)@(?P<nl><)?reduce;(?P<func>autoWrap|awrap|aWrap)(?P<context>\(.*)$`, "", replacementFunc(autoWrap)},
	{`Inline content "@<...>"`, `"@reduce;<(?P<content>.*?)>"`, `"<<@${reduce}(${content})"`},
	{"Newline expression", `@<`, `{{- $.NEWLINE }}@`},

	// Comments
	{"Pseudo line comments - #! @", `(?m)(?:[sp](?:#|//)![sp])@`, "@"},
	{"Pseudo block comments - /*@  @*/", `(?s)/\*@(?P<content>.*?)@\*/`, "${content}"},
	{"Real comments - ##|/// @ comment", `(?m)[sp](?:##|///)[sp]@.*$`, `{{- "" }}`},
	{"Line comment - @// or @#", `(?m)@reduce;(#|//)[sp](?P<line_comment>.*)[sp]$`, "{{${reduce1} /* ${line_comment} */ ${reduce2}}}"},
	{"Block comment - @/* */", `(?s)@reduce;/\*(?P<block_comment>.*?)\*/`, "{{${reduce1} /*${block_comment}*/ ${reduce2}}}"},
	{"", `{{ /\*`, `{{/*`}, {"", `\*/ }}`, `*/}}`}, // Gotemplate is picky about spaces around comment {{- /* comment */ -}} and {{/* comment */}} are valid, but {{-/* comment */-}} and {{ /* comment */ }} are not.

	// Commands
	{"Foreach", `@reduce;for(?:[sp]each)?[sp]\(`, "@${reduce}range("},
	{"Single line command - @command (expr) action;", `@reduce;(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp](?P<action>[^\n]+?)[sp];`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}${action}{{${reduce1} end ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Single line command - @command (expr) { action }", `(?m)@reduce;(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]{[sp](?P<action>[^\n]+?)}[sp]$`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}${action}{{${reduce1} end ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Command(expr)", `@reduce;(?P<command>if|else[sp]if|block|with|define|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"else", `@reduce;else`, "{{${reduce1} else ${reduce2}}}"},
	{"various ends", `@reduce;(?P<command>end[sp](if|range|define|block|with|for[sp]each|for|))endexpr;`, "{{${reduce1} end ${reduce2}}}"},

	// Assignations
	{"Assign - @var := value", `(?P<type>@(\$|\.|\$\.)?)(?P<id>[flexible_id])[sp](?P<assign>assign_op;)[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},
	{"Assign - @{var} := value", `(?P<type>@{)(?P<id>[id_comp])}[sp](?P<assign>assign_op;)[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},
	{"Assign - @{var := expr}", `(?P<type>@{)(?P<id>[id_comp])[sp](?P<assign>assign_op;)[sp](?P<expr>[expr]+?)}endexpr;`, ``, replacementFunc(assignExpressionAcceptError)},

	// Function calls
	{"Function call followed by expression - @func(args...).args", `function;selector;endexpr;`, `@${reduce}((${expr})${selection});`, replacementFunc(expressionParserSkipError)},
	{"Function call with slice - @func(args...)[...]", `function;index;endexpr;`, `{{${reduce1} ${slicer} (${expr}) ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Function call - @func(args...)", `function;endexpr;`, `{{${reduce1} ${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Function unmanaged - @func(value | func)", `@reduce;(?P<function>[id])\([sp](?P<args>[expr]+)[sp]\)endexpr;`, `{{${reduce1} ${function} ${args} ${reduce2}}}`},

	// Variables
	{"Local variables - @{var}", `@reduce;{[sp](?P<name>[\p{L}\d_\.]*)[sp]}(?P<end>endexpr;)`, `@${reduce}($$${name});`},
	{"Global variables followed by expression", `@reduce;(?P<expr>[idSel]selector;index;?)(?P<end>endexpr;)`, `@${reduce}(${expr});`, replacementFunc(expressionParserSkipError)},
	{"Context variables - @.var", `@reduce;\.(?P<name>[idSel])endexpr;`, `@${reduce}(.${name})`},
	{"Global variables with slice - @var[...]", `@reduce;(?P<name>[idSel])index;endexpr;`, `{{${reduce1} ${slicer} $$.${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Context variables special with slice", `@reduce;\.(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce1} ${slicer} (get . "${name}") ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Global variables special with slice", `@reduce;(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce1} ${slicer} (get $$ "${name}") ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Local variables with slice", `@reduce;(?P<expr>(?P<name>[\$\.][\p{L}\d_\.]*)index;)endexpr;`, `{{${reduce1} ${slicer} ${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Global variables - @var", `@reduce;(?P<name>[idSel])endexpr;`, `{{${reduce1} $$.${name} ${reduce2}}}`},
	{"Context variables special - @.var", `@reduce;\.(?P<name>[flexible_id])endexpr;`, `{{${reduce1} get . "${name}" ${reduce2}}}`},
	{"Global variables special - @var", `@reduce;(?P<name>[flexible_id])endexpr;`, `{{${reduce1} get $$ "${name}" ${reduce2}}}`},
	{"Local variables - @$var or @.var", `@reduce;(?P<name>[\$\.][\p{L}\d_\.]*)endexpr;`, `{{${reduce1} ${name} ${reduce2}}}`},

	// Expressions
	{"Expression @(var)[...]", `@reduce;(?P<expr>\([sp](?P<name>[idSel])[sp]\)index;)endexpr;`, `{{${reduce1} ${slicer} $$.${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(var)", `@reduce;\([sp](?P<expr>[idSel])[sp]\)endexpr;`, `{{${reduce1} ${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)[...]", `@reduce;\([sp](?P<expr>[expr]+)[sp]\)index;endexpr;`, `{{${reduce1} ${slicer} (${expr}) ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)", `@reduce;\([sp](?P<assign>.*?:= ?)?[sp](?P<expr>[expr]+)[sp]\)endexpr;`, `{{${reduce1} ${assign}${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},

	{"Space eater", `@-endexpr;`, `{{- "" -}}`},

	// Inline contents: Render the content without its enclosing quotes
	{`Inline content "<<..."`, `"<<(?P<content>{{[sp].*[sp]}})"`, `${content}`},

	// Restoring literals
	{"", `}}\\\.`, "}}."},
	{"", literalAt, "@"},
	{"", literalReplacement, "${"},
	{"", fmt.Sprintf(`\x60%s(?P<num>\d+)\x60`, protectString), "", replacementFunc(protectMultiLineStrings)},
}

var replacementsInit = make(map[string][]replacement)

type replacementFunc func(replacement, string) string
type replacement struct {
	name       string
	expr       string
	replace    string
	re         *regexp.Regexp
	parser     replacementFunc
	delimiters []string
}

func (t *Template) ensureInit() {
	delimiters := fmt.Sprint(t.delimiters)
	if _, ok := replacementsInit[delimiters]; !ok {
		// We must ensure that search and replacement expression are compatible with the set of delimiters
		replacements := make([]replacement, 0, len(expressions))
		for _, expr := range expressions {
			comment := expr[0].(string)
			re := strings.Replace(expr[1].(string), "@", regexp.QuoteMeta(t.delimiters[2]), -1)
			re = strings.Replace(re, "{{", regexp.QuoteMeta(t.delimiters[0]), -1)
			re = strings.Replace(re, "}}", regexp.QuoteMeta(t.delimiters[1]), -1)
			replace := strings.Replace(strings.Replace(strings.Replace(expr[2].(string), "{{", t.delimiters[0], -1), "}}", t.delimiters[1], -1), "@", t.delimiters[2], -1)
			var exprParser replacementFunc
			if len(expr) >= 4 {
				exprParser = expr[3].(replacementFunc)
			}

			// We apply replacements in regular expression to make them regex compliant
			for i := range customMetaclass {
				key, value := customMetaclass[i][0], customMetaclass[i][1]
				re = strings.Replace(re, key, value, -1)
			}

			subExpressions := []string{re}
			if strings.Contains(re, expressionKey) {
				// If regex contains the generic expression token [expr], we generate several expression evaluator
				// that go from the most generic expression to the most specific one
				subExpressions = make([]string, len(expressionList))
				for i := range expressionList {
					subExpressions[i] = strings.Replace(re, expressionKey, expressionList[i], -1)
				}
			}

			for i := range subExpressions {
				re := regexp.MustCompile(subExpressions[i])
				replacements = append(replacements, replacement{comment, subExpressions[i], replace, re, exprParser, t.delimiters})
			}

			if len(subExpressions) > 1 && len(expr) == 5 {
				// If there is a fallback expression evaluator, we apply it on the first replacement alternative
				re := regexp.MustCompile(subExpressions[0])
				replacements = append(replacements, replacement{comment, subExpressions[0], replace, re, expr[4].(replacementFunc), t.delimiters})
			}
		}
		replacementsInit[delimiters] = replacements
	}
}

func printDebugInfo(r replacement, content string) {
	if r.name == "" || !InternalLog.IsLevelEnabled(logrus.DebugLevel) {
		return
	}

	debugMode = true
	defer func() { debugMode = false }()

	// We only report each match once
	allUnique := make(map[string]int)
	for _, found := range r.re.FindAllString(content, -1) {
		if r.parser != nil {
			newContent := r.re.ReplaceAllStringFunc(found, func(match string) string {
				return r.parser(r, match)
			})
			if newContent == found && !InternalLog.IsLevelEnabled(logrus.TraceLevel) {
				// There is no change
				continue
			}
		}
		lines := strings.Split(found, "\n")
		for i := range lines {
			lines[i] = fmt.Sprintf("%v%s", iif(i > 0, "  ", ""), highlight(lines[i]))
		}
		found = strings.Join(lines, "\n")
		allUnique[found] = allUnique[found] + 1
	}

	if len(allUnique) == 0 && !InternalLog.IsLevelEnabled(logrus.TraceLevel) {
		return
	}

	matches := make([]string, 0, len(allUnique)+1)
	if len(allUnique) > 0 {
		matches = append(matches, "")
	}
	for key, count := range allUnique {
		if count > 1 {
			key = fmt.Sprintf("%s (%d)", key, count)
		}
		matches = append(matches, key)
	}
	logString := fmt.Sprintf("%s: %s%s", color.YellowString(r.name), r.expr, strings.Join(matches, "\n- "))
	if len(matches) > 0 {
		logString = logString + " \n"
	}
	InternalLog.Debug(logString)
}

var debugMode bool
