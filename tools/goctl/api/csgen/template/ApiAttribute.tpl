namespace {{.Namespace}};

[AttributeUsage(AttributeTargets.Property | AttributeTargets.Field, AllowMultiple = false)]
public sealed class HeaderPropertyName : Attribute
{
    public HeaderPropertyName(string name)
    {
        Name = name;
    }

    /// <summary>
    /// The name of the property.
    /// </summary>
    public string Name { get; }
}


[AttributeUsage(AttributeTargets.Property | AttributeTargets.Field, AllowMultiple = false)]
public sealed class PathPropertyName : Attribute
{
    public PathPropertyName(string name)
    {
        Name = name;
    }

    /// <summary>
    /// The name of the property.
    /// </summary>
    public string Name { get; }
}


[AttributeUsage(AttributeTargets.Property | AttributeTargets.Field, AllowMultiple = false)]
public sealed class FormPropertyName : Attribute
{
    public FormPropertyName(string name)
    {
        Name = name;
    }

    /// <summary>
    /// The name of the property.
    /// </summary>
    public string Name { get; }
}