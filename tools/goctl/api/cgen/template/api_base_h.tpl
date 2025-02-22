#ifndef GO_ZERO_API_BASE_H
#define GO_ZERO_API_BASE_H

#include <stddef.h>
#include <stdbool.h>
#include <curl/curl.h>

#define HTTP_METHOD_GET 1
#define HTTP_METHOD_PUT 2
#define HTTP_METHOD_POST 3

typedef int http_method_t;
typedef struct curl_slist curl_slist_t;

typedef struct __array_t {
    void* items;
    size_t count;
} array_t;

typedef struct __{{.ClientName}}_t {
    bool is_secure;
    char* host;
    short port;
    long timeout_s;
}{{.ClientName}}_t;

typedef struct __base_request_t {
    http_method_t method;
    curl_slist_t *headers;
    const char* path;
    char** query;
    const char* body;
    size_t body_size;
} base_request_t;

typedef struct __base_response_t {
    char* body;
    size_t body_size;
    char* headers;
    size_t headers_size;
} base_response_t;

bool base_client_request({{.ClientName}}_t* client, base_request_t* request, base_response_t* response);
const char* base_last_error();
void base_response_free(base_response_t* response);
void base_client_clearup();

char* replace_substr(char* source, const char* substr, const char* replacement);

#endif
