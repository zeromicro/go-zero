<template>
  <div>
    <BasicTable @register="registerTable">
      <template #tableTitle>
        <Button type="primary" danger v-if="showDeleteButton" @click="handleBatchDelete()">
          {{.deleteButtonTitle}}
        </Button>
      </template>
      <template #toolbar>
        <a-button type="primary" @click="handleCreate">
          {{.addButtonTitle}}
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                popConfirm: {
                  title: t('common.deleteConfirm'),
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <{{.modelName}}Drawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { createVNode, defineComponent, ref } from 'vue';
  import { Button, Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue/lib/icons';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';

  import { useDrawer } from '/@/components/Drawer';
  import {{.modelName}}Drawer from './{{.modelName}}Drawer.vue';
  import { useI18n } from 'vue-i18n';
  import { useMessage } from '/@/hooks/web/useMessage';

  import { columns, searchFormSchema } from './{{.modelNameLowerCase}}.data';
  import { get{{.modelName}}List, delete{{.modelName}}, batchDelete{{.modelName}} } from '/@/api/{{.folderName}}/{{.modelNameLowerCase}}';

  export default defineComponent({
    name: '{{.modelName}}Management',
    components: { BasicTable, {{.modelName}}Drawer, TableAction, Button },
    setup() {
      const { t } = useI18n();
      const selectedIds = ref<number[] | string[]>();
      const showDeleteButton = ref<boolean>(false);

      const [registerDrawer, { openDrawer }] = useDrawer();
      const { notification } = useMessage();
      const [registerTable, { reload }] = useTable({
        title: t('{{.folderName}}.{{.modelNameLowerCase}}.{{.modelNameLowerCase}}List'),
        api: get{{.modelName}}List,
        columns,
        formConfig: {
          labelWidth: 120,
          schemas: searchFormSchema,
        },
        useSearchForm: true,
        showTableSetting: true,
        bordered: true,
        showIndexColumn: false,
        actionColumn: {
          width: 30,
          title: t('common.action'),
          dataIndex: 'action',
          fixed: undefined,
        },
        rowKey: 'id',
        rowSelection: {
          type: 'checkbox',
          onChange: (selectedRowKeys, _selectedRows) => {
            selectedIds.value = selectedRowKeys;
            if (selectedRowKeys.length > 0) {
              showDeleteButton.value = true;
            } else {
              showDeleteButton.value = false;
            }
          },
        },
      });

      function handleCreate() {
        openDrawer(true, {
          isUpdate: false,
        });
      }

      function handleEdit(record: Recordable) {
        openDrawer(true, {
          record,
          isUpdate: true,
        });
      }

      async function handleDelete(record: Recordable) {
        const result = await delete{{.modelName}}({ id: record.id }, 'modal');
        notification.success({
          message: t('common.successful'),
          description: t(result.msg),
          duration: 3,
        });
        reload();
      }

      async function handleBatchDelete() {
        Modal.confirm({
          title: t('common.deleteConfirm'),
          icon: createVNode(ExclamationCircleOutlined),
          async onOk() {
            const result = await batchDelete{{.modelName}}(
              { ids: selectedIds.value as number[] },
              'modal',
            );
            if (result.code === 0) {
              reload();
              showDeleteButton.value = false;
            }
          },
          onCancel() {
            console.log('Cancel');
          },
        });
      }

      function handleSuccess() {
        reload();
      }

      return {
        t,
        registerTable,
        registerDrawer,
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        handleBatchDelete,
        showDeleteButton,
      };
    },
  });
</script>
