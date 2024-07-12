import React, {useState, useEffect} from 'react';
import {LeftOutlined, GithubFilled} from '@ant-design/icons';
import {Button, Layout, Menu, theme, ConfigProvider, Flex, Avatar, Typography, Breadcrumb} from 'antd';
import '../../Base.css'
import './App.css';
import {menuItems} from "./_defaultProps";
import zhCN from "antd/locale/zh_CN";
import enUS from "antd/locale/en_US";
import {ConverterIcon} from "../../util/icon";
import {useTranslation} from "react-i18next";
import {Outlet} from "react-router-dom";
import {useNavigate} from "react-router-dom";
import {MenuInfo} from "rc-menu/lib/interface";
import {ItemType} from "antd/es/breadcrumb/Breadcrumb";
import logo from "../../assets/logo.svg"

const {Text, Link} = Typography;
const {Header, Sider, Content} = Layout;

const App: React.FC = () => {
    const navigate = useNavigate()
    const {t, i18n} = useTranslation();
    const [collapsed, setCollapsed] = useState(false);
    const [localeZH, setLocaleZh] = useState(false);
    const [locale, setLocale] = useState(zhCN);
    const [breadcrumbItems, setBreadcrumbItems] = useState<ItemType[]>([{title: t("welcome")}]);

    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    useEffect(() => {
        setLocaleZh(i18n.language != "zh")
    }, [])

    const onLocaleClick = () => {
        const isZH = !localeZH
        setLocaleZh(isZH)
        if (isZH) {
            i18n.changeLanguage("en")
            setLocale(zhCN)
        } else {
            i18n.changeLanguage("zh")
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
        if (localeZH) {
            return <>
                <Button className="locale-btn" onClick={onLocaleClick}>中</Button>
                <Link href="https://go-zero.dev" target="_blank">
                    <ConverterIcon type={"icon-document"} className="sider-footer-icon"/>
                </Link>
            </>
        }
        return <>
            <Button className="locale-btn" onClick={onLocaleClick}>EN</Button>
            <Link href="https://go-zero.dev" target="_blank">
                <ConverterIcon type={"icon-document"} className="sider-footer-icon"/>
            </Link>
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
                            <Text
                                ellipsis
                                className={"logo-text-gradient"}
                            >{t("logoText")}</Text>
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
                                    itemActiveBg: "#d0d0d0",
                                    subMenuItemBg: "#fafafa"
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
                            items={menuItems(t)}
                            defaultOpenKeys={["api"]}
                            onClick={(info: MenuInfo) => {
                                let breadcrumbItems: ItemType[] = []
                                if (info.key !== 'welcome') {
                                    breadcrumbItems.push({
                                        title: t("welcome"),
                                    })
                                }
                                const reverseArray = info.keyPath.reverse()
                                reverseArray.forEach((val: string) => {
                                    breadcrumbItems.push({
                                        title: t(val)
                                    })
                                })
                                setBreadcrumbItems(breadcrumbItems)
                                const path = reverseArray.join("/")
                                navigate(path)
                            }}
                        />
                    </ConfigProvider>
                    <Flex
                        vertical
                        justify="center"
                        align="center"
                        gap={10}
                        style={{height: "100px", background: "white", padding: "0 20px"}}
                    >
                        <Flex
                            justify="center"
                            align="center"
                            gap={20}
                        >
                            <Link href="https://github.com/zeromicro/go-zero" target="_blank">
                                <GithubFilled className="sider-footer-icon"/>
                            </Link>
                            {renderSiderFooter()}
                        </Flex>
                        {collapsed ? <></> : <Text style={{color: "#333333", fontSize: 12}} ellipsis>
                            go-zero ©{new Date().getFullYear()} Created by zeromicro
                        </Text>
                        }
                    </Flex>
                </Sider>

                <div style={{height: "100vh", width: "1px", background: "#ebebeb"}}/>

                {collapsed ?
                    <LeftOutlined className="collapse-button-uncollapsed" onClick={onCollapsedClick}/>
                    : <LeftOutlined className="collapse-button-collapsed" onClick={onCollapsedClick}/>}
                <Layout>
                    <Breadcrumb
                        items={breadcrumbItems}
                        style={{
                            marginTop: 22,
                            marginLeft: 16
                        }}
                    />
                    <Content
                        style={{
                            margin: '24px 16px',
                            minHeight: 280,
                            background: colorBgContainer,
                            borderRadius: borderRadiusLG,
                        }}
                    >
                        <Outlet/>
                    </Content>

                </Layout>
            </Layout>
        </ConfigProvider>
    )
        ;
};

export default App;