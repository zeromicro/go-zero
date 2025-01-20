export type HttpScheme =  'http' | 'https';
export type HttpMethod = "OPTIONS" | "GET" | "HEAD" | "POST" | "PUT" | "DELETE" | "TRACE" | "CONNECT";
export type HttpQuery =  { [key: string]: string|number|null };
export type HttpHeaders = {[key: string]: string };
export type HttpJsonBody = { [key: string]: any };

export class ApiBaseClient {
	private host: string;
	private port: number;
	private scheme: HttpScheme;
	
	
	constructor(host: string, port: number, scheme: HttpScheme='http') {
	    this.host = host;
		this.port = port;
		this.scheme = scheme;
	}
	
	makeUrl(path: string, query?: HttpQuery): string {
		const head = `${this.scheme}://${this.host}:${this.port}${path}`;
		if (query) {
			const items: Array<string> = [];
			for(let key of Object.keys(query)) {
				items.push(`${key}=${query[key]}`);
			}
			return `${head}?${items.join('&')}`;
		}
		return head;
	}
	
	async request(method: HttpMethod, path: string, query?: HttpQuery, headers?: HttpHeaders, body?: HttpJsonBody): Promise<UniApp.RequestSuccessCallbackResult> {
		const header: HttpHeaders = { };
		if (headers) {
			for(const key of Object.keys(headers)) {
				if (key.toLowerCase() != 'user-agent') {
					header[key] = headers[key];
				}
			}
		}
		
		return new Promise<UniApp.RequestSuccessCallbackResult>((resolve, reject) => {
			uni.request({
				method: method,
				url: this.makeUrl(path, query),
				data: body,
				dataType: 'json',
				header: header,
				success: (res: UniApp.RequestSuccessCallbackResult) => {
					resolve(res);
				},
				fail: (error: UniApp.GeneralCallbackResult) => {
					reject(error);
				}
			})
		});
	}
}