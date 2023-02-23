syntax = "proto3";

package {{.package}};
option go_package="./{{.package}}";

// base message
message Empty {}

message IDReq {
  uint64 id = 1;
}

message IDsReq {
  repeated uint64 ids = 1;
}

message UUIDsReq {
  repeated string ids = 1;
}

message UUIDReq {
  string id = 1;
}

message BaseResp {
  string msg = 1;
}

message PageInfoReq {
  uint64 page = 1;
  uint64 page_size = 2;
}

message BaseIDResp {
  uint64 id = 1;
  string msg = 2;
}

message BaseUUIDResp {
  string id = 1;
  string msg = 2;
}


service {{.serviceName}} {
  // group: base
  rpc initDatabase (Empty) returns (BaseResp);
}
