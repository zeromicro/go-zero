#include "base.h"

static bool is_curl_inited = false;
static char last_error[1024];

size_t base_client_read_body(void* contents, size_t size, size_t nmemb, void* userp) {
    base_response_t* response = (base_response_t*)userp;
    size_t real_size = size * nmemb;
    void* memory = realloc(response->body, response->body_size + real_size + 1);
    if (memory == NULL) {
        return 0;
    }
    response->body = memory;

    memcpy(response->body + response->body_size, contents, real_size);
    response->body_size += real_size;
    response->body[response->body_size] = '\0';
    return real_size;
}

size_t base_client_read_headers(void* contents, size_t size, size_t nmemb, void* userp) {
    base_response_t* response = (base_response_t*)userp;
    size_t real_size = size * nmemb;
    void* memory = realloc(response->headers, response->headers_size + real_size + 1);
    if (memory == NULL) {
        return 0;
    }
    response->headers = memory;

    memcpy(response->headers + response->headers_size, contents, real_size);
    response->headers_size += real_size;
    response->headers[response->headers_size] = '\0';
    return real_size;
}

void base_client_clearup() {
    if (is_curl_inited) {
        curl_global_cleanup();
    }
}

const char* base_last_error() {
    return last_error;
}

bool base_client_request({{.ClientName}}_t* client, base_request_t* request, base_response_t* response) {
    if (!is_curl_inited) {
        curl_global_init(CURL_GLOBAL_ALL);
        is_curl_inited = true;
    }

    CURLU *handle = curl_easy_init();
    if (NULL == handle) {
        sprintf(last_error, "curl_easy_init init failed.");
        return false;
    }

    bool status = true;

    // setting
    curl_easy_setopt(handle, CURLOPT_TIMEOUT, client->timeout_s);

    // url
    char* url = malloc(8192);
    if (NULL == url) {
        sprintf(last_error, "malloc url failed.");
        status = false;
        goto exit_free;
    }

    sprintf(url, "%s://%s:%d%s", client->is_secure ? "https" : "http", client->host, client->port, request->path ? request->path : "/");
    curl_easy_setopt(handle, CURLOPT_URL, url);
    free(url);

    // method
    switch (request->method) {
    case HTTP_METHOD_GET:
        curl_easy_setopt(handle, CURLOPT_CUSTOMREQUEST, "GET");
        if (request->body) {
            sprintf(last_error, "http method 'get' not supported request body.");
            status = false;
            goto exit_free;
        }
        break;
    case HTTP_METHOD_PUT:
        curl_easy_setopt(handle, CURLOPT_CUSTOMREQUEST, "PUT");
        curl_easy_setopt(handle, CURLOPT_POSTFIELDS, request->body);
        break;
    case HTTP_METHOD_POST:
        curl_easy_setopt(handle, CURLOPT_CUSTOMREQUEST, "POST");
        curl_easy_setopt(handle, CURLOPT_POSTFIELDS, request->body);
        break;
    default:
        sprintf(last_error, "http method(%d) not supported.", request->method);
        status = false;
        goto exit_free;
    }

    // request headers
    if (request->headers) {
        curl_easy_setopt(handle, CURLOPT_HTTPHEADER, request->headers);
    }

    // response
    response->body = malloc(1);
    if (NULL == response->body) {
        sprintf(last_error, "malloc response body failed.");
        status = false;
        goto exit_free;
    }
    response->body[0] = '\0';
    response->body_size = 0;

    response->headers = malloc(1);
    if (NULL == response->headers) {
        sprintf(last_error, "malloc response headers failed.");
        status = false;
        goto exit_free;
    }
    response->headers[0] = '\0';
    response->headers_size = 0;

    curl_easy_setopt(handle, CURLOPT_WRITEFUNCTION, &base_client_read_body);
    curl_easy_setopt(handle, CURLOPT_WRITEDATA, response);

    curl_easy_setopt(handle, CURLOPT_HEADERFUNCTION, &base_client_read_headers);
    curl_easy_setopt(handle, CURLOPT_HEADERDATA, response);

    if (client->is_secure) {
        curl_easy_setopt(handle, CURLOPT_SSL_VERIFYPEER, 0L);
        curl_easy_setopt(handle, CURLOPT_SSL_VERIFYHOST, 1L);
    }

    // 
    CURLcode res = curl_easy_perform(handle);
    if (res != CURLE_OK) {
        sprintf(last_error, "curl_easy_perform failed(%d): %s.", res, curl_easy_strerror(res));
        status = false;
    }

exit_free:
    curl_easy_cleanup(handle);
    return status;
}

void base_response_free(base_response_t* response) {
    if (response && response->body) {
        free(response->body);
        response->body = NULL;
        response->body_size = 0;
    }
}

char* replace_substr(char* source, const char* substr, const char* replacement) {
    char* p = strstr(source, substr);
    if (p) {
        size_t sl = strlen(source);
        size_t ssl = strlen(substr);
        size_t rl = strlen(replacement);
        size_t ssi = p - source;
        size_t len = sl - ssl + rl + 1;
        char* result = (char*)malloc(len);
        if (!result) {
            return NULL;
        }
        strncpy(result, source, ssi);
        strncpy(result + ssi, replacement, rl);
        strncpy(result + ssi + rl, p + ssl, len - ssi - rl);
        return result;
    }
    return source;
}