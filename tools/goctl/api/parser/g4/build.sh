#!/bin/bash

target="antlr-4.9-complete.jar"
antlr="${PWD}/$target"
echo "检测 java 环境"
javaPath=$(which java)

echo "$antlr"
if [ ! -f "$javaPath" ]; then
  echo "java 不存在"
  exit 1
fi

echo "java 环境已支持"
echo "检测 antlr 是否存在"
if [ ! -f "$antlr" ];then
  echo "antlr 不存在，开始下载..."
  curl -O https://www.antlr.org/download/${target}
fi

#export CLASSPATH="/Users/bytedance/keson/workspace/go-zero/tools/goctl/api/parser/g4/antlr-4.9-complete.jar:$CLASSPATH"
export CLASSPATH=".:$antlr:$CLASSPATH"
alias antlr4='java -Xmx500M -cp $antlr:$CLASSPATH org.antlr.v4.Tool'
antlr4 -o ./gen/api -package api -visitor ImportParser.g4  -Dlanguage=Go -no-listener -lib . ./ApiParser.g4
