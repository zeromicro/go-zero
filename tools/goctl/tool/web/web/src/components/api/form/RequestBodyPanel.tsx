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
  Tooltip,
} from "antd";
import { CloseOutlined } from "@ant-design/icons";
import { FormListFieldData, FormListOperation } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData } from "./_defaultProps";
import CodeMirror, { EditorView } from "@uiw/react-codemirror";
import { githubLight } from "@uiw/codemirror-theme-github";
import { langs } from "@uiw/codemirror-extensions-langs";
import { ConverterIcon } from "../../../util/icon";
import type { GetRef } from "antd";
import RequestFieldPanel from "./RequestFieldPanel";

type FormInstance<T> = GetRef<typeof Form<T>>;

interface RequestBodyPanelProps {
  routeGroupField: FormListFieldData;
  requestBodyFields: FormListFieldData[];
  requestBodyOpt: FormListOperation;
  routeField: FormListFieldData;
  form: FormInstance<any>;
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
  const [modalOpen, setModalOpen] = useState(false);
  const [requestBodyParseCode, setRequestBodyParseCode] = useState("");
  const [api, contextHolder] = notification.useNotification();

  return (
    <div>
      {contextHolder}
      {/*request body import modal*/}
      <Modal
        title={t("formRequestBodyFieldBtnImport")}
        centered
        open={modalOpen}
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
            setModalOpen(false);
          } catch (err) {
            api.error({
              message: t("tipsInvalidJSON") + ": " + err,
            });
            return;
          }
        }}
        onCancel={() => setModalOpen(false)}
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
        <span
          style={{
            position: "absolute",
            top: -25,
            right: 0,
            zIndex: 1000,
          }}
        >
          <Tooltip title={t("formRequestBodyFieldBtnImport")}>
            <ConverterIcon
              type={"icon-import"}
              style={{ cursor: "pointer", fontSize: 18, color: "#333333" }}
              onClick={() => {
                setModalOpen(true);
              }}
            />
          </Tooltip>
        </span>

        {requestBodyFields.map((requestBodyField) => (
          <RequestFieldPanel
            requestBodyField={requestBodyField}
            requestBodyOpt={requestBodyOpt}
          />
        ))}
        <Button
          type="dashed"
          onClick={() => {
            requestBodyOpt.add();
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
