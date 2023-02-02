import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import { useI18n } from '/@/hooks/web/useI18n';
import { formatToDateTime } from '/@/utils/dateUtil';
{{if .hasStatus}}import { set{{.modelName}}Status } from '/@/api/{{.folderName}}/{{.modelNameLowerCase}}';
import { Switch } from 'ant-design-vue';
import { useMessage } from '/@/hooks/web/useMessage';
import { h } from 'vue';{{end}}

const { t } = useI18n();

export const columns: BasicColumn[] = [{{.basicData}}
{{if .useBaseInfo}}  {
    title: t('common.createTime'),
    dataIndex: 'createdAt',
    width: 70,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt);
    },
  },{{end}}
];

export const searchFormSchema: FormSchema[] = [{{.searchFormData}}
];

export const formSchema: FormSchema[] = [
  {
    field: 'id',
    label: 'ID',
    component: 'Input',
    show: false,
  },
{{.formData}}
];
