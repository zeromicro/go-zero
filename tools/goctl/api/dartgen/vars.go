package dartgen

import "text/template"

var funcMap = template.FuncMap{
	"appendNullCoalescing":            appendNullCoalescing,
	"appendDefaultEmptyValue":         appendDefaultEmptyValue,
	"extractPositionalParamsFromPath": extractPositionalParamsFromPath,
	"getBaseName":                     getBaseName,
	"getCoreType":                     getCoreType,
	"getPropertyFromMember":           getPropertyFromMember,
	"hasUrlPathParams":                hasUrlPathParams,
	"isAtomicListType":                isAtomicListType,
	"isAtomicType":                    isAtomicType,
	"isDirectType":                    isDirectType,
	"isClassListType":                 isClassListType,
	"isListItemsNullable":             isListItemsNullable,
	"isMapType":                       isMapType,
	"isNullableType":                  isNullableType,
	"isNumberType":                    isNumberType,
	"lowCamelCase":                    lowCamelCase,
	"makeDartRequestUrlPath":          makeDartRequestUrlPath,
	"normalizeHandlerName":            normalizeHandlerName,
}

const (
	apiFileContent = `import 'dart:io';
import 'dart:convert';
import '../vars/kv.dart';
import '../vars/vars.dart';

/// Send GET request.
///
/// ok: the function that will be called on success.
/// fail：the fuction that will be called on failure.
/// eventually：the function that will be called regardless of success or failure.
Future apiGet(String path,
    {Map<String, String> header,
    Function(Map<String, dynamic>) ok,
    Function(String) fail,
    Function eventually}) async {
  await _apiRequest('GET', path, null,
      header: header, ok: ok, fail: fail, eventually: eventually);
}

/// Send POST request.
///
/// data: the data to post, it will be marshaled to json automatically.
/// ok: the function that will be called on success.
/// fail：the fuction that will be called on failure.
/// eventually：the function that will be called regardless of success or failure.
Future apiPost(String path, dynamic data,
    {Map<String, String> header,
    Function(Map<String, dynamic>) ok,
    Function(String) fail,
    Function eventually}) async {
  await _apiRequest('POST', path, data,
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
      r = await client.postUrl(Uri.parse(serverHost + path));
    } else {
      r = await client.getUrl(Uri.parse(serverHost + path));
    }

    var strData = '';
    if (data != null) {
      strData = jsonEncode(data);
    }

    if (method == 'POST') {
      r.headers.set('Content-Type', 'application/json; charset=utf-8');
      r.headers.set('Content-Length', utf8.encode(strData).length);
    }

    if (tokens != null) {
      r.headers.set('Authorization', tokens.accessToken);
    }
    if (header != null) {
      header.forEach((k, v) {
        r.headers.set(k, v);
      });
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

	apiFileContentV2 = `import 'dart:io';
	import 'dart:convert';
	import '../vars/kv.dart';
	import '../vars/vars.dart';

	/// send request with post method
	///
	/// data: any request class that will be converted to json automatically
	/// ok: is called when request succeeds
	/// fail: is called when request fails
	/// eventually: is always called until the nearby functions returns
	Future apiPost(String path, dynamic data,
			{Map<String, String>? header,
			Function(Map<String, dynamic>)? ok,
			Function(String)? fail,
			Function? eventually}) async {
		await _apiRequest('POST', path, data,
				header: header, ok: ok, fail: fail, eventually: eventually);
	}

	/// send request with get method
	///
	/// ok: is called when request succeeds
	/// fail: is called when request fails
	/// eventually: is always called until the nearby functions returns
	Future apiGet(String path,
			{Map<String, String>? header,
			Function(Map<String, dynamic>)? ok,
			Function(String)? fail,
			Function? eventually}) async {
		await _apiRequest('GET', path, null,
				header: header, ok: ok, fail: fail, eventually: eventually);
	}

	Future _apiRequest(String method, String path, dynamic data,
			{Map<String, String>? header,
			Function(Map<String, dynamic>)? ok,
			Function(String)? fail,
			Function? eventually}) async {
		var tokens = await getTokens();
		try {
			var client = HttpClient();
			HttpClientRequest r;
			if (method == 'POST') {
				r = await client.postUrl(Uri.parse(serverHost + path));
			} else {
				r = await client.getUrl(Uri.parse(serverHost + path));
			}

      var strData = '';
			if (data != null) {
				strData = jsonEncode(data);
			}
			if (method == 'POST') {
        r.headers.set('Content-Type', 'application/json; charset=utf-8');
        r.headers.set('Content-Length', utf8.encode(strData).length);
      }
			if (tokens != null) {
				r.headers.set('Authorization', tokens.accessToken);
			}
			if (header != null) {
				header.forEach((k, v) {
					r.headers.set(k, v);
				});
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
	}`

	tokensFileContent = `class Tokens {
  /// the token used to access, it must be carried in the header of each request
  final String accessToken;
  final int accessExpire;

  /// the token used to refresh
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

	tokensFileContentV2 = `class Tokens {
  /// the token used to access, it must be carried in the header of each request
  final String accessToken;
  final int accessExpire;

  /// the token used to refresh
  final String refreshToken;
  final int refreshExpire;
  final int refreshAfter;
  Tokens({
		required this.accessToken,
		required this.accessExpire,
		required this.refreshToken,
		required this.refreshExpire,
		required this.refreshAfter
	});
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
)
