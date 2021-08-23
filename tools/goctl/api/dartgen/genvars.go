package dartgen

import (
	"io/ioutil"
	"os"
)

const varTemplate = `import 'dart:convert';
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
`

func genVars(dir string) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	if !fileExists(dir + "vars.dart") {
		err = ioutil.WriteFile(dir+"vars.dart", []byte(`const serverHost='demo-crm.xiaoheiban.cn';`), 0o644)
		if err != nil {
			return err
		}
	}

	if !fileExists(dir + "kv.dart") {
		err = ioutil.WriteFile(dir+"kv.dart", []byte(varTemplate), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}
