<?php

class ApiBaseClient
{
    private $host;
    private $port;
    private $scheme;

    public function __construct(
        $host,
        $port,
        $scheme
    ) {
        $this->host = $host;
        $this->port = $port;
        $this->scheme = $scheme;
    }

    public function getHost()
    {
        return $this->host;
    }
    public function getPort()
    {
        return $this->port;
    }

    public function getAddress()
    {
        return "$this->scheme://$this->host:$this->port";
    }

    protected function request(
        $path, // 请求路径
        $method, // 请求方法 get post delete put
        $params, // 请求路径参数
        $query,  // 请求字符串
        $headers, // 头部字段
        $body // 内容体
    ) {
        $address = $this->getAddress();
        if (!$headers) {
            $headers = [];
        } else {
            $headers = $headers->toAssocArray();
        }

        // path
        if ($params) {
            $path = self::replacePathParams($path, $params->toAssocArray());
        }
        $url = "$address$path";

        // query
        if ($query) {
            $queryString = $query->toQueryString();
            $url = "$url?$queryString";
        }

        $ch = curl_init();

        // 2. 设置请求选项, 包括具体的url
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_1); // HTTP/1.1
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 4); // 在发起连接前等待的时间，如果设置为0，则无限等待。
        curl_setopt($ch, CURLOPT_TIMEOUT, 20); // 设置cURL允许执行的最长秒数
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_HEADER, 1);

        // POST
        if ($method == 'post') {
            curl_setopt($ch, CURLOPT_POST, true);
        } else {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, strtoupper($method));
        }

        // body
        if ($body) {
            $headers['Content-Type'] = 'application/json';
            curl_setopt($ch, CURLOPT_POSTFIELDS, $body->toJsonString());
        }

        // header
        if (!empty($headers)) {
            $header = [];
            foreach ($headers as $k => $v) {
                $header[] = "$k: $v";
            }
            curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
        }

        // 3. 执行一个cURL会话并且获取响应消息。
        curl_setopt($ch, CURLINFO_HEADER_OUT, true);
        $response = curl_exec($ch);
        if ($response === false) {
            $error = curl_error($ch);
            throw new ApiException("curl exec failed." . var_export($error, true), -1);
        }

        // 4. 释放cURL句柄,关闭一个cURL会话。
        curl_close($ch);

        return self::parseResponse($response);
    }

    private static function parseResponse($response)
    {
        $statusEnd = strpos($response, "\r\n");
        $status = substr($response, 0, $statusEnd);
        $status = explode(' ', $status, 3);
        $statusCode = intval($status[1] ?? 0);

        if ($statusCode < 200 || $statusCode > 299) {
            throw new ApiException("response failed.", -2, null, $response);
        }

        $headerEnd = strpos($response, "\r\n\r\n");
        $header = substr($response, $statusEnd + 2, $headerEnd - ($statusEnd + 2));
        $header = explode("\r\n", $header);
        $headers = [];
        foreach ($header as $row) {
            $kw = explode(':', $row, 2);
            $headers[strtolower($kw[0])] = $kw[1] ?? null;
        }

        $body = json_decode(substr($response, $headerEnd + 4), true);

        return [
            'status' => $status[2],
            'statusCode' => $statusCode,
            'headers' => $headers,
            'body' => $body,
        ];
    }

    // 填入路径参数
    private static function replacePathParams($path, $kw)
    {
        $map = [];
        foreach ($kw as $k => $v) {
            $map[":$k"] = $v;
        }
        $path = str_replace(array_keys($map), $map, $path);
        return $path;
    }
}
