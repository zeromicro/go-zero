import React, {useState} from "react";
import {
    Button,
    Col,
    Collapse,
    Flex,
    Form,
    Input,
    Row,
    Select,
    Modal,
    notification,
    Tooltip,
    type MenuProps, Switch
} from "antd";
import {CloseOutlined, FullscreenOutlined} from "@ant-design/icons";
import {FormListFieldData} from "antd/es/form/FormList";
import {useTranslation} from "react-i18next";
import {RoutePanelData, Method, ContentType, GolangType} from "./_defaultProps";
import CodeMirror, {EditorView} from '@uiw/react-codemirror';
import {githubLight} from '@uiw/codemirror-theme-github';
import {langs} from '@uiw/codemirror-extensions-langs';
import type {FormInstance} from "antd/es/form/hooks/useForm";

const {TextArea} = Input;

interface RoutePanelProps {
    routeGroupField: FormListFieldData
    form: FormInstance
}


const RoutePanel: React.FC<RoutePanelProps & React.RefAttributes<HTMLDivElement>> = (props) => {
    const {t, i18n} = useTranslation();
    const routeGroupField = props.routeGroupField
    const form = props.form
    const [initRequestValues, setInitRequestValues] = useState([]);
    const [requestBodyModalOpen, setRequestBodyModalOpen] = useState(false);
    const [responseBodyModalOpen, setResponseBodyModalOpen] = useState(false);
    const [requestBodyParseCode, setRequestBodyParseCode] = useState('');
    const [responseCode, setResponseCode] = useState('');
    const [api, contextHolder] = notification.useNotification();
    const [showImportButton, setShowImportButton] = useState(true);

    const canChowImportButton = (routeIdx: number) => {
        const routeGroups = form.getFieldValue(`routeGroups`)
        if (!routeGroups) {
            setShowImportButton(true)
            return
        }

        if (routeGroups.length <= routeGroupField.key) {
            setShowImportButton(true)
            return
        }

        const routeGroup = routeGroups[routeGroupField.key]

        if (!routeGroup) {
            setShowImportButton(true)
            return
        }
        if (!routeGroup.routes) {
            setShowImportButton(true)
            return
        }

        if (routeGroup.routes.length <= routeIdx) {
            setShowImportButton(true)
            return
        }

        const route = routeGroup.routes[routeIdx]
        if (!route) {
            setShowImportButton(true)
            return
        }
        if (!route.requestBodyFields) {
            setShowImportButton(true)
            return
        }


        setShowImportButton(route.requestBodyFields.length === 0)
    }
    return (
        <div>
            {contextHolder}
            {/*request body import modal*/}
            <Modal
                title={t("formRequestBodyFieldBtnImport")}
                centered
                open={requestBodyModalOpen}
                maskClosable={false}
                keyboard={false}
                closable={false}
                destroyOnClose
                onOk={() => {
                    try {
                        const obj = JSON.parse(requestBodyParseCode)
                        if (Array.isArray(obj)) {
                            api.error({
                                message: t("tipsInvalidJSONArray")
                            })
                            return
                        }

                        // todo: 从后段解析数据
                        setRequestBodyModalOpen(false)
                    } catch (err) {
                        api.error({
                            message: t("tipsInvalidJSON") + ": " + err
                        })
                        return
                    }
                }}
                onCancel={() => setRequestBodyModalOpen(false)}
                width={1000}
                cancelText={t("formRequestBodyModalCancel")}
                okText={t("formRequestBodyModalConfirm")}
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
                    onChange={(code) => {
                        setRequestBodyParseCode(code)
                    }}
                />
            </Modal>
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
                                                    {/*<Col span={4}>*/}
                                                    {/*    <Form.Item*/}
                                                    {/*        label={t("formMethodTitle")}*/}
                                                    {/*        name={[routeField.name, 'method']}>*/}
                                                    {/*        <Select*/}
                                                    {/*            options={RoutePanelData.MethodOptions}*/}
                                                    {/*        />*/}
                                                    {/*    </Form.Item>*/}
                                                    {/*</Col>*/}
                                                    <Col span={16}>
                                                        <Form.Item
                                                            label={t("formPathTitle")}
                                                            name={[routeField.name, 'path']}>
                                                            <Input addonBefore={<div>
                                                                <Form.Item
                                                                    noStyle
                                                                    name={[routeField.name, 'method']}>
                                                                    <Select
                                                                        style={{width: 100}}
                                                                        defaultValue={Method.POST}
                                                                        options={RoutePanelData.MethodOptions}
                                                                    />
                                                                </Form.Item>
                                                            </div>}/>
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={8}>
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
                                                {/*request body*/}
                                                <Form.Item
                                                    label={t("formRequestBodyTitle")}>
                                                    <Form.List
                                                        initialValue={initRequestValues}
                                                        name={[routeField.name, 'requestBodyFields']}>
                                                        {(requestBodyFields, requestBodyOpt) => (
                                                            <div
                                                                style={{
                                                                    display: 'flex',
                                                                    flexDirection: 'column',
                                                                }}>

                                                                {requestBodyFields.map((requestBodyField) => (
                                                                    <Flex
                                                                        key={requestBodyField.key}
                                                                        gap={10}
                                                                        wrap
                                                                    >
                                                                        <Form.Item
                                                                            label={t("formRequestBodyFieldNameTitle")}
                                                                            name={[requestBodyField.name, 'name']}
                                                                            style={{flex: 1}}
                                                                            tooltip={t("formRequestBodyFieldNameTooltip")}
                                                                        >
                                                                            <Input/>
                                                                        </Form.Item>
                                                                        <Form.Item
                                                                            label={t("formRequestBodyFieldTypeTitle")}
                                                                            name={[requestBodyField.name, 'type']}
                                                                            style={{flex: 1}}
                                                                        >
                                                                            <Select
                                                                                options={RoutePanelData.GolangTypeOptions}
                                                                            />
                                                                        </Form.Item>
                                                                        <Form.Item
                                                                            label={t("formRequestBodyFieldOptionalTitle")}
                                                                            name={[requestBodyField.name, 'optional']}
                                                                        >
                                                                            <Switch/>

                                                                        </Form.Item>
                                                                        <CloseOutlined
                                                                            onClick={() => {
                                                                                requestBodyOpt.remove(requestBodyField.name);
                                                                                canChowImportButton(routeField.key)
                                                                            }}
                                                                        />
                                                                    </Flex>
                                                                ))}
                                                                {showImportButton ? <Button
                                                                    style={{marginBottom: 16}}
                                                                    type="dashed"
                                                                    onClick={() => setRequestBodyModalOpen(true)}
                                                                    block>
                                                                    🔍 {t("formRequestBodyFieldBtnImport")}
                                                                </Button> : <></>
                                                                }
                                                                <Button
                                                                    type="dashed"
                                                                    onClick={() => {
                                                                        requestBodyOpt.add()
                                                                        canChowImportButton(routeField.key)

                                                                    }}
                                                                    block>
                                                                    + {t("formRequestBodyFieldBtnAdd")}
                                                                </Button>

                                                            </div>

                                                        )}
                                                    </Form.List>
                                                </Form.Item>
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
                                                            const routeGroups = form.getFieldValue("routeGroups")
                                                            if (!routeGroups) {
                                                                return
                                                            }
                                                            const routeGroup = routeGroups[routeGroupField.key]
                                                            if (!routeGroup) {
                                                                return
                                                            }
                                                            const routes = routeGroup.routes
                                                            if (!routes) {
                                                                return;
                                                            }
                                                            const route = routes[routeField.key]
                                                            if (!route) {
                                                                return
                                                            }
                                                            route.responseBody = code
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