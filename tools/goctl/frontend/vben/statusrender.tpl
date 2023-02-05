
  {
    title: t('common.statusName'),
    dataIndex: 'status',
    width: 50,
    customRender: ({ record }) => {
      if (!Reflect.has(record, 'pendingStatus')) {
        record.pendingStatus = false;
      }
      return h(Switch, {
        checked: record.status === 1,
        checkedChildren: t('common.on'),
        unCheckedChildren: t('common.off'),
        loading: record.pendingStatus,
        onChange(checked: boolean) {
          const { createMessage } = useMessage();
          record.pendingStatus = true;
          const newStatus = checked ? 1 : 0;
          set{{.modelName}}Status(record.id, newStatus)
            .then((data) => {
              record.status = newStatus;
              if (data.code == 0) createMessage.success(t('common.changeStatusSuccess'));
            })
            .catch(() => {
              createMessage.error(t('common.changeStatusFailed'));
            })
            .finally(() => {
              record.pendingStatus = false;
            });
        },
      });
    },
  },