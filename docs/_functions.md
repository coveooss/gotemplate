```text
Base go template functions

and                     gt                      le                      not                     println                 
call                    html                    len                     or                      urlquery                
eq                      index                   lt                      print                   
ge                      js                      ne                      printf                  

Data Conversion

data                    toHcl                   toPrettyJson            toQuotedTFVars          
hcl                     toInternalHcl           toPrettyTFVars          toTFVars                
json                    toJson                  toQuotedHcl             toYaml                  
toBash                  toPrettyHcl             toQuotedJson            yaml                    

Data Manipulation

String                  dict                    keys                    pickv                   string                  
array                   extract                 lenc                    pluck                   undef                   
bool                    get                     merge                   safeIndex               unset                   
char                    hasKey                  omit                    set                     
content                 key                     pick                    slice                   

Logging

critical                error                   info                    panic                   
criticalf               errorf                  infof                   panicf                  
debug                   fatal                   notice                  warning                 
debugf                  fatalf                  noticef                 warningf                

Mathematic Bit Operations

band                    bor                     lshift                  
bclear                  bxor                    rshift                  

Mathematic Fundamental

add                     div                     floor                   pow                     trunc                   
cbrt                    exp                     mod                     pow10                   
ceil                    exp2                    modf                    rem                     
dim                     expm1                   mul                     sub                     

Mathematic Stats

avg                     max                     min                     

Mathematic Trigonometry

acos                    atanh                   j1                      logb                    tanh                    
acosh                   cos                     jn                      rad                     y0                      
asin                    cosh                    log                     sin                     y1                      
asinh                   deg                     log10                   sincos                  yn                      
atan                    ilogb                   log1p                   sinh                    
atan2                   j0                      log2                    tan                     

Mathematic Utilities

abs                     gamma                   isInf                   lgamma                  sqrt                    
dec                     hex                     isNaN                   nextAfter               to                      
frexp                   hypot                   ldexp                   signBit                 until                   

Other utilities

center                  formatList              joinLines               repeat                  
color                   glob                    lorem                   sIndent                 
concat                  id                      mergeList               splitLines              
diff                    iif                     pwd                     wrap                    

Runtime

alias                   current                 func                    localAlias              templates               
aliases                 ellipsis                function                run                     
allFunctions            exec                    functions               substitute              
categories              exit                    include                 templateNames           

Sprig Cryptographic & Security http://masterminds.github.io/sprig/crypto.html

derivePassword          genPrivateKey           genSignedCert           
genCA                   genSelfSignedCert       sha256sum               

Sprig Date http://masterminds.github.io/sprig/date.html

ago                     dateInZone              htmlDate                now                     
date                    dateModify              htmlDateInZone          toDate                  

Sprig Default http://masterminds.github.io/sprig/defaults.html

coalesce                compact                 default                 empty                   

Sprig Dictionnary http://masterminds.github.io/sprig/dicst.html

list                    

Sprig Encoding http://masterminds.github.io/sprig/encoding.html

b32dec                  b32enc                  b64dec                  b64enc                  

Sprig File Path http://masterminds.github.io/sprig/paths.html

base                    clean                   dir                     ext                     isAbs                   

Sprig Flow Control http://masterminds.github.io/sprig/flow_control.html

fail                    

Sprig General http://masterminds.github.io/sprig/

hello                   uuidv4                  

Sprig List http://masterminds.github.io/sprig/lists.html

append                  initial                 prepend                 reverse                 without                 
first                   last                    rest                    uniq                    

Sprig OS http://masterminds.github.io/sprig/defaults.html

env                     expandenv               

Sprig Reflection http://masterminds.github.io/sprig/reflection.html

kindIs                  kindOf                  typeIs                  typeIsLike              typeOf                  

Sprig Regex http://masterminds.github.io/sprig/strings.html

regexFind               regexMatch              regexReplaceAllLiteral  
regexFindAll            regexReplaceAll         regexSplit              

Sprig Sprig Mathematics http://masterminds.github.io/sprig/math.html

add1                    round                   untilStep               

Sprig String Slice http://masterminds.github.io/sprig/string_slice.html

join                    sortAlpha               split                   splitList               toStrings               

Sprig Strings http://masterminds.github.io/sprig/strings.html

abbrev                  indent                  randAlpha               squote                  trimPrefix              
abbrevboth              initials                randAlphaNum            substr                  trimSuffix              
camelcase               lower                   randAscii               swapcase                untitle                 
cat                     nindent                 randNumeric             title                   upper                   
contains                nospace                 replace                 toString                wrapWith                
hasPrefix               plural                  shuffle                 trim                    
hasSuffix               quote                   snakecase               trimAll                 

Sprig Type Conversion http://masterminds.github.io/sprig/conversion.html

atoi                    float64                 int64                   

Sprig Version comparison http://masterminds.github.io/sprig/semver.html

semver                  semverCompare           
```
