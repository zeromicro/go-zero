lexer grammar ApiLexer;

// Keywords
GOTYPE:
                    BOOL
                    |UINT8
                    |UINT16
                    |UINT32
                    |UINT64
                    |INT8
                    |INT16
                    |INT32
                    |INT64
                    |FLOAT32
                    |FLOAT64
                    |COMPLEX64
                    |COMPLEX128
                    |STRING
                    |INT
                    |UINT
                    |UINTPTR
                    |BYTE
                    |RUNE
                    |TIME
                    ;


SYNTAX:             'syntax';
INFO:               'info';
MAP:                'map';
STRUCT:             'struct';
INTERFACE:          'interface{}';
TYPE:               'type';
ATSERVER:           '@server';
ATDOC:              '@doc';
ATHANDLER:          '@handler';
SERVICE:            'service';
RETURNS:            'returns';
IMPORT:             'import';

HTTPMETHOD:         GET
                    |HEAD
                    |POST
                    |PUT
                    |PATCH
                    |DELETE
                    |CONNECT
                    |OPTIONS
                    |TRACE
                    ;



// separators
LPAREN:             '(';
RPAREN:             ')';
LBRACE:             '{';
RBRACE:             '}';
LBRACK:             '[';
RBRACK:             ']';
COMMA:              ',';
DOT:                '.';
SLASH:              '/';
QUESTION:           '?';
BITAND:             '&';

// Operators
ASSIGN:             '=';
SUB:                '-';
COLON:              ':';
STAR:                '*';

// Whitespace and comments
WS:                 [ \t\r\n\u000C]+ -> channel(HIDDEN);
COMMENT:            '/*' .*? '*/'    -> channel(HIDDEN);
LINE_COMMENT:       '//' ~[\r\n]*    -> channel(HIDDEN);

// Literals
SYNTAX_VERSION:     '"' 'v'[1-9][0-9]* '"';
IMPORT_PATH:        '"' '/'? ID ('/' ID)* '.api' '"';
STRING_LIT:         ('"' (~["\\] | EscapeSequence)* '"');
RAW_STRING:         '`' (~[`\\\r\n] | EscapeSequence)* '`';

ID:         Letter LetterOrDigit*;


fragment ExponentPart
    : [eE] [+-]? Digits
    ;

fragment EscapeSequence
    : '\\' [btnfr"'\\]
    | '\\' ([0-3]? [0-7])? [0-7]
    | '\\' 'u'+ HexDigit HexDigit HexDigit HexDigit
    ;
fragment HexDigits
    : HexDigit ((HexDigit | '_')* HexDigit)?
    ;
fragment HexDigit
    : [0-9a-fA-F]
    ;
fragment Digits
    : [0-9] ([0-9_]* [0-9])?
    ;

fragment LetterOrDigit
    : Letter
    | [0-9]
    ;
fragment Letter
    : [a-zA-Z$_] // these are the "java letters" below 0x7F
    | ~[\u0000-\u007F\uD800-\uDBFF] // covers all characters above 0x7F which are not a surrogate
    | [\uD800-\uDBFF] [\uDC00-\uDFFF] // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
    ;

fragment BOOL:               'bool';
fragment UINT8:              'uint8';
fragment UINT16:             'uint16';
fragment UINT32:             'uint32';
fragment UINT64:             'uint64';
fragment INT8:               'int8';
fragment INT16:              'int16';
fragment INT32:              'int32';
fragment INT64:              'int64';
fragment FLOAT32:            'float32';
fragment FLOAT64:            'float64';
fragment COMPLEX64:          'complex64';
fragment COMPLEX128:         'complex128';
fragment STRING:             'string';
fragment INT:                'int';
fragment UINT:               'uint';
fragment UINTPTR:            'uintptr';
fragment BYTE:               'byte';
fragment RUNE:               'rune';
fragment TIME:               'time.Time';
fragment GET:                'get';
fragment HEAD:               'head';
fragment POST:               'post';
fragment PUT:                'put';
fragment PATCH:              'patch';
fragment DELETE:             'delete';
fragment CONNECT:            'connect';
fragment OPTIONS:            'options';
fragment TRACE:              'trace';