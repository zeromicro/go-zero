import React, {useState} from 'react';
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    UploadOutlined,
    UserOutlined,
    VideoCameraOutlined,
} from '@ant-design/icons';
import {Button, Layout, Menu, theme, ConfigProvider, Space, Flex, Avatar, Typography} from 'antd';
import logo from './assets/logo.svg'
import './Base.css'

const {Title, Paragraph, Text, Link} = Typography;
const {Header, Sider, Content} = Layout;

const App: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    return (
        <Layout>
            <ConfigProvider
                theme={{
                    components: {
                        Layout: {
                            siderBg: "#fafafa",
                        },
                        Menu: {
                            itemSelectedBg: "#ebebeb",
                            itemSelectedColor: "rgba(0, 0, 0, 0.88)",
                        }
                    }
                }}
            >
                <Sider
                    trigger={null}
                    collapsible
                    collapsed={collapsed}
                    theme={"light"}
                    width={256}
                    collapsedWidth={66}
                >
                    <Flex wrap
                          gap={10}
                          style={{
                              minHeight: "60px",
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                              width: "100%",
                              padding: "10px 0"
                          }}
                    >
                        <Avatar src={<img src={logo} alt="avatar"/>} size={30}/>
                        <Text ellipsis>Goctl Web</Text>
                    </Flex>
                    <div style={{height: "1px", background: "#f1f1f1"}}/>

                    <Menu
                        style={{height: 'calc(100vh - 240px)', overflowY: 'auto', padding: "20px 0"}}
                        theme="light"
                        mode="inline"
                        defaultSelectedKeys={['1']}
                        items={[
                            {
                                key: '1',
                                icon: <UserOutlined/>,
                                label: 'nav 1',
                                children: [
                                    {
                                        key: '1',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '2',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '3',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                    {
                                        key: '4',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '5',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '6',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                ]
                            },
                            {
                                key: '2',
                                icon: <VideoCameraOutlined/>,
                                label: 'nav 2',
                                children: [
                                    {
                                        key: '1',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '2',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '3',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                    {
                                        key: '4',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '5',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '6',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                ]
                            },
                            {
                                key: '3',
                                icon: <UploadOutlined/>,
                                label: 'nav 3',
                                children: [
                                    {
                                        key: '1',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '2',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '3',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                    {
                                        key: '4',
                                        icon: <UserOutlined/>,
                                        label: 'nav 1',
                                    },
                                    {
                                        key: '5',
                                        icon: <VideoCameraOutlined/>,
                                        label: 'nav 2',
                                    },
                                    {
                                        key: '6',
                                        icon: <UploadOutlined/>,
                                        label: 'nav 3',
                                    },
                                ]
                            },
                        ]}
                    />


                    <div style={{height: "100px", background: "black"}}/>
                </Sider>
            </ConfigProvider>
            <div style={{height: "100vh", width: "1px", background: "#f1f1f1"}}/>
            <Layout>
                <Header style={{padding: 0, background: colorBgContainer}}>
                    <Button
                        type="text"
                        icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                        onClick={() => setCollapsed(!collapsed)}
                        style={{
                            fontSize: '16px',
                            width: 26,
                            height: 26,
                            marginLeft: -15,
                            borderRadius: 26,
                            background: "white",
                            boxShadow: "1px 1px 10px #f1f1f1",
                            borderStyle: "solid",
                            borderWidth: 1
                        }}
                    />
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
    );
};

export default App;