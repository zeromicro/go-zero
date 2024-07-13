import React, {useState, useEffect} from 'react';
import {
    Layout,
    theme,
    Flex,
    Typography,
    Form,
    Input,
    Space,
    Button,
    Card,
    Select, Col, Row, Switch, InputNumber
} from 'antd';
import {CloseOutlined} from '@ant-design/icons';
import '../../Base.css'
import './API.css'
import {useTranslation} from "react-i18next";
import {ConverterIcon} from "../../util/icon";
import {useNavigate} from "react-router-dom";

const {Text, Title, Link} = Typography;
const {Header, Sider, Content} = Layout;
const { TextArea } = Input;

const App: React.FC = () => {
    const navigate = useNavigate()
    const {t, i18n} = useTranslation();
    const [form] = Form.useForm();
    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    useEffect(() => {
    }, [])


    return (
        <Layout className="api">
            <Flex wrap className="api-container" gap={1}>
                <Flex className={"api-panel"} flex={1}>
                    <Form
                        name="basic"
                        autoComplete="off"
                        className={"api-panel-form"}
                        layout={"vertical"}
                        form={form}
                        initialValues={{items: [{}]}}
                    >
                        <Form.Item
                            label={t("formServiceTitle")}
                            name="serviceName"
                            rules={[{required: true, message: t("formServiceTips")}]}
                            className={"api-panel-form-item"}
                        >
                            <Input/>
                        </Form.Item>

                        <Form.List name="items">
                            {
                                (fields, {add, remove}) => (
                                    <div style={{display: 'flex', rowGap: 16, flexDirection: 'column'}}>
                                        {fields.map((field) => (
                                            <Card
                                                size="small"
                                                title={t("formRouteGroupTitle") + `${field.name + 1}`}
                                                key={field.key}
                                                extra={
                                                    <CloseOutlined
                                                        onClick={() => {
                                                            remove(field.name);
                                                        }}
                                                    />
                                                }
                                            >

                                                <Row gutter={16}>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formJwtTitle")}
                                                                   name={[field.name, 'jwt']}>
                                                            <Switch defaultChecked/>
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formPrefixTitle")}
                                                                   name={[field.name, 'prefix']}>
                                                            <Input prefix={"/"}/>
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formGroupTitle")}
                                                                   name={[field.name, 'group']}>
                                                            <Input/>
                                                        </Form.Item>
                                                    </Col>
                                                </Row>

                                                <Row gutter={16}>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formTimeoutTitle")}
                                                                   name={[field.name, 'timeout']}>
                                                            <InputNumber addonAfter="ms" defaultValue={2000}/>
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formMiddlewareTitle")}
                                                                   name={[field.name, 'prefix']}>
                                                            <Input/>
                                                        </Form.Item>
                                                    </Col>
                                                    <Col span={8}>
                                                        <Form.Item label={t("formMaxByteTitle")}
                                                                   name={[field.name, 'group']}>
                                                            <Input/>
                                                        </Form.Item>
                                                    </Col>
                                                </Row>

                                                {/* Nest Form.List */}
                                                <Form.Item label={t("formRouteListTitle")}>
                                                    <Form.List name={[field.name, 'list']}>
                                                        {(subFields, subOpt) => (
                                                            <div style={{
                                                                display: 'flex',
                                                                rowGap: 16,
                                                                flexDirection: 'column'
                                                            }}>

                                                                {subFields.map((subField) => (
                                                                    <Card
                                                                        size="small"
                                                                        title={t("formRouteTitle") + `${subField.name + 1}`}
                                                                        key={subField.key}
                                                                        extra={
                                                                            <CloseOutlined
                                                                                onClick={() => {
                                                                                    subOpt.remove(subField.name);
                                                                                }}
                                                                            />
                                                                        }
                                                                    >
                                                                        <Row gutter={16}>
                                                                            <Col span={12}>
                                                                                <Form.Item label={t("formMethodTitle")}
                                                                                           name={[subField.name, 'method']}>
                                                                                    <Select
                                                                                        defaultValue="POST"
                                                                                        options={[
                                                                                            {
                                                                                                value: 'GET',
                                                                                                label: 'GET'
                                                                                            },
                                                                                            {
                                                                                                value: 'HEAD',
                                                                                                label: 'HEAD'
                                                                                            },
                                                                                            {
                                                                                                value: 'POST',
                                                                                                label: 'POST'
                                                                                            },
                                                                                            {
                                                                                                value: 'PUT',
                                                                                                label: 'PUT'
                                                                                            },
                                                                                            {
                                                                                                value: 'PATCH',
                                                                                                label: 'PATCH'
                                                                                            },
                                                                                            {
                                                                                                value: 'DELETE',
                                                                                                label: 'DELETE'
                                                                                            },
                                                                                            {
                                                                                                value: 'CONNECT',
                                                                                                label: 'CONNECT'
                                                                                            },
                                                                                            {
                                                                                                value: 'OPTIONS',
                                                                                                label: 'OPTIONS'
                                                                                            },
                                                                                            {
                                                                                                value: 'TRACE',
                                                                                                label: 'TRACE'
                                                                                            },
                                                                                        ]}
                                                                                    />
                                                                                </Form.Item>
                                                                            </Col>
                                                                            <Col span={12}>
                                                                                <Form.Item
                                                                                    label={t("formContentTypeTitle")}
                                                                                    name={[subField.name, 'contentType']}>
                                                                                    <Select
                                                                                        defaultValue="application/json"
                                                                                        options={[
                                                                                            {
                                                                                                value: 'application/json',
                                                                                                label: 'application/json'
                                                                                            },
                                                                                            {
                                                                                                value: 'application/x-www-form-urlencoded',
                                                                                                label: 'application/x-www-form-urlencoded'
                                                                                            },
                                                                                        ]}
                                                                                    />
                                                                                </Form.Item>
                                                                            </Col>
                                                                        </Row>

                                                                        <Form.Item label={t("formPathTitle")}
                                                                                   name={[subField.name, 'path']}>
                                                                            <Input/>
                                                                        </Form.Item>

                                                                        {/*request body*/}
                                                                        <Form.Item label={t("formRequestBodyTitle")}>
                                                                            <Form.List
                                                                                name={[subField.name, 'requestBody']}>
                                                                                {(requestBodyFields, requestBodyOpt) => (
                                                                                    <div style={{
                                                                                        display: 'flex',
                                                                                        flexDirection: 'column'
                                                                                    }}>

                                                                                        {requestBodyFields.map((requestBodyField) => (
                                                                                            <Flex
                                                                                                key={requestBodyField.key}
                                                                                                gap={10}
                                                                                            >
                                                                                                <Form.Item
                                                                                                    label={t("formRequestBodyFieldNameTitle")}
                                                                                                    name={[requestBodyField.name, 'name']}
                                                                                                    style={{flex:1}}
                                                                                                >
                                                                                                    <Input/>
                                                                                                </Form.Item>
                                                                                                <Form.Item
                                                                                                    label={t("formRequestBodyFieldTypeTitle")}
                                                                                                    name={[requestBodyField.name, 'type']}
                                                                                                    style={{flex:1}}
                                                                                                >
                                                                                                    <Select
                                                                                                        defaultValue="string"
                                                                                                        options={[
                                                                                                            {
                                                                                                                value: 'bool',
                                                                                                                label: 'bool'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'uint8',
                                                                                                                label: 'uint8'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'uint16',
                                                                                                                label: 'uint16'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'uint32',
                                                                                                                label: 'uint32'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'uint64',
                                                                                                                label: 'uint64'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'int8',
                                                                                                                label: 'int8'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'int16',
                                                                                                                label: 'int16'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'int32',
                                                                                                                label: 'int32'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'int64',
                                                                                                                label: 'int64'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'float32',
                                                                                                                label: 'float32'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'float64',
                                                                                                                label: 'float64'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'string',
                                                                                                                label: 'string'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'int',
                                                                                                                label: 'int'
                                                                                                            },
                                                                                                            {
                                                                                                                value: 'uint',
                                                                                                                label: 'uint'
                                                                                                            },
                                                                                                        ]}
                                                                                                    />
                                                                                                </Form.Item>
                                                                                                <CloseOutlined
                                                                                                    onClick={() => {
                                                                                                        requestBodyOpt.remove(requestBodyField.name);
                                                                                                    }}
                                                                                                />
                                                                                            </Flex>
                                                                                        ))}
                                                                                        <Button type="dashed"
                                                                                                onClick={() => requestBodyOpt.add()}
                                                                                                block>
                                                                                            + {t("formRequestBodyFieldBtnAdd")}
                                                                                        </Button>

                                                                                    </div>

                                                                                )}
                                                                            </Form.List>
                                                                        </Form.Item>
                                                                    {/*  response body  */}
                                                                        <Form.Item label={t("formResponseBodyTitle")}>
                                                                            <TextArea
                                                                                autoSize={{ minRows: 3, maxRows: 5 }}
                                                                                placeholder={t("formResponseBodyPlaceholder")}/>
                                                                        </Form.Item>
                                                                    </Card>
                                                                ))}
                                                                <Button type="dashed" onClick={() => subOpt.add()}
                                                                        block>
                                                                    + {t("formButtonRouteAdd")}
                                                                </Button>

                                                            </div>

                                                        )}
                                                    </Form.List>
                                                </Form.Item>
                                            </Card>
                                        ))}

                                        <Button type="dashed" onClick={() => add()} block>
                                            + {t("formButtonRouteGroupAdd")}
                                        </Button>
                                    </div>
                                )}
                        </Form.List>

                    </Form>
                </Flex>
                <Flex className={"api-panel"} flex={1}>
                    right
                </Flex>
            </Flex>
        </Layout>
    )
};

export default App;