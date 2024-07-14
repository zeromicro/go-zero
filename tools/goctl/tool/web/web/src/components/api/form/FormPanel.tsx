import React from "react";
import {
    Button,
    Collapse,
    Flex,
    Form,
    Input,
    message,
    Typography
} from "antd";
import {CloseOutlined, CopyOutlined} from "@ant-design/icons";
import {useTranslation} from "react-i18next";
import './FormPanel.css'
import RouteGroupPanel from './RouteGroupPanel'
import {ConverterIcon} from "../../../util/icon";
import {FormInstance} from "rc-field-form/es/interface";

const {Title} = Typography;

interface FormPanelProps {
    onBuild?: (data: FormInstance) => void
}

const FormPanel: React.FC<FormPanelProps> = (props) => {
    const [api, contextHolder] = message.useMessage();
    const {t, i18n} = useTranslation();
    const [form] = Form.useForm();
    const onBuild = () => {
        const obj = form.getFieldsValue()
        if(props.onBuild) {
            props.onBuild(obj)
        }
        api.open({
            type: 'success',
            content: '转换成功',
        });
    }
    return (
        <Flex vertical className={"form-panel"} flex={1}>
            {contextHolder}
            <Flex justify={"space-between"} align={"center"} className={"form-container-header"}>
                <Title level={4}>{t("builderPanelTitle")}</Title>
                <Button size={"middle"} onClick={onBuild} type={"primary"}>
                    <ConverterIcon type={"icon-terminal"} className="welcome-start-icon"/>{t("btnBuild")}
                </Button>
            </Flex>
            <div className={"form-container-divider"}/>
            <Form
                name="basic"
                autoComplete="off"
                className={"form-panel-form"}
                layout={"vertical"}
                form={form}
                initialValues={
                    {
                        serviceName: "demo",
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
                    <Input placeholder={t("formServicePlaceholder")} allowClear/>
                </Form.Item>

                <Form.List
                    name="routeGroups">
                    {
                        (routeGroupFields, routeGroupOperation) => (
                            <div style={{display: 'flex', rowGap: 16, flexDirection: 'column'}}>
                                {routeGroupFields.map((routeGroupField) => (
                                    <Collapse
                                        key={routeGroupField.key}
                                        defaultActiveKey={[routeGroupField.key]}
                                        items={
                                            [
                                                {
                                                    key: routeGroupField.key,
                                                    label: t("formRouteGroupTitle") + `${routeGroupField.name + 1}`,
                                                    children: <RouteGroupPanel routeGroupField={routeGroupField}
                                                                               form={form}/>,
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