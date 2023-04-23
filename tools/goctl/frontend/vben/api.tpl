import { defHttp } from '/@/utils/http/axios';
import { ErrorMessageMode } from '/#/axios';
import { BaseDataResp, BaseListReq, BaseResp, Base{{if .useUUID}}UU{{end}}IDsReq, Base{{if .useUUID}}UU{{end}}IDReq } from '/@/api/model/baseModel';
import { {{.modelName}}Info, {{.modelName}}ListResp } from './model/{{.modelNameLowerCamel}}Model';

enum Api {
  Create{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/create',
  Update{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/update',
  Get{{.modelName}}List = '/{{.prefix}}/{{.modelNameSnake}}/list',
  Delete{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/delete',
  Get{{.modelName}}ById = '/{{.prefix}}/{{.modelNameSnake}}',
}

/**
 * @description: Get {{.modelNameSpace}} list
 */

export const get{{.modelName}}List = (params: BaseListReq, mode: ErrorMessageMode = 'notice') => {
  return defHttp.post<BaseDataResp<{{.modelName}}ListResp>>(
    { url: Api.Get{{.modelName}}List, params },
    { errorMessageMode: mode },
  );
};

/**
 *  @description: Create a new {{.modelNameSpace}}
 */
export const create{{.modelName}} = (params: {{.modelName}}Info, mode: ErrorMessageMode = 'notice') => {
  return defHttp.post<BaseResp>(
    { url: Api.Create{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
      successMessageMode: mode,
    },
  );
};

/**
 *  @description: Update the {{.modelNameSpace}}
 */
export const update{{.modelName}} = (params: {{.modelName}}Info, mode: ErrorMessageMode = 'notice') => {
  return defHttp.post<BaseResp>(
    { url: Api.Update{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
      successMessageMode: mode,
    },
  );
};

/**
 *  @description: Delete {{.modelNameSpace}}s
 */
export const delete{{.modelName}} = (params: Base{{if .useUUID}}UU{{end}}IDsReq, mode: ErrorMessageMode = 'notice') => {
  return defHttp.post<BaseResp>(
    { url: Api.Delete{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
      successMessageMode: mode,
    },
  );
};

/**
 *  @description: Get {{.modelNameSpace}} By ID
 */
export const get{{.modelName}}ById = (params: Base{{if .useUUID}}UU{{end}}IDReq, mode: ErrorMessageMode = 'notice') => {
  return defHttp.post<BaseDataResp<{{.modelName}}Info>>(
    { url: Api.Get{{.modelName}}ById, params: params },
    {
      errorMessageMode: mode,
    },
  );
};
