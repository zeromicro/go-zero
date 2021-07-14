grammar ApiParser;

import ApiLexer;

@lexer::members{
    const COMEMNTS = 88
}

api:            spec*;
spec:           syntaxLit
                |importSpec
                |infoSpec
                |typeSpec
                |serviceSpec
                ;

// syntax
syntaxLit:      {match(p,"syntax")}syntaxToken=ID assign='=' {checkVersion(p)}version=STRING;

// import
importSpec:     importLit|importBlock;
importLit:      {match(p,"import")}importToken=ID importValue ;
importBlock:    {match(p,"import")}importToken=ID '(' importBlockValue+ ')';
importBlockValue:   importValue;
importValue:    {checkImportValue(p)}STRING;

// info
infoSpec:       {match(p,"info")}infoToken=ID lp='(' kvLit+ rp=')';

// type
typeSpec:       typeLit
                |typeBlock;

// eg: type Foo int
typeLit:        {match(p,"type")}typeToken=ID  typeLitBody;
// eg: type (...)
typeBlock:      {match(p,"type")}typeToken=ID lp='(' typeBlockBody* rp=')';
typeLitBody:    typeStruct|typeAlias;
typeBlockBody:  typeBlockStruct|typeBlockAlias;
typeStruct:     {checkKeyword(p)}structName=ID structToken=ID? lbrace='{'  field* rbrace='}';
typeAlias:      {checkKeyword(p)}alias=ID assign='='? dataType;
typeBlockStruct: {checkKeyword(p)}structName=ID structToken=ID? lbrace='{'  field* rbrace='}';
typeBlockAlias: {checkKeyword(p)}alias=ID assign='='? dataType;
field:          {isNormal(p)}? normalField|anonymousFiled ;
normalField:    {checkKeyword(p)}fieldName=ID dataType tag=RAW_STRING?;
anonymousFiled: star='*'? ID;
dataType:       {isInterface(p)}ID
                |mapType
                |arrayType
                |inter='interface{}'
                |time='time.Time'
                |pointerType
                |typeStruct
                ;
pointerType:    star='*' {checkKeyword(p)}ID;
mapType:        {match(p,"map")}mapToken=ID lbrack='[' {checkKey(p)}key=ID rbrack=']' value=dataType;
arrayType:      lbrack='[' rbrack=']' dataType;

// service
serviceSpec:    atServer? serviceApi;
atServer:       ATSERVER lp='(' kvLit+ rp=')';
serviceApi:     {match(p,"service")}serviceToken=ID serviceName lbrace='{' serviceRoute* rbrace='}';
serviceRoute:   atDoc? (atServer|atHandler) route;
atDoc:          ATDOC lp='('? ((kvLit+)|STRING) rp=')'?;
atHandler:      ATHANDLER ID;
route:          {checkHttpMethod(p)}httpMethod=ID path request=body? returnToken=ID? response=replybody?;
body:           lp='(' (ID)? rp=')';
replybody:      lp='(' dataType? rp=')';
// kv
kvLit:          key=ID {checkKeyValue(p)}value=LINE_VALUE;

serviceName:    (ID '-'?)+;
path:           (('/' (ID ('-' ID)*))|('/:' (ID ('-' ID)?)))+;