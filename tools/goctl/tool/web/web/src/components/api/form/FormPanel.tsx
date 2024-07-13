import React from "react";
import {Button, Col, Collapse, Flex, Form, Input, InputNumber, Layout, Row, Select, Switch, Typography} from "antd";
import {CloseOutlined} from "@ant-design/icons";
import {useTranslation} from "react-i18next";
import './FormPanel.css'
import RouteGroupPanel from './RouteGroupPanel'

const {Title} = Typography;
const {TextArea} = Input;

const FormPanel: React.FC = () => {
    const {t, i18n} = useTranslation();
    const [form] = Form.useForm();

    return (
        <Flex vertical className={"form-panel"} flex={1}>
            <Title level={4}>{t("builderPanelTitle")}</Title>
            <div className={"form-container-divider"}/>
            <Form
                name="basic"
                autoComplete="off"
                className={"form-panel-form"}
                layout={"vertical"}
                form={form}
                initialValues={
                    {
                        serviceName: "",
                        routeGroups: [{}]
                    }
                }
            >
                <Form.Item
                    label={t("formServiceTitle")}
                    name="serviceName"
                    rules={[{required: true, message: t("formServiceTips")}]}
                    className={"form-panel-form-item"}
                >
                    <Input placeholder={t("formServicePlaceholder")}/>
                </Form.Item>

                <Form.List name="routeGroups">
                    {
                        (routeGroupFields, routeGroupOperation) => (
                            <div style={{display: 'flex', rowGap: 16, flexDirection: 'column'}}>
                                {routeGroupFields.map((routeGroupField) => (
                                    <Collapse
                                        defaultActiveKey={[routeGroupField.key]}
                                        items={
                                            [
                                                {
                                                    key: routeGroupField.key,
                                                    label: t("formRouteGroupTitle") + `${routeGroupField.name + 1}`,
                                                    children: <RouteGroupPanel routeGroupField={routeGroupField} form={form}/>,
                                                    extra: <CloseOutlined
                                                        onClick={() => {
                                                            routeGroupOperation.remove(routeGroupField.name);
                                                        }}
                                                    />
                                                }
                                            ]
                                        }
                                    />
                                ))}

                                <Button type="dashed" onClick={() => routeGroupOperation.add()} block>
                                    + {t("formButtonRouteGroupAdd")}
                                </Button>
                            </div>
                        )}
                </Form.List>

            </Form>
        </Flex>
    )
}

export default FormPanel;