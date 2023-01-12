<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="500px"
    @ok="handleSubmit"
  >
    <BasicForm @register="registerForm" />
  </BasicDrawer>
</template>
<script lang="ts">
  import { defineComponent, ref, computed, unref } from 'vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { formSchema } from './{{.modelNameLowerCase}}.data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import { useI18n } from 'vue-i18n';

  import { {{.modelName}}Info } from '/@/api/{{.folderName}}/model/{{.modelNameLowerCase}}Model';
  import { createOrUpdate{{.modelName}} } from '/@/api/{{.folderName}}/{{.modelNameLowerCase}}';

  export default defineComponent({
    name: '{{.modelName}}Drawer',
    components: { BasicDrawer, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const isUpdate = ref(true);
      const { t } = useI18n();

      const [registerForm, { resetFields, setFieldsValue, validate }] = useForm({
        labelWidth: 90,
        baseColProps: { span: 24 },
        schemas: formSchema,
        showActionButtonGroup: false,
      });

      const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
        resetFields();
        setDrawerProps({ confirmLoading: false });

        isUpdate.value = !!data?.isUpdate;

        if (unref(isUpdate)) {
          setFieldsValue({
            ...data.record,
          });
        }
      });

      const getTitle = computed(() =>
        !unref(isUpdate) ? t('{{.folderName}}.{{.modelNameLowerCase}}.add{{.modelName}}') : t('{{.folderName}}.{{.modelNameLowerCase}}.edit{{.modelName}}'),
      );

      async function handleSubmit() {
        const values = await validate();
        setDrawerProps({ confirmLoading: true });
        let id: {{if .useUUID}}string{{else}}number{{end}};
        if (unref(isUpdate)) {
          id = {{if .useUUID}}values['id'];{{else}}Number(values['id']);{{end}}
        } else {
          id = {{if .useUUID}}''{{else}}0{{end}};
        }
        let params: {{.modelName}}Info = {
          id: id,{{.infoData}}
        };
        let result = await createOrUpdate{{.modelName}}(params);
        if (result.code === 0) {
          closeDrawer();
          emit('success');
        }
        setDrawerProps({ confirmLoading: false });
      }

      return {
        registerDrawer,
        registerForm,
        getTitle,
        handleSubmit,
      };
    },
  });
</script>
