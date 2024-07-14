import React, { useState } from "react";
import { Col, Form, Input, Row, Select } from "antd";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData, Method, ContentType } from "./_defaultProps";

interface RequestLintPanelProps {
  routeField: FormListFieldData;
}

const RequestLintPanel: React.FC<
  RequestLintPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const routeField = props.routeField;
  const [initRequestValues, setInitRequestValues] = useState([]);

  return (
    <div>
      <Row gutter={16}>
        <Col span={16}>
          <Form.Item
            label={t("formPathTitle")}
            name={[routeField.name, "path"]}
          >
            <Input
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
        </Col>
        <Col span={8}>
          <Form.Item
            label={t("formContentTypeTitle")}
            name={[routeField.name, "contentType"]}
          >
            <Select
              defaultValue={ContentType.ApplicationJson}
              options={RoutePanelData.ContentTypeOptions}
            />
          </Form.Item>
        </Col>
      </Row>
    </div>
  );
};

export default RequestLintPanel;
