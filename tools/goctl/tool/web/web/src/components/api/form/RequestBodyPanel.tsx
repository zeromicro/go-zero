import React, { useState } from "react";
import {
  Button,
  Flex,
  Form,
  Input,
  Select,
  Modal,
  notification,
  Switch,
} from "antd";
import { CloseOutlined } from "@ant-design/icons";
import { FormListFieldData, FormListOperation } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData } from "./_defaultProps";
import CodeMirror, { EditorView } from "@uiw/react-codemirror";
import { githubLight } from "@uiw/codemirror-theme-github";
import { langs } from "@uiw/codemirror-extensions-langs";
import type { FormInstance } from "antd/es/form/hooks/useForm";

interface RequestBodyPanelProps {
  routeGroupField: FormListFieldData;
  requestBodyFields: FormListFieldData[];
  requestBodyOpt: FormListOperation;
  routeField: FormListFieldData;
  form: FormInstance;
}

const RequestBodyPanel: React.FC<
  RequestBodyPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const requestBodyFields = props.requestBodyFields;
  const requestBodyOpt = props.requestBodyOpt;
  const routeGroupField = props.routeGroupField;
  const routeField = props.routeField;
  const form = props.form;
  const [initRequestValues, setInitRequestValues] = useState([]);
  const [requestBodyModalOpen, setRequestBodyModalOpen] = useState(false);
  const [requestBodyParseCode, setRequestBodyParseCode] = useState("");
  const [api, contextHolder] = notification.useNotification();
  const [showImportButton, setShowImportButton] = useState(true);

  const canChowImportButton = () => {
    const routeGroups = form.getFieldValue(`routeGroups`);
    if (!routeGroups) {
      setShowImportButton(true);
      return;
    }

    if (routeGroups.length <= routeGroupField.key) {
      setShowImportButton(true);
      return;
    }

    const routeGroup = routeGroups[routeGroupField.key];

    if (!routeGroup) {
      setShowImportButton(true);
      return;
    }
    if (!routeGroup.routes) {
      setShowImportButton(true);
      return;
    }

    if (routeGroup.routes.length <= routeField.key) {
      setShowImportButton(true);
      return;
    }

    const route = routeGroup.routes[routeField.key];
    if (!route) {
      setShowImportButton(true);
      return;
    }
    if (!route.requestBodyFields) {
      setShowImportButton(true);
      return;
    }

    setShowImportButton(route.requestBodyFields.length === 0);
  };
  return (
    <div>
      {contextHolder}
      {/*request body import modal*/}
      <Modal
        title={t("formRequestBodyFieldBtnImport")}
        centered
        open={requestBodyModalOpen}
        maskClosable={false}
        keyboard={false}
        closable={false}
        destroyOnClose
        onOk={() => {
          try {
            const obj = JSON.parse(requestBodyParseCode);
            if (Array.isArray(obj)) {
              api.error({
                message: t("tipsInvalidJSONArray"),
              });
              return;
            }

            // todo: 从后段解析数据
            setRequestBodyModalOpen(false);
          } catch (err) {
            api.error({
              message: t("tipsInvalidJSON") + ": " + err,
            });
            return;
          }
        }}
        onCancel={() => setRequestBodyModalOpen(false)}
        width={1000}
        cancelText={t("formRequestBodyModalCancel")}
        okText={t("formRequestBodyModalConfirm")}
      >
        <CodeMirror
          style={{ marginTop: 10, overflow: "auto" }}
          extensions={[
            langs.json(),
            EditorView.theme({
              "&.cm-focused": {
                outline: "none",
              },
            }),
          ]}
          theme={githubLight}
          height={"70vh"}
          onChange={(code) => {
            setRequestBodyParseCode(code);
          }}
        />
      </Modal>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
        }}
      >
        {requestBodyFields.map((requestBodyField) => (
          <Flex key={requestBodyField.key} gap={10} wrap>
            <Form.Item
              label={t("formRequestBodyFieldNameTitle")}
              name={[requestBodyField.name, "name"]}
              style={{ flex: 1 }}
              tooltip={t("formRequestBodyFieldNameTooltip")}
            >
              <Input />
            </Form.Item>
            <Form.Item
              label={t("formRequestBodyFieldTypeTitle")}
              name={[requestBodyField.name, "type"]}
              style={{ flex: 1 }}
            >
              <Select options={RoutePanelData.GolangTypeOptions} />
            </Form.Item>
            <Form.Item
              label={t("formRequestBodyFieldOptionalTitle")}
              name={[requestBodyField.name, "optional"]}
            >
              <Switch />
            </Form.Item>
            <CloseOutlined
              onClick={() => {
                requestBodyOpt.remove(requestBodyField.name);
                canChowImportButton();
              }}
            />
          </Flex>
        ))}
        {showImportButton ? (
          <Button
            style={{ marginBottom: 16 }}
            type="dashed"
            onClick={() => setRequestBodyModalOpen(true)}
            block
          >
            🔍 {t("formRequestBodyFieldBtnImport")}
          </Button>
        ) : (
          <></>
        )}
        <Button
          type="dashed"
          onClick={() => {
            requestBodyOpt.add();
            canChowImportButton();
          }}
          block
        >
          + {t("formRequestBodyFieldBtnAdd")}
        </Button>
      </div>
    </div>
  );
};

export default RequestBodyPanel;
