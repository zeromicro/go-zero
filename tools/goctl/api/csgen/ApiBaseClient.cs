using System.Net.Http.Json;
using System.Reflection;
using System.Web;

public abstract class ApiBaseClient
{
    protected readonly HttpClient httpClient = new HttpClient();


    public ApiBaseClient(string host = "localhost", short port = 8080, string scheme = "http")
    {
        httpClient.BaseAddress = new Uri($"{scheme}://{host}:{port}");
    }

    public async Task<R> CallResultAsync<R>(HttpMethod method, string path, CancellationToken cancellationToken) where R : new()
    {
        var response = await CallAsync(HttpMethod.Post, path, cancellationToken);
        return await ParseResponseAsync<R>(response);
    }

    public async Task<HttpResponseMessage> CallAsync(HttpMethod method, string path, CancellationToken cancellationToken)
    {
        using var request = new HttpRequestMessage(method, path);
        return await httpClient.SendAsync(request, cancellationToken);
    }


    public async Task<R> RequestResultAsync<T, R>(HttpMethod method, string path, T? param, CancellationToken cancellationToken) where R : new()
    {
        var response = await RequestAsync(HttpMethod.Post, path, param, cancellationToken);
        return await ParseResponseAsync<R>(response);
    }

    public async Task<HttpResponseMessage> RequestAsync<T>(HttpMethod method, string path, T? param, CancellationToken cancellationToken)
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
                    request.Headers.Add(hpn.Name, pv?.ToString() ?? "");
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
                request.Content = JsonContent.Create(param);
            }
        }

        return await httpClient.SendAsync(request, cancellationToken);
    }

    public static async Task<R> ParseResponseAsync<R>(HttpResponseMessage response) where R : new()
    {
        R result = await response.Content.ReadFromJsonAsync<R>() ?? new R();
        foreach (var p in typeof(R).GetProperties())
        {
            // 头部
            var hpn = p.GetCustomAttribute<HeaderPropertyName>();
            if (hpn != null && response.Headers.Contains(hpn.Name))
            {
                // TODO 这里应该有类型问题
                var v = response.Headers.GetValues(hpn.Name);
                p.SetValue(result, v);
                continue;
            }
        }
        return result;
    }
}
