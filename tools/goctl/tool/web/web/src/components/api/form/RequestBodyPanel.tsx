import React, { useState } from "react";
import { Button, Form, Modal, notification, Tooltip } from "antd";
import { FormListFieldData, FormListOperation } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import CodeMirror, { EditorView } from "@uiw/react-codemirror";
import { githubLight } from "@uiw/codemirror-theme-github";
import { langs } from "@uiw/codemirror-extensions-langs";
import { ConverterIcon } from "../../../util/icon";
import type { GetRef } from "antd";
import RequestFieldPanel from "./RequestFieldPanel";
import { Http, ParseBodyForm } from "../../../util/http";

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
  const [initRequestValues, setInitRequestValues] = useState<ParseBodyForm[]>(
    [],
  );
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

            Http.ParseBodyFromJson(
              requestBodyParseCode,
              (data: ParseBodyForm[]) => {
                if (!data || data.length === 0) {
                  return;
                }
                let routeGroups = form.getFieldValue("routeGroups");
                if (!routeGroups) {
                  return;
                }
                let routeGroup = routeGroups[routeGroupField.key];
                if (!routeGroup) {
                  return;
                }
                let routes = routeGroup.routes;
                if (!routes.length) {
                  return;
                }
                let route = routes[routeField.key];
                if (!route) {
                  return;
                }
                let requestBodyFields = route.requestBodyFields;
                if (!requestBodyFields) {
                  return;
                }
                for (let i = 0; i < data.length; i++) {
                  const item = data[i];
                  requestBodyFields.push(item);
                }

                form.setFieldValue("routeGroups", routeGroups);
                setModalOpen(false);
              },
              (error) => {
                api.error({
                  message: error,
                });
              },
            );
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
              style={{ cursor: "pointer", fontSize: 30 }}
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
            routeField={routeField}
            form={form}
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
