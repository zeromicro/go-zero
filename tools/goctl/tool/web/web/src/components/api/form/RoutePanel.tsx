import React, { useState } from "react";
import { Button, Collapse, Form, Modal, notification } from "antd";
import { CloseOutlined } from "@ant-design/icons";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import CodeMirror, { EditorView } from "@uiw/react-codemirror";
import { githubLight } from "@uiw/codemirror-theme-github";
import { langs } from "@uiw/codemirror-extensions-langs";
import type { FormInstance } from "antd/es/form/hooks/useForm";
import RequestLinePanel from "./RequestLinePanel";
import RequestBodyPanel from "./RequestBodyPanel";
import CodeMirrorPanel from "./CodeMirrorPanel";

interface RoutePanelProps {
  routeGroupField: FormListFieldData;
  form: FormInstance;
}

const RoutePanel: React.FC<
  RoutePanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const routeGroupField = props.routeGroupField;
  const form = props.form;
  const [initRequestValues, setInitRequestValues] = useState([
    {
      path: "/your/path",
      contentType: "application/json",
      method: "POST",
      requestBodyFields: [
        {
          name: "FieldName",
          type: "string",
          optional: false,
        },
      ],
      responseBody: "{}",
    },
  ]);
  const [requestBodyModalOpen, setRequestBodyModalOpen] = useState(false);
  const [responseBodyModalOpen, setResponseBodyModalOpen] = useState(false);
  const [requestBodyParseCode, setRequestBodyParseCode] = useState("");
  const [responseCode, setResponseCode] = useState("");
  const [api, contextHolder] = notification.useNotification();
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
      {/*response body editor*/}
      <Modal
        title={t("formResponseBodyModelTitle")}
        centered
        open={responseBodyModalOpen}
        maskClosable={false}
        keyboard={false}
        closable={false}
        destroyOnClose
        onOk={() => {
          Modal.destroyAll();
          setResponseBodyModalOpen(false);
        }}
        onCancel={() => {
          Modal.destroyAll();
          setResponseBodyModalOpen(false);
        }}
        width={1000}
        cancelText={t("formResponseBodyModalCancel")}
        okText={t("formResponseBodyModalConfirm")}
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
          value={responseCode}
          onChange={(code) => {
            setResponseCode(code);
          }}
        />
      </Modal>
      <Form.Item label={t("formRouteListTitle")}>
        <Form.List
          initialValue={initRequestValues}
          name={[routeGroupField.name, "routes"]}
        >
          {(routeFields, routeOpt) => (
            <div
              style={{
                display: "flex",
                rowGap: 16,
                flexDirection: "column",
              }}
            >
              {routeFields.map((routeField) => (
                <Collapse
                  defaultActiveKey={[routeField.key]}
                  items={[
                    {
                      key: routeField.key,
                      label: t("formRouteTitle") + `${routeField.name + 1}`,
                      children: (
                        <div>
                          {/*request line component*/}
                          <RequestLinePanel routeField={routeField} />
                          {/*request body*/}
                          <Form.Item
                            tooltip={t("formRequestBodyTooltip")}
                            label={t("formRequestBodyTitle")}
                          >
                            <Form.List
                              name={[routeField.name, "requestBodyFields"]}
                            >
                              {(requestBodyFields, requestBodyOpt) => (
                                <RequestBodyPanel
                                  routeGroupField={routeGroupField}
                                  requestBodyFields={requestBodyFields}
                                  requestBodyOpt={requestBodyOpt}
                                  routeField={routeField}
                                  form={form}
                                />
                              )}
                            </Form.List>
                          </Form.Item>
                          {/*  response body  */}
                          <Form.Item
                            label={t("formResponseBodyTitle")}
                            name={[routeField.name, "responseBody"]}
                          >
                            <CodeMirrorPanel
                              value={responseCode}
                              placeholder={t("formResponseBodyPlaceholder")}
                              onChange={(code) => {
                                setResponseCode(code);
                              }}
                            />
                          </Form.Item>
                        </div>
                      ),
                      extra: (
                        <CloseOutlined
                          onClick={() => {
                            routeOpt.remove(routeField.name);
                          }}
                        />
                      ),
                    },
                  ]}
                />
              ))}
              <Button type="dashed" onClick={() => routeOpt.add()} block>
                + {t("formButtonRouteAdd")}
              </Button>
            </div>
          )}
        </Form.List>
      </Form.Item>
    </div>
  );
};

export default RoutePanel;
