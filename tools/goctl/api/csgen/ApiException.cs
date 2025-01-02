public class ApiException : Exception
{
    public ApiException(string message, Exception? inner=null): base(message, inner)
    {

    }
}
