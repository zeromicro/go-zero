package dartgen

import (
	"fmt"
	"os"
)

const (
	varTemplate = `import 'dart:convert';
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
    if (str == null || str.isEmpty) {
      return null;
    }
    return Tokens.fromJson(jsonDecode(str));
  } catch (e) {
    print(e);
    return null;
  }
}
`

	varTemplateV2 = `import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../data/tokens.dart';

const String _tokenKey = 'tokens';

/// Saves tokens
Future<bool> setTokens(Tokens tokens) async {
  var sp = await SharedPreferences.getInstance();
  return await sp.setString(_tokenKey, jsonEncode(tokens.toJson()));
}

/// remove tokens
Future<bool> removeTokens() async {
  var sp = await SharedPreferences.getInstance();
  return sp.remove(_tokenKey);
}

/// Reads tokens
Future<Tokens?> getTokens() async {
  try {
    var sp = await SharedPreferences.getInstance();
    var str = sp.getString('tokens');
    if (str == null || str.isEmpty) {
      return null;
    }
    return Tokens.fromJson(jsonDecode(str));
  } catch (e) {
    print(e);
    return null;
  }
}`
)

func genVars(dir string, isLegacy bool, scheme string, hostname string) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	if !fileExists(dir + "vars.dart") {
		err = os.WriteFile(dir+"vars.dart", []byte(fmt.Sprintf(`const serverHost='%s://%s';`, scheme, hostname)), 0o644)
		if err != nil {
			return err
		}
	}

	if !fileExists(dir + "kv.dart") {
		tpl := varTemplateV2
		if isLegacy {
			tpl = varTemplate
		}
		err = os.WriteFile(dir+"kv.dart", []byte(tpl), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}
