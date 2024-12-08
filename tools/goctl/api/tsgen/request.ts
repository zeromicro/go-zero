export type Method =
    | "get"
    | "GET"
    | "delete"
    | "DELETE"
    | "head"
    | "HEAD"
    | "options"
    | "OPTIONS"
    | "post"
    | "POST"
    | "put"
    | "PUT"
    | "patch"
    | "PATCH";

export type QueryParams = {
    [key: string | symbol | number]: string|number;
} | null;

/**
 * Parse route parameters for responseType
 */
const reg = /:[a-z|A-Z]+/g;

export function parseParams(url: string): Array<string> {
    const ps = url.match(reg);
    if (!ps) {
        return [];
    }
    return ps.map((k) => k.replace(/:/, ""));
}

/**
 * Generate url and parameters
 * @param url
 * @param params
 */
export function genUrl(url: string, params: QueryParams) {
    if (!params) {
        return url;
    }

    const ps = parseParams(url);
    ps.forEach((k) => {
        const reg = new RegExp(`:${k}`);
        url = url.replace(reg, params[k].toString());
    });

    const path: Array<string> = [];
    for (const key of Object.keys(params)) {
        if (!ps.find((k) => k === key)) {
            path.push(`${key}=${params[key]}`);
        }
    }

    return url + (path.length > 0 ? `?${path.join("&")}` : "");
}

export async function request(
    method: Method,
    url: string,
    config?: RequestInit
) {
    if (config?.body && /get|head/i.test(method)) {
        throw new Error(
            "Request with GET/HEAD method cannot have body. *.api service use other method, example: POST or PUT."
        );
    }
    const response = await fetch(url, {
        method: method.toLocaleUpperCase(),
        credentials: "include",
        ...config,
        headers: {
            "Content-Type": "application/json",
            ...config?.headers,
        },
    });

    return response.json();
}

function api<T>(
    method: Method = "get",
    url: string,
    params?: QueryParams,
    config?: RequestInit
): Promise<T> {
    if (params) {
        url = genUrl(url, params);
    }
    method = method.toLocaleLowerCase() as Method;

    switch (method) {
        case "get":
            return request("get", url, config);
        case "delete":
            return request("delete", url, config);
        case "put":
            return request("put", url, config);
        case "post":
            return request("post", url, config);
        case "patch":
            return request("patch", url, config);
        default:
            return request("post", url, config);
    }
}

export const webapi = {
    get<T>(url: string, params?: QueryParams, config?: RequestInit): Promise<T> {
        return api<T>("get", url, params, config);
    },
    delete<T>(
        url: string,
        params?: QueryParams,
        config?: RequestInit
    ): Promise<T> {
        return api<T>("delete", url, params, config);
    },
    put<T>(url: string, params?: QueryParams, config?: RequestInit): Promise<T> {
        return api<T>("put", url, params, config);
    },
    post<T>(url: string, params?: QueryParams, config?: RequestInit): Promise<T> {
        return api<T>("post", url, params, config);
    },
    patch<T>(
        url: string,
        params?: QueryParams,
        config?: RequestInit
    ): Promise<T> {
        return api<T>("patch", url, params, config);
    },
};

export default webapi;
