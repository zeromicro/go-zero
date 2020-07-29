package dartgen

import "text/template"

var funcMap = template.FuncMap{
	"tagGet":          tagGet,
	"isDirectType":    isDirectType,
	"isClassListType": isClassListType,
	"getCoreType":     getCoreType,
	"pathToFuncName":  pathToFuncName,
	"lowCamelCase":    lowCamelCase,
}

const apiFileContent = `import 'dart:io';
import 'dart:convert';
import '../vars/kv.dart';
import '../vars/vars.dart';

/// 发送POST请求.
///
/// data:为你要post的结构体，我们会帮你转换成json字符串;
/// ok函数:请求成功的时候调用，fail函数：请求失败的时候会调用，eventually函数：无论成功失败都会调用
Future apiPost(String path, dynamic data,
    {Map<String, String> header,
    Function(Map<String, dynamic>) ok,
    Function(String) fail,
    Function eventually}) async {
  await _apiRequest('POST', path, data,
      header: header, ok: ok, fail: fail, eventually: eventually);
}

/// 发送GET请求.
///
/// ok函数:请求成功的时候调用，fail函数：请求失败的时候会调用，eventually函数：无论成功失败都会调用
Future apiGet(String path,
    {Map<String, String> header,
    Function(Map<String, dynamic>) ok,
    Function(String) fail,
    Function eventually}) async {
  await _apiRequest('GET', path, null,
      header: header, ok: ok, fail: fail, eventually: eventually);
}

Future _apiRequest(String method, String path, dynamic data,
    {Map<String, String> header,
    Function(Map<String, dynamic>) ok,
    Function(String) fail,
    Function eventually}) async {
  var tokens = await getTokens();
  try {
    var client = HttpClient();
    HttpClientRequest r;
    if (method == 'POST') {
      r = await client.postUrl(Uri.parse('https://' + serverHost + path));
    } else {
      r = await client.getUrl(Uri.parse('https://' + serverHost + path));
    }

    r.headers.set('Content-Type', 'application/json');
    if (tokens != null) {
      r.headers.set('Authorization', tokens.accessToken);
    }
    if (header != null) {
      header.forEach((k, v) {
        r.headers.set(k, v);
      });
    }
    var strData = '';
    if (data != null) {
      strData = jsonEncode(data);
    }
    r.write(strData);
    var rp = await r.close();
    var body = await rp.transform(utf8.decoder).join();
    print('${rp.statusCode} - $path');
    print('-- request --');
    print(strData);
    print('-- response --');
    print('$body \n');
    if (rp.statusCode == 404) {
      if (fail != null) fail('404 not found');
    } else {
      Map<String, dynamic> base = jsonDecode(body);
      if (rp.statusCode == 200) {
        if (base['code'] != 0) {
          if (fail != null) fail(base['desc']);
        } else {
          if (ok != null) ok(base['data']);
        }
      } else if (base['code'] != 0) {
        if (fail != null) fail(base['desc']);
      }
    }
  } catch (e) {
    if (fail != null) fail(e.toString());
  }
  if (eventually != null) eventually();
}
`
const tokensFileContent = `class Tokens {
  /// 用于访问的token, 每次请求都必须带在Header里面
  final String accessToken;
  final int accessExpire;

  /// 用于刷新token
  final String refreshToken;
  final int refreshExpire;
  final int refreshAfter;
  Tokens(
      {this.accessToken,
      this.accessExpire,
      this.refreshToken,
      this.refreshExpire,
      this.refreshAfter});
  factory Tokens.fromJson(Map<String, dynamic> m) {
    return Tokens(
        accessToken: m['access_token'],
        accessExpire: m['access_expire'],
        refreshToken: m['refresh_token'],
        refreshExpire: m['refresh_expire'],
        refreshAfter: m['refresh_after']);
  }
  Map<String, dynamic> toJson() {
    return {
      'access_token': accessToken,
      'access_expire': accessExpire,
      'refresh_token': refreshToken,
      'refresh_expire': refreshExpire,
      'refresh_after': refreshAfter,
    };
  }
}
`
