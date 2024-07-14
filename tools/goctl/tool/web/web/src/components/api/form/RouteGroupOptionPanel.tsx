import React from "react";
import { Col, Form, Input, InputNumber, Row, Switch } from "antd";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";

interface RouteGroupOptionPanelProps {
  routeGroupField: FormListFieldData;
}

const RouteGroupOptionPanel: React.FC<
  RouteGroupOptionPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const routeGroupField = props.routeGroupField;
  return (
    <div>
      <Row gutter={16}>
        <Col span={8}>
          <Form.Item
            label={t("formJwtTitle")}
            name={[routeGroupField.name, "jwt"]}
            tooltip={t("formJWTTips")}
          >
            <Switch defaultChecked />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            label={t("formPrefixTitle")}
            name={[routeGroupField.name, "prefix"]}
          >
            <Input
              allowClear
              prefix={"/"}
              placeholder={t("formPrefixPlaceholder")}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            label={t("formGroupTitle")}
            name={[routeGroupField.name, "group"]}
            tooltip={t("formRouteGroupTooltip")}
          >
            <Input allowClear placeholder={t("formGroupPlaceholder")} />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={8}>
          <Form.Item
            label={t("formTimeoutTitle")}
            name={[routeGroupField.name, "timeout"]}
          >
            <InputNumber min={0} precision={0} addonAfter="ms" />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            label={t("formMiddlewareTitle")}
            name={[routeGroupField.name, "middleware"]}
            tooltip={t("formMiddlewareTips")}
          >
            <Input allowClear placeholder={t("formMiddlewarePlaceholder")} />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item
            label={t("formMaxByteTitle")}
            name={[routeGroupField.name, "maxByte"]}
          >
            <InputNumber min={0} precision={0} addonAfter="byte" />
          </Form.Item>
        </Col>
      </Row>
    </div>
  );
};

export default RouteGroupOptionPanel;
