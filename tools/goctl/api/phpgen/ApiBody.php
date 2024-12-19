<?php


/**
 * 提供 go-zero api 之外自定义 http 内容体的方式。
 * 可以作为可选参数替代 api 请求的内容体。
 * 因为 go-zero api 内容体都是 json 所以 $data 必须是可 json 序列化的。
 */
class ApiBody
{
    private $data;

    public function __construct($data)
    {
        $this->data = $data;
    }

    public function toJsonString()
    {
        return json_encode($this->data, JSON_UNESCAPED_UNICODE);
    }

    public function toAssocArray()
    {
        return $this->data;
    }
}
