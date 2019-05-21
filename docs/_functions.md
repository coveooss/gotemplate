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

String                  extract                 key                     pluck                   undef                   
append                  get                     keys                    prepend                 union                   
array                   hasKey                  lenc                    rest                    unique                  
bool                    initial                 list                    reverse                 unset                   
char                    intersect               merge                   safeIndex               values                  
contains                isNil                   omit                    set                     without                 
content                 isSet                   pick                    slice                   
dict                    isZero                  pickv                   string                  

Logging

critical                error                   info                    panic                   
debug                   fatal                   notice                  warning                 

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

Net

httpDoc                 httpGet                 

Operating systems functions

diff                    home                    isReadable              mode                    stat                    
exists                  isDir                   isWriteable             pwd                     user                    
glob                    isExecutable            lastMod                 save                    username                
group                   isFile                  lookPath                size                    

Other utilities

center                  formatList              indent                  mergeList               sIndent                 
color                   id                      joinLines               nIndent                 splitLines              
concat                  iif                     lorem                   repeat                  wrap                    

Runtime

alias                   categories              func                    getSignature            substitute              
aliases                 current                 function                include                 templateNames           
allFunctions            ellipsis                functions               localAlias              templates               
assert                  exec                    getAttributes           raise                   
assertWarning           exit                    getMethods              run                     

Sprig Cryptographic & Security http://masterminds.github.io/sprig/crypto.html

adler32sum              derivePassword          genPrivateKey           genSignedCert           sha256sum               
buildCustomCert         genCA                   genSelfSignedCert       sha1sum                 

Sprig Date http://masterminds.github.io/sprig/date.html

ago                     dateInZone              htmlDate                now                     
date                    dateModify              htmlDateInZone          toDate                  

Sprig Default http://masterminds.github.io/sprig/defaults.html

coalesce                default                 ternarySprig            toPrettyJsonSprig       
compact                 empty                   toJsonSprig             

Sprig Dictionnary http://masterminds.github.io/sprig/dicst.html

dictSprig               listSprig               omitSprig               setSprig                
hasKeySprig             mergeOverwrite          pickSprig               unsetSprig              
keysSprig               mergeSprig              pluckSprig              valuesSprig             

Sprig Encoding http://masterminds.github.io/sprig/encoding.html

b32dec                  b32enc                  b64dec                  b64enc                  

Sprig File Path http://masterminds.github.io/sprig/paths.html

base                    clean                   dir                     ext                     isAbs                   

Sprig Flow Control http://masterminds.github.io/sprig/flow_control.html

fail                    

Sprig General http://masterminds.github.io/sprig/

hello                   uuidv4                  

Sprig List http://masterminds.github.io/sprig/lists.html

appendSprig             initialSprig            restSprig               uniqSprig               
first                   last                    reverseSprig            withoutSprig            
hasSprig                prependSprig            sliceSprig              

Sprig Mathematics http://masterminds.github.io/sprig/math.html

add1                    divSprig                minSprig                round                   
addSprig                floorSprig              modSprig                subSprig                
ceilSprig               maxSprig                mulSprig                untilStep               

Sprig OS http://masterminds.github.io/sprig/defaults.html

env                     expandenv               

Sprig Reflection http://masterminds.github.io/sprig/reflection.html

kindIs                  kindOf                  typeIs                  typeIsLike              typeOf                  

Sprig Regex http://masterminds.github.io/sprig/strings.html

regexFind               regexMatch              regexReplaceAllLiteral  
regexFindAll            regexReplaceAll         regexSplit              

Sprig String Slice http://masterminds.github.io/sprig/string_slice.html

join                    split                   splitn                  
sortAlpha               splitList               toStrings               

Sprig Strings http://masterminds.github.io/sprig/strings.html

abbrev                  initials                randAlphaNum            substr                  truncSprig              
abbrevboth              kebabcase               randAscii               swapcase                untitle                 
camelcase               lower                   randNumeric             title                   upper                   
cat                     nindentSprig            repeatSprig             toString                wrapSprig               
containsSprig           nospace                 replace                 trim                    wrapWith                
hasPrefix               plural                  shuffle                 trimAll                 
hasSuffix               quote                   snakecase               trimPrefix              
indentSprig             randAlpha               squote                  trimSuffix              

Sprig Type Conversion http://masterminds.github.io/sprig/conversion.html

atoi                    float64                 int64                   intSprig                

Sprig Version comparison http://masterminds.github.io/sprig/semver.html

semver                  semverCompare           

```
