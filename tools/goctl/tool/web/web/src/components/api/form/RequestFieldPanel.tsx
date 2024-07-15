import React, { useEffect, useState } from "react";
import {
  Flex,
  Form,
  Input,
  Select,
  Modal,
  notification,
  Switch,
  Tooltip,
  Radio,
  RadioChangeEvent,
  type GetRef,
} from "antd";
import { CloseOutlined, SettingOutlined } from "@ant-design/icons";
import { FormListFieldData, FormListOperation } from "antd/es/form/FormList";
import { useTranslation } from "react-i18next";
import { RoutePanelData } from "./_defaultProps";
import { InputNumber } from "antd";

type FormInstance<T> = GetRef<typeof Form<T>>;

interface RequestFieldPanelProps {
  requestBodyField: FormListFieldData;
  requestBodyOpt: FormListOperation;

  routeField: FormListFieldData;
  form: FormInstance<any>;
}

const RequestFieldPanel: React.FC<
  RequestFieldPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const form = props.form;
  const routeField = props.routeField;
  const requestBodyField = props.requestBodyField;
  const requestBodyOpt = props.requestBodyOpt;
  const [initRequestValues, setInitRequestValues] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  useState<FormListFieldData>();
  const [api, contextHolder] = notification.useNotification();
  const [typeIsNumber, setTypeIsNumber] = useState(false);
  const [checkValue, setCheckValue] = useState("enum");

  return (
    <div>
      {contextHolder}
      {/*request body import modal*/}
      <Modal
        title={t("formRequestBodyFieldModelTitle")}
        centered
        open={modalOpen}
        maskClosable={false}
        keyboard={false}
        closable={false}
        destroyOnClose
        cancelButtonProps={{
          style: {
            visibility: "hidden",
          },
        }}
        onOk={() => {
          setModalOpen(false);
        }}
        onCancel={() => setModalOpen(false)}
        width={500}
        okText={"OK"}
      >
        <Flex vertical gap={6} wrap>
          <Form.Item
            label={t("formRequestBodyFieldNameTitle")}
            name={[requestBodyField.name, "name"]}
            style={{ flex: 1, marginTop: 20 }}
            tooltip={t("formRequestBodyFieldNameTooltip")}
            rules={[
              {
                required: true,
                message: `${t("formInputPrefix")}${t("formRequestBodyFieldNameTitle")}`,
              },
              {
                pattern: RoutePanelData.IDPattern,
                message: `${t("formRequestBodyFieldNameTitle")}${t("formRegexTooltip")}: ${RoutePanelData.IDPattern}`,
              },
            ]}
          >
            <Input
              allowClear
              placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldNameTitle")}`}
            />
          </Form.Item>
          <Form.Item
            shouldUpdate
            label={t("formRequestBodyFieldTypeTitle")}
            name={[requestBodyField.name, "type"]}
            style={{ flex: 1 }}
            rules={[
              {
                required: true,
                message: `${t("formInputPrefix")}${t("formRequestBodyFieldTypeTitle")}`,
              },
            ]}
          >
            <Select
              allowClear
              placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldTypeTitle")}`}
              options={RoutePanelData.GolangTypeOptions}
              showSearch
              onSelect={(value) => {
                const isNumberType = RoutePanelData.IsNumberType(value);
                setTypeIsNumber(isNumberType);
                if (!isNumberType) {
                  setCheckValue("enum");
                }
              }}
            />
          </Form.Item>
          <Form.Item
            label={t("formRequestBodyFieldOptionalTitle")}
            name={[requestBodyField.name, "optional"]}
          >
            <Switch />
          </Form.Item>
          <Form.Item
            label={t("formRequestBodyFieldDefaultTitle")}
            name={[requestBodyField.name, "defaultValue"]}
          >
            <Input
              allowClear
              placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldDefaultTitle")}`}
            />
          </Form.Item>
          {typeIsNumber ? (
            <Form.Item name={[requestBodyField.name, "checkEnum"]}>
              <Radio.Group
                onChange={(e: RadioChangeEvent) => {
                  setCheckValue(e.target.value);
                }}
                defaultValue={"enum"}
              >
                <Radio value={"enum"}>
                  {t("formRequestBodyFieldEnumTitle")}
                </Radio>
                <Radio value={"range"}>
                  {t("formRequestBodyFieldRangeTitle")}
                </Radio>
              </Radio.Group>
            </Form.Item>
          ) : (
            <></>
          )}

          {checkValue === "enum" ? (
            <Form.Item
              label={t("formRequestBodyFieldEnumTitle")}
              name={[requestBodyField.name, "enumValue"]}
              tooltip={t("formRequestBodyFieldEnumTooltip")}
              rules={[
                {
                  pattern: RoutePanelData.EnumCommaPattern,
                  message: `${t("formRequestBodyFieldEnumTitle")}${t("formRegexTooltip")}: ${RoutePanelData.EnumCommaPattern}`,
                },
              ]}
            >
              <Input
                allowClear
                placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldEnumTitle")}`}
              />
            </Form.Item>
          ) : (
            <Form.Item
              label={t("formRequestBodyFieldRangeTitle")}
              tooltip={t("formRequestBodyFieldRangeTooltip")}
            >
              <Flex justify={"space-between"} align={"center"}>
                <Form.Item noStyle name={[requestBodyField.name, "lowerBound"]}>
                  <InputNumber style={{ width: "50%" }} />
                </Form.Item>
                <div
                  style={{
                    width: 10,
                    height: 1,
                    background: "#c1c1c1",
                    marginLeft: 8,
                    marginRight: 8,
                  }}
                />
                <Form.Item noStyle name={[requestBodyField.name, "upperBound"]}>
                  <InputNumber style={{ width: "50%" }} />
                </Form.Item>
              </Flex>
            </Form.Item>
          )}
        </Flex>
      </Modal>

      <Flex align={"center"} key={requestBodyField.key} gap={6} wrap>
        <Form.Item
          label={t("formRequestBodyFieldNameTitle")}
          name={[requestBodyField.name, "name"]}
          style={{ flex: 1 }}
          tooltip={t("formRequestBodyFieldNameTooltip")}
          rules={[
            {
              required: true,
              message: `${t("formInputPrefix")}${t("formRequestBodyFieldNameTitle")}`,
            },
            {
              pattern: RoutePanelData.IDPattern,
              message: `${t("formRequestBodyFieldNameTitle")}${t("formRegexTooltip")}: ${RoutePanelData.IDPattern}`,
            },
          ]}
        >
          <Input
            allowClear
            placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldNameTitle")}`}
          />
        </Form.Item>
        <Form.Item
          label={t("formRequestBodyFieldTypeTitle")}
          name={[requestBodyField.name, "type"]}
          style={{ flex: 1 }}
          rules={[
            {
              required: true,
              message: `${t("formInputPrefix")}${t("formRequestBodyFieldTypeTitle")}`,
            },
          ]}
        >
          <Select
            allowClear
            placeholder={`${t("formInputPrefix")}${t("formRequestBodyFieldTypeTitle")}`}
            options={RoutePanelData.GolangTypeOptions}
            showSearch
            onSelect={(value) => {
              const isNumberType = RoutePanelData.IsNumberType(value);
              setTypeIsNumber(isNumberType);
              if (!isNumberType) {
                setCheckValue("enum");
              }
            }}
          />
        </Form.Item>
        <Tooltip title={t("formRequestBodySettings")}>
          <SettingOutlined
            style={{
              padding: 4,
              border: "dashed 1px #c1c1c1",
              marginTop: 4,
            }}
            onClick={() => {
              setModalOpen(true);
              let routeGroups = form.getFieldValue("routeGroups");
              if (!routeGroups) {
                return;
              }
              let routeGroup = routeGroups[routeField.key];
              if (!routeGroup) {
                return;
              }
              let route = routeGroup.routes[routeField.key];
              if (!route) {
                return;
              }
              let field = route.requestBodyFields[requestBodyField.key];
              if (!field) {
                return;
              }
              setTypeIsNumber(RoutePanelData.IsNumberType(field.type));
            }}
          />
        </Tooltip>
        <CloseOutlined
          onClick={() => {
            requestBodyOpt.remove(requestBodyField.name);
          }}
        />
      </Flex>
    </div>
  );
};

export default RequestFieldPanel;
