import React, {useState} from "react";
import {
    Button,
    Col,
    Collapse,
    Form,
    Input,
    Row,
    Select,
    Modal,
    notification,
    Tooltip,
} from "antd";
import {CloseOutlined, FullscreenOutlined} from "@ant-design/icons";
import {FormListFieldData} from "antd/es/form/FormList";
import {useTranslation} from "react-i18next";
import {RoutePanelData, Method, ContentType} from "./_defaultProps";
import CodeMirror, {EditorView} from '@uiw/react-codemirror';
import {githubLight} from '@uiw/codemirror-theme-github';
import {langs} from '@uiw/codemirror-extensions-langs';
import type {FormInstance} from "antd/es/form/hooks/useForm";
import RequestBodyPanel from "./RequestBodyPanel";

interface RoutePanelProps {
    routeGroupField: FormListFieldData
    form: FormInstance
}


const RoutePanel: React.FC<RoutePanelProps & React.RefAttributes<HTMLDivElement>> = (props) => {
    const {t} = useTranslation();
    const routeGroupField = props.routeGroupField
    const form = props.form
    const [initRequestValues, setInitRequestValues] = useState([]);
    const [responseBodyModalOpen, setResponseBodyModalOpen] = useState(false);
    const [responseCode, setResponseCode] = useState('');
    const [api, contextHolder] = notification.useNotification();


    return (
        <div>
            {contextHolder}
            {/*response body editor*/}
            <Modal
                title={t("formResponseBodyModelTitle")}
                centered
                open={responseBodyModalOpen}
                maskClosable={false}
                keyboard={false}
                closable={false}
                destroyOnClose
                onOk={() => {
                    setResponseBodyModalOpen(false)
                }}
                onCancel={() => setResponseBodyModalOpen(false)}
                width={1000}
                cancelText={t("formResponseBodyModalCancel")}
                okText={t("formResponseBodyModalConfirm")}
            >
                <CodeMirror
                    style={{marginTop: 10, overflow: "auto"}}
                    extensions={[langs.json(), EditorView.theme({
                        "&.cm-focused": {
                            outline: "none",
                        },
                    })]}
                    theme={githubLight}
                    height={'70vh'}
                    value={responseCode}
                    onChange={(code) => {
                        setResponseCode(code)
                    }}
                />
            </Modal>
            <Form.Item label={t("formRouteListTitle")}>
                <Form.List
                    name={[routeGroupField.name, 'routes']}>
                    {(routeFields, routeOpt) => (
                        <div style={{
                            display: 'flex',
                            rowGap: 16,
                            flexDirection: 'column'
                        }}>

                            {routeFields.map((routeField) => (
                                <Collapse
                                    defaultActiveKey={[routeField.key]}
                                    items={[
                                        {
                                            key: routeField.key,
                                            label: t("formRouteTitle") + `${routeField.name + 1}`,
                                            children: <div>
                                                <Row gutter={16}>
                                                    <Col span={12}>
                                                        <Form.Item
                                                            label={t("formMethodTitle")}
                                                            name={[routeField.name, 'method']}>
                                                            <Select
                                                                defaultValue={Method.POST}
                                                                options={RoutePanelData.MethodOptions}
                                                            />
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={12}>
                                                        <Form.Item
                                                            label={t("formContentTypeTitle")}
                                                            name={[routeField.name, 'contentType']}>
                                                            <Select
                                                                defaultValue={ContentType.ApplicationJson}
                                                                options={RoutePanelData.ContentTypeOptions}
                                                            />
                                                        </Form.Item>
                                                    </Col>
                                                </Row>

                                                <Form.Item
                                                    label={t("formPathTitle")}
                                                    name={[routeField.name, 'path']}>
                                                    <Input/>
                                                </Form.Item>

                                                {/*request body*/}
                                                <RequestBodyPanel
                                                    form={form}
                                                    routeField={routeField}
                                                    routeGroupField={routeGroupField}
                                                />
                                                {/*  response body  */}
                                                <Form.Item
                                                    label={t("formResponseBodyTitle")}
                                                    name={[routeField.name, 'responseBody']}>
                                                    <span style={{
                                                        position: "absolute",
                                                        top: -30,
                                                        right: 0,
                                                        zIndex: 1000
                                                    }}>
                                                        <Tooltip title={t("tooltipFullScreen")}>
                                                            <FullscreenOutlined
                                                                style={{cursor: "pointer"}}
                                                                onClick={() => {
                                                                    setResponseBodyModalOpen(true)
                                                                }}
                                                            />
                                                        </Tooltip>
                                                    </span>
                                                    <CodeMirror
                                                        style={{overflow: "scroll", minHeight: 100, maxHeight: 200}}
                                                        extensions={[langs.json(), EditorView.theme({
                                                            "&.cm-focused": {
                                                                outline: "none",
                                                            },
                                                        })]}
                                                        value={responseCode}
                                                        placeholder={t("formResponseBodyPlaceholder")}
                                                        theme={githubLight}
                                                        onChange={(code) => {
                                                            setResponseCode(code)
                                                        }}
                                                    />
                                                </Form.Item>
                                            </div>,
                                            extra: <CloseOutlined
                                                onClick={() => {
                                                    routeOpt.remove(routeField.name);
                                                }}
                                            />
                                        }
                                    ]}
                                />
                            ))}
                            <Button type="dashed"
                                    onClick={() => routeOpt.add()}
                                    block>
                                + {t("formButtonRouteAdd")}
                            </Button>

                        </div>

                    )}
                </Form.List>
            </Form.Item>
        </div>
    )
}

export default RoutePanel;