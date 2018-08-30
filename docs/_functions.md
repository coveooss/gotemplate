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

String                  extract                 isZero                  pluck                   true                    
append                  false                   key                     prepend                 undef                   
array                   get                     keys                    push                    uniq                    
bool                    has                     lenc                    remove                  unique                  
char                    hasKey                  list                    rest                    unset                   
contains                ifUndef                 merge                   reverse                 values                  
content                 initial                 nbChars                 safeIndex               without                 
delete                  isEmpty                 omit                    set                     
dict                    isNil                   pick                    slice                   
dictionary              isNull                  pickv                   string                  

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

Other utilities

assert                  diff                    identifier              raise                   splitLines              
center                  difference              iif                     raiseError              ternary                 
centered                enhanced                joinLines               repeat                  wrap                    
color                   expand                  lorem                   sIndent                 wrapped                 
colored                 formatList              loremIpsum              save                    write                   
concat                  glob                    mergeList               sindent                 writeTo                 
currentDir              id                      pwd                     spaceIndent             

Runtime

alias                   current                 function                localAlias              templateNames           
aliases                 ellipsis                functions               methods                 templates               
allFunctions            exec                    getAttributes           run                     
attr                    execute                 getMethods              sign                    
attributes              exit                    getSignature            signature               
categories              func                    include                 substitute              

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

abbrev                  lower                   randNumeric             title                   untitle                 
abbrevboth              nindent                 repeatSprig             toString                upper                   
camelcase               nospace                 replace                 trim                    wrapSprig               
cat                     plural                  shuffle                 trimAll                 wrapWith                
hasPrefix               quote                   snakecase               trimPrefix              
hasSuffix               randAlpha               squote                  trimSuffix              
indent                  randAlphaNum            substr                  trimall                 
initials                randAscii               swapcase                truncSprig              

Sprig Type Conversion http://masterminds.github.io/sprig/conversion.html

atoi                    float64                 int64                   intSprig                

Sprig Version comparison http://masterminds.github.io/sprig/semver.html

semver                  semverCompare           

```
