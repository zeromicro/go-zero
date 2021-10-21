#!/bin/bash

target="antlr-4.9-complete.jar"
antlr="${PWD}/$target"
echo "check env JAVA"
javaPath=$(which java)

echo "$antlr"
if [ ! -f "$javaPath" ]; then
  echo "JAVA is not install"
  exit 1
fi

echo "JAVA installed"
echo "check library ANTLR"
if [ ! -f "$antlr" ];then
  echo "ANTLR not exists, start to download ..."
  curl -O https://www.antlr.org/download/${target}
fi

export CLASSPATH=".:$antlr:$CLASSPATH"
alias antlr4='java -Xmx500M -cp $antlr:$CLASSPATH org.antlr.v4.Tool'
antlr4 -o ./gen/api -package api -visitor  -Dlanguage=Go -no-listener -lib . ./ApiParser.g4
