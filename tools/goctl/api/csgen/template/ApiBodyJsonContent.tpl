using System.Net;
using System.Text.Json;

namespace {{.Namespace}};

public sealed class ApiBodyJsonContent : HttpContent
{
    public MemoryStream Stream { get; private set; }

    private ApiBodyJsonContent()
    {
        Stream = new MemoryStream();
        Headers.Add("Content-Type", "application/json");
    }

    protected override Task SerializeToStreamAsync(Stream stream, TransportContext? context)
    {
        Stream.CopyTo(stream);
        return Task.CompletedTask;
    }

    protected override bool TryComputeLength(out long length)
    {
        length = Stream.Length;
        return true;
    }

    public static ApiBodyJsonContent Create<T>(T t, JsonSerializerOptions? settings=null)
    {
        var content = new ApiBodyJsonContent();
        JsonSerializer.Serialize(content.Stream, t, settings);
        content.Stream.Position = 0;
        return content;
    }
}
