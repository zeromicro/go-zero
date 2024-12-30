using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Reflection;
using System.Web;

public abstract class ApiBaseClient
{
    protected readonly HttpClient httpClient = new HttpClient();


    public ApiBaseClient(string host, short port, string scheme)
    {
        httpClient.BaseAddress = new Uri($"{scheme}://{host}:{port}");
    }

    protected async Task<R> CallResultAsync<R>(HttpMethod method, string path, CancellationToken cancellationToken, HttpContent? body) where R : new()
    {
        var response = await CallAsync(HttpMethod.Post, path, cancellationToken, body);
        return await ParseResponseAsync<R>(response, cancellationToken);
    }

    protected async Task<HttpResponseMessage> CallAsync(HttpMethod method, string path, CancellationToken cancellationToken, HttpContent? body)
    {
        using var request = new HttpRequestMessage(method, path);
        if (body != null)
        {
            request.Content = body;
        }
        return await httpClient.SendAsync(request, cancellationToken);
    }


    protected async Task<R> RequestResultAsync<T, R>(HttpMethod method, string path, T? param, CancellationToken cancellationToken, HttpContent? body) where R : new()
    {
        var response = await RequestAsync(HttpMethod.Post, path, param, cancellationToken, body);
        return await ParseResponseAsync<R>(response, cancellationToken);
    }

    protected async Task<HttpResponseMessage> RequestAsync<T>(HttpMethod method, string path, T? param, CancellationToken cancellationToken, HttpContent? body)
    {
        using var request = new HttpRequestMessage();
        request.Method = method;
        var query = HttpUtility.ParseQueryString(string.Empty);

        if (param != null)
        {
            foreach (var p in param.GetType().GetProperties())
            {
                var pv = p.GetValue(param);

                // 头部
                var hpn = p.GetCustomAttribute<HeaderPropertyName>();
                if (hpn != null)
                {
                    var hv = pv?.ToString() ?? "";
                    if (hpn.Name.ToLower() == "user-agent")
                    {
                        request.Headers.UserAgent.Clear();
                        if (!string.IsNullOrEmpty(hv))
                        {
                            request.Headers.UserAgent.Add(new ProductInfoHeaderValue(hv));
                        }
                    }
                    else
                    {
                        request.Headers.Add(hpn.Name, hv);
                    }
                    continue;
                }

                // 请求参数
                var fpn = p.GetCustomAttribute<FormPropertyName>();
                if (fpn != null)
                {
                    query.Add(fpn.Name, pv?.ToString() ?? "");
                    continue;
                }

                // 路径参数
                var ppn = p.GetCustomAttribute<PathPropertyName>();
                if (ppn != null)
                {
                    path.Replace($":{ppn.Name}", pv?.ToString() ?? "");
                    continue;
                }
            }

            // 请求链接
            request.RequestUri = new Uri(httpClient.BaseAddress!, query.Count > 0 ? $"{path}?{query}" : path);

            // JSON 内容
            if (HttpMethod.Get != method)
            {
                request.Content = body ?? ApiBodyJsonContent.Create(param);
            }
        }

        return await httpClient.SendAsync(request, cancellationToken);
    }

    protected static async Task<R> ParseResponseAsync<R>(HttpResponseMessage response, CancellationToken cancellationToken) where R : new()
    {
        R result = await response.Content.ReadFromJsonAsync<R>(cancellationToken) ?? new R();
        foreach (var p in typeof(R).GetProperties())
        {
            // 头部
            var hpn = p.GetCustomAttribute<HeaderPropertyName>();
            if (hpn != null && response.Headers.Contains(hpn.Name))
            {
                var v = response.Headers.GetValues(hpn.Name);
                if (v != null && v.Count() > 0)
                {
                    p.SetValue(result, v.First());
                }
                continue;
            }
        }
        return result;
    }
}
