package com.xhb.logic.http.packet.{{.packet}};

import com.xhb.core.packet.HttpPacket;
import com.xhb.core.network.HttpRequestClient;
{{.imports}}

{{.doc}}
public class {{.packetName}} extends HttpPacket<{{.responseType}}> {
	{{.paramsDeclaration}}

	public {{.packetName}}({{.params}}{{if .HasRequestBody}}{{.requestType}} request{{end}}) {
		{{if .HasRequestBody}}super(request);{{else}}super(EmptyRequest.instance);{{end}}
		{{if .HasRequestBody}}this.request = request;{{end}}{{.paramsSetter}}
    }

	@Override
    public HttpRequestClient.Method requestMethod() {
        return HttpRequestClient.Method.{{.method}};
    }

	@Override
    public String requestUri() {
        return {{.uri}};
    }
}
