import React from "react";
import { Col, Flex, Form, Input, InputNumber, Row, Switch } from "antd";
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
      <Flex gap={16} wrap>
        <Form.Item
          label={t("formJwtTitle")}
          name={[routeGroupField.name, "jwt"]}
          tooltip={t("formJWTTips")}
        >
          <Switch defaultChecked />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formPrefixTitle")}
          name={[routeGroupField.name, "prefix"]}
        >
          <Input
            allowClear
            prefix={"/"}
            placeholder={t("formPrefixPlaceholder")}
          />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formGroupTitle")}
          name={[routeGroupField.name, "group"]}
          tooltip={t("formRouteGroupTooltip")}
        >
          <Input allowClear placeholder={t("formGroupPlaceholder")} />
        </Form.Item>
      </Flex>

      <Flex gap={16} wrap>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formTimeoutTitle")}
          name={[routeGroupField.name, "timeout"]}
        >
          <InputNumber min={0} precision={0} addonAfter="ms" />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formMiddlewareTitle")}
          name={[routeGroupField.name, "middleware"]}
          tooltip={t("formMiddlewareTips")}
        >
          <Input allowClear placeholder={t("formMiddlewarePlaceholder")} />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formMaxByteTitle")}
          name={[routeGroupField.name, "maxByte"]}
        >
          <InputNumber min={0} precision={0} addonAfter="byte" />
        </Form.Item>
      </Flex>
    </div>
  );
};

export default RouteGroupOptionPanel;
