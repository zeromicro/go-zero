import React, {useState} from 'react';
import {LeftOutlined, GithubFilled} from '@ant-design/icons';
import {Button, Layout, Menu, theme, ConfigProvider, Space, Flex, Avatar, Typography} from 'antd';
import logo from '../assets/logo.svg'
import '../Base.css'
import './App.css';
import {menuItems} from "./_defaultProps";
import type {Locale} from "antd/lib/locale";
import zhCN from "antd/locale/zh_CN";
import enUS from "antd/locale/en_US";
import {ConverterIcon} from "../util/icon";

const {Title, Paragraph, Text, Link} = Typography;
const {Header, Sider, Content} = Layout;

const App: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const [localeZH, setLocaleZh] = useState(true);
    const [locale, setLocale] = useState(zhCN);

    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    const onLocaleClick = () => {
        const isZH = !localeZH
        setLocaleZh(isZH)
        if (isZH) {
            setLocale(zhCN)
        } else {
            setLocale(enUS)
        }
    }

    const onCollapsedClick = () => {
        setCollapsed(!collapsed)
    }
    const renderSiderFooter = () => {
        if (collapsed) {
            return <></>
        }
        if (localeZH){
            return <>
                <Button className="locale-btn" onClick={onLocaleClick}>中</Button>
                <ConverterIcon type={"icon-document"}/>
            </>
        }
        return <>
            <Button className="locale-btn" onClick={onLocaleClick}>EN</Button>
            <ConverterIcon type={"icon-document"}/>
        </>
    }
    return (
        <ConfigProvider
            locale={locale}
            theme={{
                token: {
                    colorPrimary: "#333333",
                    colorInfo: "#333333",
                },
            }}
        >
            <Layout hasSider>
                <Sider
                    trigger={null}
                    collapsible
                    collapsed={collapsed}
                    theme={"light"}
                    width={256}
                    collapsedWidth={66}
                    style={{
                        background: '#fafafa',
                    }}
                >
                    <Flex wrap
                          gap={10}
                          style={{
                              height: "60px",
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                              width: "100%",
                              padding: "10px 0",
                              background: '#fafafa',
                          }}
                    >
                        <Avatar src={<img src={logo} alt="avatar"/>} size={30}/>
                        {collapsed ?
                            <></> :
                            <Text ellipsis style={{
                                fontFamily: "阿里妈妈方圆体 VF Regular",
                                fontSize: 18,
                                color: "#333333",
                                paddingLeft: "10px",
                                paddingRight: "10px",
                            }}>goctl 网页端</Text>
                        }
                    </Flex>
                    <div style={{height: "1px", background: "#f1f1f1", margin: "0 20px"}}/>
                    <ConfigProvider
                        theme={{
                            components: {
                                Menu: {
                                    itemSelectedBg: "#ebebeb",
                                    itemSelectedColor: "rgba(0, 0, 0, 0.88)",
                                    itemMarginInline: 20,
                                    itemMarginBlock: 8,
                                    activeBarBorderWidth: 0,
                                    itemActiveBg: "#d0d0d0"
                                }
                            }
                        }}
                    >
                        <Menu
                            style={{
                                height: 'calc(100vh - 160px)',
                                overflowY: 'auto',
                                padding: "20px 0",
                                background: '#fafafa',
                            }}
                            theme="light"
                            mode="inline"
                            defaultSelectedKeys={['1']}
                            items={menuItems}
                        />
                    </ConfigProvider>
                    <Flex
                        justify="center"
                        align="center"
                        gap={10}
                        style={{height: "100px", background: "white", padding: "0 20px"}}>
                        <Link href="https://github.com/zeromicro/go-zero" target="_blank">
                            <GithubFilled style={{fontSize: 20}}/>
                        </Link>
                        {renderSiderFooter()}
                    </Flex>
                </Sider>

                <div style={{height: "100vh", width: "1px", background: "#ebebeb"}}/>

                {collapsed ?
                    <LeftOutlined className="collapse-button-uncollapsed" onClick={onCollapsedClick}/>
                    : <LeftOutlined className="collapse-button-collapsed" onClick={onCollapsedClick}/>}
                <Layout>
                    <Header style={{padding: 0, background: colorBgContainer}}>
                    </Header>
                    <Content
                        style={{
                            margin: '24px 16px',
                            padding: 24,
                            minHeight: 280,
                            background: colorBgContainer,
                            borderRadius: borderRadiusLG,
                        }}
                    >
                        Content
                    </Content>
                </Layout>
            </Layout>
        </ConfigProvider>
    )
        ;
};

export default App;