import React, { useState } from "react";
import { Col, Flex, Form, Input, Row, Select } from "antd";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData, Method, ContentType } from "./_defaultProps";

interface RequestLinePanelProps {
  routeField: FormListFieldData;
}

const RequestLinePanel: React.FC<
  RequestLinePanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const routeField = props.routeField;
  const [initRequestValues, setInitRequestValues] = useState([]);

  return (
    <div>
      <Flex gap={16} wrap>
        <Form.Item
          style={{ flex: "0.75" }}
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
                  />
                </Form.Item>
              </div>
            }
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
          />
        </Form.Item>
      </Flex>
    </div>
  );
};

export default RequestLinePanel;
