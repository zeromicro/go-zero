import React, { useState } from "react";
import { Flex, Form, Input, Select } from "antd";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { ContentType, Method, RoutePanelData } from "./_defaultProps";

interface RequestLinePanelProps {
  routeField: FormListFieldData;
}

const RequestLinePanel: React.FC<
  RequestLinePanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const routeField = props.routeField;
  const [disableContentType, setDisableContentType] = useState(false);
  const [contentType, setContentType] = useState(ContentType.ApplicationJson);

  return (
    <div>
      <Flex vertical wrap>
        <Flex gap={16} wrap>
          <Form.Item
            style={{ flex: 1 }}
            label={t("formPathTitle")}
            name={[routeField.name, "path"]}
            rules={[
              {
                required: true,
                message: `${t("formInputPrefix")}${t("formPathTitle")}`,
              },
              {
                pattern: RoutePanelData.PathPattern,
                message: `${t("formPathTitle")}${t("formRegexTooltip")}: ${RoutePanelData.PathPattern}`,
              },
            ]}
          >
            <Input
              placeholder={`${t("formInputPrefix")}${t("formPathTitle")}`}
              allowClear
              addonBefore={
                <div>
                  <Form.Item noStyle name={[routeField.name, "method"]}>
                    <Select
                      style={{ width: 100 }}
                      defaultValue={Method.POST}
                      options={RoutePanelData.MethodOptions}
                      onSelect={(value) => {
                        if (value === Method.GET.toLowerCase()) {
                          setDisableContentType(true);
                          setContentType(ContentType.ApplicationForm);
                        } else {
                          setDisableContentType(false);
                        }
                      }}
                    />
                  </Form.Item>
                </div>
              }
            />
          </Form.Item>
        </Flex>
        <Flex gap={8} wrap>
          <Form.Item
            style={{ flex: "0.75" }}
            label={t("formHandlerTitle")}
            name={[routeField.name, "handler"]}
            tooltip={t("formHandlerTooltip")}
            rules={[
              {
                pattern: RoutePanelData.IDPattern,
                message: `${t("formHandlerTitle")}${t("formRegexTooltip")}: ${RoutePanelData.IDPattern}`,
              },
            ]}
          >
            <Input
              placeholder={`${t("formInputPrefix")}${t("formHandlerTitle")}`}
              allowClear
            />
          </Form.Item>
          <Form.Item
            style={{ flex: "0.25" }}
            label={t("formContentTypeTitle")}
            name={[routeField.name, "contentType"]}
          >
            <Select
              defaultValue={ContentType.ApplicationJson}
              options={RoutePanelData.ContentTypeOptions}
              disabled={disableContentType}
              value={contentType}
            />
          </Form.Item>
        </Flex>
      </Flex>
    </div>
  );
};

export default RequestLinePanel;
