# 请先安装 antlr4-tools  https://github.com/antlr/antlr4/blob/master/doc/getting-started.md
# pip install antlr4-tools

antlr4 -Dlanguage=Go ApiParser.g4 -visitor -no-listener -o ./gen/api