import React from "react";
import {Button, Col, Collapse, Flex, Form, Input, InputNumber, Row, Select, Switch} from "antd";
import {CloseOutlined} from "@ant-design/icons";
import {FormListFieldData} from "antd/es/form/FormList";
import {useTranslation} from "react-i18next";

const {TextArea} = Input;

interface RouteGroupOptionPanelProps {
    routeGroupField: FormListFieldData
}

const RouteGroupOptionPanel: React.FC<RouteGroupOptionPanelProps & React.RefAttributes<HTMLDivElement>> = (props) => {
    const {t, i18n} = useTranslation();
    const routeGroupField = props.routeGroupField
    return (
        <div>
            <Row gutter={16}>
                <Col span={8}>
                    <Form.Item label={t("formJwtTitle")}
                               name={[routeGroupField.name, 'jwt']}>
                        <Switch defaultChecked/>
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item label={t("formPrefixTitle")}
                               name={[routeGroupField.name, 'prefix']}>
                        <Input prefix={"/"}
                               placeholder={t("formPrefixPlaceholder")}/>
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item label={t("formGroupTitle")}
                               name={[routeGroupField.name, 'group']}>
                        <Input/>
                    </Form.Item>
                </Col>
            </Row>

            <Row gutter={16}>
                <Col span={8}>
                    <Form.Item label={t("formTimeoutTitle")}
                               name={[routeGroupField.name, 'timeout']}>
                        <InputNumber addonAfter="ms"
                                     defaultValue={2000}/>
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item label={t("formMiddlewareTitle")}
                               name={[routeGroupField.name, 'prefix']}>
                        <Input/>
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item label={t("formMaxByteTitle")}
                               name={[routeGroupField.name, 'group']}>
                        <Input/>
                    </Form.Item>
                </Col>
            </Row>
        </div>
    )
}

export default RouteGroupOptionPanel;