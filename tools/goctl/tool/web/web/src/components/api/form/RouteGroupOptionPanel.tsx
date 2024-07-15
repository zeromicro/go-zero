import React from "react";
import { Flex, Form, Input, InputNumber, Switch } from "antd";
import { FormListFieldData } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData } from "./_defaultProps";

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
      <Flex gap={8} wrap>
        <Form.Item
          label={t("formJwtTitle")}
          name={[routeGroupField.name, "jwt"]}
          tooltip={t("formJWTTips")}
        >
          <Switch />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formPrefixTitle")}
          name={[routeGroupField.name, "prefix"]}
          rules={[
            {
              pattern: RoutePanelData.PrefixPathPattern,
              message: `${t("formPrefixTitle")}${t("formRegexTooltip")}: ${RoutePanelData.PrefixPathPattern}`,
            },
          ]}
        >
          <Input
            allowClear
            placeholder={`${t("formInputPrefix")}${t("formPrefixTitle")}`}
          />
        </Form.Item>
        <Form.Item
          style={{ flex: 1 }}
          label={t("formGroupTitle")}
          name={[routeGroupField.name, "group"]}
          tooltip={t("formRouteGroupTooltip")}
          rules={[
            {
              pattern: RoutePanelData.IDPattern,
              message: `${t("formGroupTitle")}${t("formRegexTooltip")}: ${RoutePanelData.IDPattern}`,
            },
          ]}
        >
          <Input
            allowClear
            placeholder={`${t("formInputPrefix")}${t("formGroupTitle")}`}
          />
        </Form.Item>
      </Flex>

      <Flex gap={8} wrap>
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
          rules={[
            {
              pattern: RoutePanelData.IDCommaPattern,
              message: `${t("formMiddlewareTitle")}${t("formRegexTooltip")}: ${RoutePanelData.IDCommaPattern}`,
            },
          ]}
        >
          <Input
            allowClear
            placeholder={`${t("formInputPrefix")}${t("formMiddlewareTitle")}`}
          />
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
