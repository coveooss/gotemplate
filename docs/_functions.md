```text
Base go template functions

and                     gt                      le                      not                     println                 
call                    html                    len                     or                      urlquery                
eq                      index                   lt                      print                   
ge                      js                      ne                      printf                  

Data Conversion

data                    toBash                  toPrettyHcl             toQuotedHcl             toTFVars                
hcl                     toHcl                   toPrettyJson            toQuotedJson            toYaml                  
json                    toJson                  toPrettyTFVars          toQuotedTFVars          yaml                    

Data Manipulation

array                   extract                 merge                   safeIndex               undef                   
bool                    get                     omit                    set                     
char                    key                     pick                    slice                   
content                 lenc                    pickv                   string                  

Logging

debug                   errorf                  info                    noticef                 
debugf                  fatal                   infof                   warning                 
error                   fatalf                  notice                  warningf                

Mathematic Bit Operations

band                    bor                     lshift                  
bclear                  bxor                    rshift                  

Mathematic Fundamental

add                     exp                     mod                     pow10                   
ceil                    exp2                    modf                    rem                     
dim                     expm1                   mul                     sub                     
div                     floor                   pow                     trunc                   

Mathematic Stats

avg                     max                     min                     

Mathematic Trigonometry

acos                    atan2                   ilogb                   logb                    tan                     
acosh                   atanh                   log                     rad                     tanh                    
asin                    cos                     log10                   sin                     
asinh                   cosh                    log1p                   sincos                  
atan                    deg                     log2                    sinh                    

Mathematic Utilities

abs                     gamma                   isInf                   lgamma                  sqrt                    
dec                     hex                     isNaN                   nextAfter               to                      
frexp                   hypot                   ldexp                   signBit                 until                   

Other utilities

center                  diff                    id                      lorem                   repeat                  
color                   formatList              iif                     mergeList               splitLines              
concat                  glob                    joinLines               pwd                     wrap                    

Runtime

alias                   exec                    function                localAlias              templateNames           
current                 exit                    functions               run                     
ellipsis                func                    include                 substitute              

Sprig Cryptographic & Security http://masterminds.github.io/sprig/crypto.html

derivePassword          genPrivateKey           genSignedCert           
genCA                   genSelfSignedCert       sha256sum               

Sprig Date http://masterminds.github.io/sprig/date.html

ago                     dateInZone              htmlDate                now                     
date                    dateModify              htmlDateInZone          toDate                  

Sprig Default http://masterminds.github.io/sprig/defaults.html

coalesce                compact                 default                 empty                   

Sprig Dictionnary http://masterminds.github.io/sprig/dicst.html

dict                    keys                    pluck                   
hasKey                  list                    unset                   

Sprig Encoding http://masterminds.github.io/sprig/encoding.html

b32dec                  b32enc                  b64dec                  b64enc                  

Sprig File Path http://masterminds.github.io/sprig/paths.html

base                    clean                   dir                     ext                     isAbs                   

Sprig Flow Control http://masterminds.github.io/sprig/flow_control.html

fail                    

Sprig General http://masterminds.github.io/sprig/

hello                   uuidv4                  

Sprig List http://masterminds.github.io/sprig/lists.html

append                  has                     last                    rest                    uniq                    
first                   initial                 prepend                 reverse                 without                 

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
