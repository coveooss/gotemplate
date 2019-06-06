```text
usage: gotemplate [<flags>] <command> [<args> ...]

An extended template processor for go.

See: https://coveo.github.io/gotemplate for complete documentation.

Flags:
  -h, --help                     Show context-sensitive help (also try --help-man). or set GOTEMPLATE_HELP
      --color                    Force rendering of colors event if output is redirected or set GOTEMPLATE_COLOR (alias: --c)
  -v, --version                  Get the current version of gotemplate
      --base                     Turn off all addons (they could then be enabled explicitly) or set GOTEMPLATE_BASE
      --razor                    Razor Addon (ON by default) or set GOTEMPLATE_RAZOR (off: --no-razor)
      --extension                Extension Addon (ON by default) or set GOTEMPLATE_EXTENSION (alias: --ext) (off: --next, --no-ext, --no-extension)
      --math                     Math Addon (ON by default) or set GOTEMPLATE_MATH (off: --no-math)
      --sprig                    Sprig Addon (ON by default) or set GOTEMPLATE_SPRIG (off: --no-sprig)
      --data                     Data Addon (ON by default) or set GOTEMPLATE_DATA (off: --no-data)
      --logging                  Logging Addon (ON by default) or set GOTEMPLATE_LOGGING (off: --no-logging)
      --runtime                  Runtime Addon (ON by default) or set GOTEMPLATE_RUNTIME (off: --no-runtime)
      --utils                    Utils Addon (ON by default) or set GOTEMPLATE_UTILS (off: --no-utils)
      --net                      Net Addon (ON by default) or set GOTEMPLATE_NET (off: --nnet, --no-net)
      --os                       OS Addon (ON by default) or set GOTEMPLATE_OS (off: --no-os, --nos)

Args:
  [<templates>]  Template files or commands to process

Commands:
  help [<command>...]
    Show help.


  run [<flags>] [<templates>...]

        --delimiters={{,}},@       Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) or set
                                   GOTEMPLATE_DELIMITERS (alias: --d, --del)
    -i, --import=file ...          Import variables files (could be any of YAML, JSON or HCL format) or set GOTEMPLATE_IMPORT
        --import-if-exist=file ...  
                                   Import variables files (do not consider missing file as an error) or set GOTEMPLATE_IMPORT_IF_EXIST (alias: --iie)
    -V, --var=values ...           Import named variables (if value is a file, the content is loaded) or set GOTEMPLATE_VAR
    -t, --type=TYPE                Force the type used for the main context (Json, Yaml, Hcl) or set GOTEMPLATE_TYPE
    -p, --patterns=pattern ...     Additional patterns that should be processed by gotemplate or set GOTEMPLATE_PATTERNS
    -e, --exclude=pattern ...      Exclude file patterns (comma separated) when applying gotemplate recursively or set GOTEMPLATE_EXCLUDE
    -o, --overwrite                Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target
                                   folder) or set GOTEMPLATE_OVERWRITE
    -s, --substitute=exp ...       Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution,
                                   the first character acts as separator like in sed, see: Go regexp) or set GOTEMPLATE_SUBSTITUTE
    -E, --remove-empty-lines       Remove empty lines from the result or set GOTEMPLATE_REMOVE_EMPTY_LINES (alias: --re, --rel, --remove-empty)
    -r, --recursive                Process all template files recursively or set GOTEMPLATE_RECURSIVE
    -R, --recursion-depth=depth    Process template files recursively specifying depth or set GOTEMPLATE_RECURSION_DEPTH (alias: --rd)
        --source=folder            Specify a source folder (default to the current folder) or set GOTEMPLATE_SOURCE (alias: --s)
        --target=folder            Specify a target folder (default to source folder) or set GOTEMPLATE_TARGET (alias: --t)
    -I, --stdin                    Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set) or set
                                   GOTEMPLATE_STDIN
    -f, --follow-symlinks          Follow the symbolic links while using the recursive option or set GOTEMPLATE_FOLLOW_SYMLINKS (alias: --fs)
    -P, --print                    Output the result directly to stdout or set GOTEMPLATE_PRINT
    -d, --disable                  Disable go template rendering (used to view razor conversion) or set GOTEMPLATE_DISABLE
        --accept-no-value          Do not consider rendering <no value> as an error or set GOTEMPLATE_NO_VALUE (alias: --anv, --no-value, --nv)
    -S, --strict-error-validation  Consider error encountered in any file as real error or set GOTEMPLATE_STRICT_ERROR (alias: --sev, --strict)
    -L, --log-level=level          Set the logging level CRITICAL (0), ERROR (1), WARNING (2), NOTICE (3), INFO (4), DEBUG (5) or set
                                   GOTEMPLATE_LOG_LEVEL (alias: --ll)
        --debug-log-level=level    Set the debug logging level 0-9 or set GOTEMPLATE_DEBUG_LOG_LEVEL (alias: --debug-level, --dl, --dll)
        --log-simple               Disable the extended logging, i.e. no color, no date or set GOTEMPLATE_LOG_SIMPLE (alias: --ls)
        --ignore-missing-import    Exit with code 0 even if import does not exist or set GOTEMPLATE_IGNORE_MISSING_IMPORT (alias: --imi)
        --ignore-missing-source    Exit with code 0 even if source does not exist or set GOTEMPLATE_IGNORE_MISSING_SOURCE (alias: --ims)
        --ignore-missing-paths     Exit with code 0 even if import or source do not exist or set GOTEMPLATE_IGNORE_MISSING_PATHS (alias: --imp)

  list [<flags>] [<filters>...]
    Get detailed help on gotemplate functions

    -f, --functions  Get detailed help on function
    -t, --templates  List the available templates
    -l, --long       Get detailed list
    -a, --all        List all
    -c, --category   Group functions by category
```
