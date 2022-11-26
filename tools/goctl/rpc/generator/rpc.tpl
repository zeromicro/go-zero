syntax = "proto3";

package {{.package}};
option go_package="./{{.package}}";

message Empty {}

message IDReq {
  uint64 id = 1;
}

message IDsReq {
  repeated uint64 ids = 1;
}

message UUIDReq {
  string uuid = 1;
}

message BaseResp {
  string msg = 1;
}

message PageInfoReq {
  uint64 page = 1;
  uint64 page_size = 2;
}

service {{.serviceName}} {
  rpc initDatabase (Empty) returns (BaseResp);
}
