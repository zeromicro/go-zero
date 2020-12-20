parser grammar ApiParser;

options {
    tokenVocab=ApiLexer;
}

api:        syntaxLit body* EOF;

body:       importSpec
            |infoBlock
            |typeBlock
            |serviceBlock
            ;

syntaxLit:      SYNTAX ASSIGN version=SYNTAX_VERSION;
importSpec:     importLit|importLitGroup;
importLit:      IMPORT importPath=IMPORT_PATH;
importLitGroup:     IMPORT '(' (importPath=IMPORT_PATH)* ')';

infoBlock: INFO '(' (key=ID COLON value=STRING_LIT?)* ')';

typeBlock:      typeLit|typeGroup;
typeLit:        TYPE typeSpec;
typeGroup:      TYPE '(' typeSpec* ')';
typeSpec:       typeAlias|typeStruct;
typeAlias:      alias=ID '='? dataType;
typeStruct:     name=ID STRUCT? '{' typeField* '}';
typeField:       name=ID filed?;
filed:      (dataType|innerStruct) tag=RAW_STRING?;
innerStruct:        STRUCT? '{' typeField* '}';
dataType:       pointer
                |mapType
                |arrayType
                |INTERFACE
                ;
mapType:        MAP '[' key=GOTYPE ']' value=dataType;
arrayType:      '['']'lit=dataType;
pointer:        STAR* (GOTYPE|ID);

serviceBlock:       serverMeta? serviceBody;
serverMeta:     ATSERVER '(' (key=ID COLON value=idValue?)* ')';
idValue:        (ID ('/' ID)?)+;
serviceBody:        SERVICE serviceName '{' routes=serviceRoute* '}';
serviceName:        ID ('-' ID)?;
serviceRoute:       routeDoc? (serverMeta|routeHandler) routePath ;
routeDoc:       doc|lineDoc;
doc:        ATDOC '(' (key=ID COLON value=STRING_LIT?)* ')';
lineDoc:        ATDOC STRING_LIT;
routeHandler:       ATHANDLER ID;
routePath:      HTTPMETHOD path request? reply?;
path:      ('/' ':'? ID (('?'|'&'|'=') ID)?)+;
request:       '(' ID ')';
reply:      RETURNS '(' ID ')';