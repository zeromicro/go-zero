export const enum Method {
  GET = "GET",
  HEAD = "HEAD",
  POST = "POST",
  PUT = "PUT",
  PATCH = "PATCH",
  DELETE = "DELETE",
  CONNECT = "CONNECT",
  OPTIONS = "OPTIONS",
  TRACE = "TRACE",
}

export const enum ContentType {
  ApplicationJson = "application/json",
  ApplicationForm = "application/x-www-form-urlencoded",
}

export const enum GolangType {
  Bool = "bool",
  Uint8 = "uint8",
  Uint16 = "uint16",
  Uint32 = "uint32",
  Uint64 = "uint64",
  Int8 = "int8",
  Int16 = "int16",
  Int32 = "int32",
  Int64 = "int64",
  Float32 = "float32",
  Float64 = "float64",
  String = "string",
  Int = "int",
  Uint = "uint",
}

export const RoutePanelData = {
  GolangTypeOptions: [
    {
      value: GolangType.Bool,
      label: GolangType.Bool,
    },
    {
      value: GolangType.Uint8,
      label: GolangType.Uint8,
    },
    {
      value: GolangType.Uint16,
      label: GolangType.Uint16,
    },
    {
      value: GolangType.Uint32,
      label: GolangType.Uint32,
    },
    {
      value: GolangType.Uint64,
      label: GolangType.Uint64,
    },
    {
      value: GolangType.Int8,
      label: GolangType.Int8,
    },
    {
      value: GolangType.Int16,
      label: GolangType.Int16,
    },
    {
      value: GolangType.Int32,
      label: GolangType.Int32,
    },
    {
      value: GolangType.Int64,
      label: GolangType.Int64,
    },
    {
      value: GolangType.Float32,
      label: GolangType.Float32,
    },
    {
      value: GolangType.Float64,
      label: GolangType.Float64,
    },
    {
      value: GolangType.String,
      label: GolangType.String,
    },
    {
      value: GolangType.Int,
      label: GolangType.Int,
    },
    {
      value: GolangType.Uint,
      label: GolangType.Uint,
    },
  ],
  ContentTypeOptions: [
    {
      value: ContentType.ApplicationJson,
      label: ContentType.ApplicationJson,
    },
    {
      value: ContentType.ApplicationForm,
      label: ContentType.ApplicationForm,
    },
  ],
  MethodOptions: [
    {
      value: Method.GET,
      label: Method.GET,
    },
    {
      value: Method.HEAD,
      label: Method.HEAD,
    },
    {
      value: Method.POST,
      label: Method.POST,
    },
    {
      value: Method.PUT,
      label: Method.PUT,
    },
    {
      value: Method.PATCH,
      label: Method.PATCH,
    },
    {
      value: Method.DELETE,
      label: Method.DELETE,
    },
    {
      value: Method.CONNECT,
      label: Method.CONNECT,
    },
    {
      value: Method.OPTIONS,
      label: Method.OPTIONS,
    },
    {
      value: Method.TRACE,
      label: Method.TRACE,
    },
  ],
};
