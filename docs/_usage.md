```text
usage: gotemplate [<flags>] <command> [<args> ...]

An extended template processor for go.

See: https://github.com/coveo/gotemplate/blob/master/README.md for complete documentation.

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --color                    Force rendering of colors event if output is redirected
      --no-color                 Force rendering of colors event if output is redirected
  -v, --version                  Get the current version of gotemplate
      --base                     Turn off all extensions (they could then be enabled explicitly)
      --razor                    Option Razor, on by default, --no-razor to disable
      --extension                Option Extension, on by default, --no-extension to disable
      --math                     Option Math, on by default, --no-math to disable
      --sprig                    Option Sprig, on by default, --no-sprig to disable
      --data                     Option Data, on by default, --no-data to disable
      --logging                  Option Logging, on by default, --no-logging to disable
      --runtime                  Option Runtime, on by default, --no-runtime to disable
      --utils                    Option Utils, on by default, --no-utils to disable
      --net                      Option Net, on by default, --no-net to disable
      --os                       Option OS, on by default, --no-os to disable
      --delimiters={{,}},@       Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)
  -i, --import=file ...          Import variables files (could be any of YAML, JSON or HCL format)
  -V, --var=values ...           Import named variables (if value is a file, the content is loaded)
  -t, --type=TYPE                Force the type used for the main context (Json, Yaml, Hcl)
  -p, --patterns=pattern ...     Additional patterns that should be processed by gotemplate
  -e, --exclude=pattern ...      Exclude file patterns (comma separated) when applying gotemplate recursively
  -o, --overwrite                Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
  -s, --substitute=exp ...       Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp) or specify that value through
                                 GOTEMPLATE_SUBSTITUTES where each substitute is separated by a newline
  -E, --remove-empty-lines       Remove empty lines from the result (--re)
  -r, --recursive                Process all template files recursively
  -R, --recursion-depth=depth    Process template files recursively specifying depth
      --source=folder            Specify a source folder (default to the current folder)
      --target=folder            Specify a target folder (default to source folder)
  -I, --stdin                    Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)
  -f, --follow-symlinks          Follow the symbolic links while using the recursive option
  -P, --print                    Output the result directly to stdout
  -d, --disable                  Disable go template rendering (used to view razor conversion)
      --accept-no-value          Do not consider rendering <no value> as an error (--nv) or env: GOTEMPLATE_NO_VALUE
  -S, --strict-error-validation  Consider error encountered in any file as real error (--strict) or env: GOTEMPLATE_STRICT_ERROR
  -L, --log-level=level          Set the logging level CRITICAL (0), ERROR (1), WARNING (2), NOTICE (3), INFO (4), DEBUG (5) (--ll)
      --debug-log-level=level    Set the debug logging level 0-9 (--dl) or env: GOTEMPLATE_DEBUG
      --log-simple               Disable the extended logging, i.e. no color, no date (--ls)
      --ignore-missing-import    Exit with code 0 even if import does not exist
      --ignore-missing-source    Exit with code 0 even if source does not exist
      --ignore-missing-paths     Exit with code 0 even if import or source do not exist

Args:
  [<templates>]  Template files or commands to process

Commands:
  help [<command>...]
    Show help.


  run [<flags>] [<templates>...]

        --delimiters={{,}},@       Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)
    -i, --import=file ...          Import variables files (could be any of YAML, JSON or HCL format)
    -V, --var=values ...           Import named variables (if value is a file, the content is loaded)
    -t, --type=TYPE                Force the type used for the main context (Json, Yaml, Hcl)
    -p, --patterns=pattern ...     Additional patterns that should be processed by gotemplate
    -e, --exclude=pattern ...      Exclude file patterns (comma separated) when applying gotemplate recursively
    -o, --overwrite                Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
    -s, --substitute=exp ...       Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp) or specify that value through
                                   GOTEMPLATE_SUBSTITUTES where each substitute is separated by a newline
    -E, --remove-empty-lines       Remove empty lines from the result (--re)
    -r, --recursive                Process all template files recursively
    -R, --recursion-depth=depth    Process template files recursively specifying depth
        --source=folder            Specify a source folder (default to the current folder)
        --target=folder            Specify a target folder (default to source folder)
    -I, --stdin                    Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)
    -f, --follow-symlinks          Follow the symbolic links while using the recursive option
    -P, --print                    Output the result directly to stdout
    -d, --disable                  Disable go template rendering (used to view razor conversion)
        --accept-no-value          Do not consider rendering <no value> as an error (--nv) or env: GOTEMPLATE_NO_VALUE
    -S, --strict-error-validation  Consider error encountered in any file as real error (--strict) or env: GOTEMPLATE_STRICT_ERROR
    -L, --log-level=level          Set the logging level CRITICAL (0), ERROR (1), WARNING (2), NOTICE (3), INFO (4), DEBUG (5) (--ll)
        --debug-log-level=level    Set the debug logging level 0-9 (--dl) or env: GOTEMPLATE_DEBUG
        --log-simple               Disable the extended logging, i.e. no color, no date (--ls)
        --ignore-missing-import    Exit with code 0 even if import does not exist
        --ignore-missing-source    Exit with code 0 even if source does not exist
        --ignore-missing-paths     Exit with code 0 even if import or source do not exist

  list [<flags>] [<filters>...]
    Get detailed help on gotemplate functions

    -f, --functions  Get detailed help on function
    -t, --templates  List the available templates
    -l, --long       Get detailed list
    -a, --all        List all
    -c, --category   Group functions by category
```
