import { defHttp } from '/@/utils/http/axios';
import { ErrorMessageMode } from '/#/axios';
import { BaseDataResp, BaseIdReq, BasePageReq, BaseResp, BaseIdsReq } from '/@/api/model/baseModel';
import { {{.modelName}}Info, {{.modelName}}ListResp } from './model/{{.modelNameLowerCase}}Model';

enum Api {
  CreateOrUpdate{{.modelName}} = '/{{.prefix}}/{{.modelNameLowerCase}}/create_or_update',
  Get{{.modelName}}List = '/{{.prefix}}/{{.modelNameLowerCase}}/list',
  Delete{{.modelName}} = '/{{.prefix}}/{{.modelNameLowerCase}}/delete',
  BatchDelete{{.modelName}} = '/{{.prefix}}/{{.modelNameLowerCase}}/batch_delete',
}

/**
 * @description: Get {{.modelNameLowerCase}} list
 */

export const get{{.modelName}}List = (params: BasePageReq) => {
  return defHttp.post<BaseDataResp<{{.modelName}}ListResp>>({ url: Api.Get{{.modelName}}List, params });
};

/**
 *  @description: create a new {{.modelNameLowerCase}}
 */
export const createOrUpdate{{.modelName}} = (params: {{.modelName}}Info, mode: ErrorMessageMode = 'modal') => {
  return defHttp.post<BaseResp>(
    { url: Api.CreateOrUpdate{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
    },
  );
};

/**
 *  @description: delete {{.modelNameLowerCase}}
 */
export const delete{{.modelName}} = (params: BaseIdReq, mode: ErrorMessageMode = 'modal') => {
  return defHttp.post<BaseResp>(
    { url: Api.Delete{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
    },
  );
};

/**
 *  @description: batch delete {{.modelNameLowerCase}}s
 */
export const batchDelete{{.modelName}} = (params: BaseIdsReq, mode: ErrorMessageMode = 'modal') => {
  return defHttp.post<BaseResp>(
    { url: Api.BatchDelete{{.modelName}}, params: params },
    {
      errorMessageMode: mode,
    },
  );
};
