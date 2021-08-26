lexer grammar ApiLexer;

// Keywords
ATDOC:              '@doc';
ATHANDLER:          '@handler';
INTERFACE:          'interface{}';
ATSERVER:           '@server';

// Whitespace and comments
WS:                 [ \t\r\n\u000C]+ -> channel(HIDDEN);
COMMENT:            '/*' .*? '*/' -> channel(88);
LINE_COMMENT:       '//' ~[\r\n]* -> channel(88);
STRING:             '"' (~["\\] | EscapeSequence)* '"';
RAW_STRING:         '`' (~[`\\\r\n] | EscapeSequence)+ '`';
LINE_VALUE:         ':' [ \t]* (STRING|(~[\r\n"`]*));
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