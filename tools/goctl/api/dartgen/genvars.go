package dartgen

import (
	"io/ioutil"
	"os"

	"github.com/tal-tech/go-zero/core/logx"
)

func genVars(dir string) error {
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		logx.Error(e)
		return e
	}

	if !fileExists(dir + "vars.dart") {
		e = ioutil.WriteFile(dir+"vars.dart", []byte(`const serverHost='demo-crm.xiaoheiban.cn';`), 0644)
		if e != nil {
			logx.Error(e)
			return e
		}
	}

	if !fileExists(dir + "kv.dart") {
		e = ioutil.WriteFile(dir+"kv.dart", []byte(`import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../data/tokens.dart';

/// 保存tokens到本地
///
/// 传入null则删除本地tokens
/// 返回：true：设置成功  false：设置失败
Future<bool> setTokens(Tokens tokens) async {
  var sp = await SharedPreferences.getInstance();
  if (tokens == null) {
    sp.remove('tokens');
    return true;
  }
  return await sp.setString('tokens', jsonEncode(tokens.toJson()));
}

/// 获取本地存储的tokens
///
/// 如果没有，则返回null
Future<Tokens> getTokens() async {
  try {
    var sp = await SharedPreferences.getInstance();
    var str = sp.getString('tokens');
    if (str.isEmpty) {
      return null;
    }
    return Tokens.fromJson(jsonDecode(str));
  } catch (e) {
    print(e);
    return null;
  }
}
`), 0644)
		if e != nil {
			logx.Error(e)
			return e
		}
	}
	return nil
}
