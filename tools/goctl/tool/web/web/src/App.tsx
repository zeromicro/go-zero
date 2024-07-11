import React from 'react';
import {Layout, Typography, ConfigProvider, theme} from 'antd';
import {Col, Row} from 'antd';
import './Base.css'
import {ConverterIcon} from './util/icon'
import {Outlet} from "react-router-dom";
import {useNavigate} from "react-router-dom";

const {Text, Link} = Typography;
const {Header} = Layout;

function App() {
    const navigate = useNavigate()
    return (
        <>
            <ConfigProvider
                theme={{
                    algorithm: theme.compactAlgorithm,
                    token: {
                        colorPrimary: "#1890ff",
                        colorInfo: "#54aeff"
                    },
                    components: {
                        Layout: {
                            headerHeight: 50,
                            headerBg: 'rgb(22, 119, 255)'
                        },
                        Menu: {
                            darkItemBg: 'rgb(22, 119, 255)',
                            darkItemSelectedBg: 'rgb(3, 91, 215)',
                            itemHoverBg: 'rgb(230, 244, 255)',
                            activeBarWidth: 2,
                            iconSize: 20,
                            collapsedIconSize: 26,
                            fontSize: 14
                        },
                        Alert: {
                            colorInfo: 'rgb(22, 119, 255)'
                        },
                        Message: {
                            colorInfo: 'rgb(22, 119, 255)'
                        },
                        Notification: {
                            colorPrimary: 'rgb(22, 119, 255)'
                        },
                        Progress: {
                            algorithm: true
                        },
                        Card: {
                            headerBg: "#1890ff",
                            colorText: "#0062bd",
                            colorTextHeading: "white",
                            fontFamily: "阿里巴巴普惠体 2.0 35 Thin",
                            headerFontSize: 16
                        }
                    }
                }}
            >
                <Layout style={{height: '100vh', overflowY: "hidden"}}>
                    <Header style={{display: 'flex', alignItems: 'center', height: '50px', background: "#1890ff"}}>
                        <Row style={{width: '100%', display: 'flex', justifyContent: "space-between"}}>
                            <Col style={{display: 'flex', alignItems: 'center', cursor: "pointer"}} onClick={() => {
                                navigate("/")
                            }}>
                                <ConverterIcon type="icon-jiandao" style={{
                                    fontSize: 30,
                                }}/>
                                <Text style={{
                                    color: "white",
                                    fontSize: 22,
                                    marginLeft: '10px',
                                    letterSpacing: 1,
                                    fontFamily: "钉钉进步体 Regular",
                                }}>converter</Text>
                            </Col>
                            <Col>
                                <Link href="https://github.com/kesonan/converter" target="_blank"
                                      style={{color: "white", width: "100%"}}>
                                    <Row style={{
                                        height: '100%',
                                        display: "flex",
                                        alignItems: 'center',
                                        justifyContent: "flex-end"
                                    }}>
                                        <Col style={{display: 'flex', alignItems: 'center'}}>
                                            <img alt="GitHub Repo stars"
                                                 src="https://img.shields.io/github/stars/kesonan/converter"/>
                                        </Col>
                                        <Col style={{marginLeft: '10px'}}>
                                            <Text style={{
                                                color: "white",
                                                fontSize: 14,
                                                letterSpacing: 1,
                                                fontFamily: "阿里巴巴普惠体 2.0 35 Thin",
                                                fontWeight: 400
                                            }}> GitHub</Text>
                                        </Col>
                                    </Row>
                                </Link>
                            </Col>
                        </Row>
                    </Header>
                    <Outlet/>
                </Layout>
            </ConfigProvider>
        </>
    );
}

export default App;
