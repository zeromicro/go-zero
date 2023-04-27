import { BasicColumn, FormSchema } from '/@/components/Table';
import { useI18n } from '/@/hooks/web/useI18n';
import { formatToDateTime } from '/@/utils/dateUtil';
{{if .hasStatus}}import { update{{.modelName}} } from '/@/api/{{.folderName}}/{{.modelNameLowerCamel}}';
import { Switch } from 'ant-design-vue';
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
