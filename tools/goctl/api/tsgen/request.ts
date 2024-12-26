export type Method =
    | 'get'
    | 'GET'
    | 'delete'
    | 'DELETE'
    | 'head'
    | 'HEAD'
    | 'options'
    | 'OPTIONS'
    | 'post'
    | 'POST'
    | 'put'
    | 'PUT'
    | 'patch'
    | 'PATCH';

/**
 * Parse route parameters for responseType
 */
const reg = /:[a-z|A-Z]+/g;

export function parseParams(url: string): Array<string> {
    const ps = url.match(reg);
    if (!ps) {
        return [];
    }
    return ps.map((k) => k.replace(/:/, ''));
}

/**
 * Generate url and parameters
 * @param url
 * @param params
 */
export function genUrl(url: string, params: any) {
    if (!params) {
        return url;
    }

    const ps = parseParams(url);
    ps.forEach((k) => {
        const reg = new RegExp(`:${k}`);
        url = url.replace(reg, params[k]);
    });

    const path: Array<string> = [];
    for (const key of Object.keys(params)) {
        if (!ps.find((k) => k === key)) {
            path.push(`${key}=${params[key]}`);
        }
    }

    return url + (path.length > 0 ? `?${path.join('&')}` : '');
}

export async function request({
    method,
    url,
    data,
    config = {}
}: {
    method: Method;
    url: string;
    data?: unknown;
    config?: unknown;
}) {
    const response = await fetch(url, {
        method: method.toLocaleUpperCase(),
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json'
        },
        body: data ? JSON.stringify(data) : undefined,
        // @ts-ignore
        ...config
    });

    return response.json();
}

function api<T>(
    method: Method = 'get',
    url: string,
    req: any,
    config?: unknown
): Promise<T> {
    if (url.match(/:/) || method.match(/get|delete/i)) {
        url = genUrl(url, req.params || req.forms);
    }
    method = method.toLocaleLowerCase() as Method;

    switch (method) {
        case 'get':
            return request({method: 'get', url, data: req, config});
        case 'delete':
            return request({method: 'delete', url, data: req, config});
        case 'put':
            return request({method: 'put', url, data: req, config});
        case 'post':
            return request({method: 'post', url, data: req, config});
        case 'patch':
            return request({method: 'patch', url, data: req, config});
        default:
            return request({method: 'post', url, data: req, config});
    }
}

export const webapi = {
    get<T>(url: string, req: unknown, config?: unknown): Promise<T> {
        return api<T>('get', url, req, config);
    },
    delete<T>(url: string, req: unknown, config?: unknown): Promise<T> {
        return api<T>('delete', url, req, config);
    },
    put<T>(url: string, req: unknown, config?: unknown): Promise<T> {
        return api<T>('put', url, req, config);
    },
    post<T>(url: string, req: unknown, config?: unknown): Promise<T> {
        return api<T>('post', url, req, config);
    },
    patch<T>(url: string, req: unknown, config?: unknown): Promise<T> {
        return api<T>('patch', url, req, config);
    }
};

export default webapi
