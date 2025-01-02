<?php

class ApiException extends \Exception
{
    private $responseContent;

    public function __construct($message, $code, $previous = null, $responseContent = null)
    {
        $this->responseContent = $responseContent;
        parent::__construct($message, $code, $previous);
    }

    public function getResponseContent() {
        return $this->responseContent;
    }
}
