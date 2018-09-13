```text
Base go template functions

and                     gt                      le                      not                     println                 
call                    html                    len                     or                      urlquery                
eq                      index                   lt                      print                   
ge                      js                      ne                      printf                  

Data Conversion

DATA                    fromHcl                 tfvars                  toJson                  toQuotedJson            
HCL                     fromJSON                toBash                  toPrettyHCL             toQuotedTFVars          
JSON                    fromJson                toHCL                   toPrettyHcl             toTFVars                
TFVARS                  fromTFVARS              toHcl                   toPrettyJSON            toYAML                  
YAML                    fromTFVars              toIHCL                  toPrettyJson            toYaml                  
data                    fromYAML                toIHcl                  toPrettyTFVars          yaml                    
fromDATA                fromYaml                toInternalHCL           toQuotedHCL             
fromData                hcl                     toInternalHcl           toQuotedHcl             
fromHCL                 json                    toJSON                  toQuotedJSON            

Data Manipulation

String                  extract                 isSet                   pickv                   string                  
append                  get                     isZero                  pluck                   undef                   
array                   has                     key                     prepend                 union                   
bool                    hasKey                  keys                    push                    uniq                    
char                    ifUndef                 lenc                    remove                  unique                  
contains                initial                 list                    rest                    unset                   
content                 intersect               merge                   reverse                 values                  
delete                  isEmpty                 nbChars                 safeIndex               without                 
dict                    isNil                   omit                    set                     
dictionary              isNull                  pick                    slice                   

Logging

critical                error                   info                    panic                   warning                 
criticalf               errorf                  infof                   panicf                  warningf                
debug                   fatal                   notice                  warn                    
debugf                  fatalf                  noticef                 warnf                   

Mathematic Bit Operations

band                    bitwiseClear            bor                     lshift                  
bclear                  bitwiseOR               bxor                    rightShift              
bitwiseAND              bitwiseXOR              leftShift               rshift                  

Mathematic Fundamental

add                     expm1                   modulo                  product                 sub                     
cbrt                    exponent                mul                     quotient                subtract                
ceil                    exponent2               multiply                rem                     sum                     
dim                     floor                   pow                     remainder               trunc                   
div                     int                     pow10                   roundDown               truncate                
divide                  integer                 power                   roundUp                 
exp                     mod                     power10                 rounddown               
exp2                    modf                    prod                    roundup                 

Mathematic Stats

average                 biggest                 maximum                 minimum                 
avg                     max                     min                     smallest                

Mathematic Trigonometry

acos                    arcTangent2             firstBessel0            log                     sine                    
acosh                   asin                    firstBessel1            log10                   sineCosine              
arcCosine               asinh                   firstBesselN            log1p                   sinh                    
arcCosinus              atan                    hyperbolicCosine        log2                    sinus                   
arcHyperbolicCosine     atan2                   hyperbolicCosinus       logb                    sinusCosinus            
arcHyperbolicCosinus    atanh                   hyperbolicSine          rad                     tan                     
arcHyperbolicSine       cos                     hyperbolicSinus         radian                  tangent                 
arcHyperbolicSinus      cosh                    hyperbolicTangent       secondBessel0           tanh                    
arcHyperbolicTangent    cosine                  ilogb                   secondBessel1           y0                      
arcSine                 cosinus                 j0                      secondBesselN           y1                      
arcSinus                deg                     j1                      sin                     yn                      
arcTangent              degree                  jn                      sincos                  

Mathematic Utilities

abs                     gamma                   hypotenuse              lgamma                  to                      
absolute                hex                     isInf                   nextAfter               until                   
dec                     hexa                    isInfinity              signBit                 
decimal                 hexaDecimal             isNaN                   sqrt                    
frexp                   hypot                   ldexp                   squareRoot              

Net

curl                    httpDoc                 httpDocument            httpGet                 

Operating systems functions

currentDir              fileStat                isExist                 lookPath                username                
currentUser             glob                    isFile                  mode                    whereIs                 
diff                    group                   isFolder                pwd                     which                   
difference              home                    isReadable              save                    write                   
exists                  homeDir                 isWriteable             size                    writeTo                 
expand                  homeFolder              lastMod                 stat                    
fileExists              isDir                   lastModification        type                    
fileMode                isDirectory             lastModificationTime    user                    
fileSize                isExecutable            look                    userGroup               

Other utilities

center                  enhanced                joinLines               sIndent                 wrap                    
centered                formatList              lorem                   sindent                 wrapped                 
color                   id                      loremIpsum              spaceIndent             
colored                 identifier              mergeList               splitLines              
concat                  iif                     repeat                  ternary                 

Runtime

alias                   attributes              exit                    getSignature            run                     
aliases                 categories              func                    include                 sign                    
allFunctions            current                 function                localAlias              signature               
assert                  ellipsis                functions               methods                 substitute              
assertion               exec                    getAttributes           raise                   templateNames           
attr                    execute                 getMethods              raiseError              templates               

Sprig Cryptographic & Security http://masterminds.github.io/sprig/crypto.html

buildCustomCert         genCA                   genSelfSignedCert       sha1sum                 
derivePassword          genPrivateKey           genSignedCert           sha256sum               

Sprig Date http://masterminds.github.io/sprig/date.html

ago                     dateInZone              date_in_zone            htmlDate                now                     
date                    dateModify              date_modify             htmlDateInZone          toDate                  

Sprig Default http://masterminds.github.io/sprig/defaults.html

coalesce                default                 ternarySprig            
compact                 empty                   toJsonSprig             

Sprig Dictionnary http://masterminds.github.io/sprig/dicst.html

dictSprig               mergeSprig              pluckSprig              unsetSprig              
hasKeySprig             omitSprig               setSprig                valuesSprig             
keysSprig               pickSprig               tuple                   

Sprig Encoding http://masterminds.github.io/sprig/encoding.html

b32dec                  b32enc                  b64dec                  b64enc                  

Sprig File Path http://masterminds.github.io/sprig/paths.html

base                    clean                   dir                     ext                     isAbs                   

Sprig Flow Control http://masterminds.github.io/sprig/flow_control.html

fail                    

Sprig General http://masterminds.github.io/sprig/

GUID                    guid                    hello                   uuid                    uuidv4                  

Sprig List http://masterminds.github.io/sprig/lists.html

appendSprig             initialSprig            restSprig               uniqSprig               
first                   last                    reverseSprig            withoutSprig            
hasSprig                prependSprig            sliceSprig              

Sprig OS http://masterminds.github.io/sprig/defaults.html

env                     expandenv               

Sprig Reflection http://masterminds.github.io/sprig/reflection.html

kindIs                  kindOf                  typeIs                  typeIsLike              typeOf                  

Sprig Regex http://masterminds.github.io/sprig/strings.html

regexFind               regexMatch              regexReplaceAllLiteral  
regexFindAll            regexReplaceAll         regexSplit              

Sprig Sprig Mathematics http://masterminds.github.io/sprig/math.html

add1                    divSprig                minSprig                round                   untilStep               
addSprig                floorSprig              modSprig                subSprig                
ceilSprig               maxSprig                mulSprig                untilSprig              

Sprig String Slice http://masterminds.github.io/sprig/string_slice.html

join                    split                   splitn                  
sortAlpha               splitList               toStrings               

Sprig Strings http://masterminds.github.io/sprig/strings.html

abbrev                  initials                randAscii               swapcase                truncSprig              
abbrevboth              lower                   randNumeric             title                   untitle                 
camelcase               nindent                 repeatSprig             toString                upper                   
cat                     nospace                 replace                 trim                    wrapSprig               
containsSprig           plural                  shuffle                 trimAll                 wrapWith                
hasPrefix               quote                   snakecase               trimPrefix              
hasSuffix               randAlpha               squote                  trimSuffix              
indent                  randAlphaNum            substr                  trimall                 

Sprig Type Conversion http://masterminds.github.io/sprig/conversion.html

atoi                    float64                 int64                   intSprig                

Sprig Version comparison http://masterminds.github.io/sprig/semver.html

semver                  semverCompare           

```
