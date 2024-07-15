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
  // array
  SUint8 = "[]uint8",
  SUint16 = "[]uint16",
  SUint32 = "[]uint32",
  SUint64 = "[]uint64",
  SInt8 = "[]int8",
  SInt16 = "[]int16",
  SInt32 = "[]int32",
  SInt64 = "[]int64",
  SFloat32 = "[]float32",
  SFloat64 = "[]float64",
  SString = "[]string",
  SInt = "[]int",
  SUint = "[]uint",
  // any
  Any = "interface{}",
}

export const RoutePanelData = {
  IDPattern: /^[a-zA-Z][\w]*$/gm,
  IDCommaPattern: /^([a-zA-Z](\w)*)+(,([a-zA-Z](\w)*)+)*$/gm,
  EnumCommaPattern: /([^|]*)(\|([^|]+))*[^|]+$/gm,
  PrefixPathPattern: /^\/[\w\/-]*[\w]$/gm,
  PathPattern: /^\/([\w\/-]|(:\w+))*[\w]$/gm,
  IsNumberType: (type: GolangType) => {
    return (
      type === GolangType.Uint8 ||
      type === GolangType.Uint16 ||
      type === GolangType.Uint32 ||
      type === GolangType.Uint64 ||
      type === GolangType.Int8 ||
      type === GolangType.Int16 ||
      type === GolangType.Int32 ||
      type === GolangType.Int64 ||
      type === GolangType.Int ||
      type === GolangType.Uint
    );
  },
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
    {
      value: GolangType.SUint8,
      label: GolangType.SUint8,
    },
    {
      value: GolangType.SUint16,
      label: GolangType.SUint16,
    },
    {
      value: GolangType.SUint32,
      label: GolangType.SUint32,
    },
    {
      value: GolangType.SUint64,
      label: GolangType.SUint64,
    },
    {
      value: GolangType.SInt8,
      label: GolangType.SInt8,
    },
    {
      value: GolangType.SInt16,
      label: GolangType.SInt16,
    },
    {
      value: GolangType.SInt32,
      label: GolangType.SInt32,
    },
    {
      value: GolangType.SInt64,
      label: GolangType.SInt64,
    },
    {
      value: GolangType.SFloat32,
      label: GolangType.SFloat32,
    },
    {
      value: GolangType.SFloat64,
      label: GolangType.SFloat64,
    },
    {
      value: GolangType.SString,
      label: GolangType.SString,
    },
    {
      value: GolangType.SInt,
      label: GolangType.SInt,
    },
    {
      value: GolangType.SUint,
      label: GolangType.SUint,
    },
    {
      value: GolangType.Any,
      label: GolangType.Any,
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
