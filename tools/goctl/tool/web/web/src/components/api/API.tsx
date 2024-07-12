import React, {useState, useEffect} from 'react';
import {
    Button,
    Layout,
    Menu,
    theme,
    ConfigProvider,
    Flex,
    Avatar,
    Typography,
    Tag,
    Dropdown,
    Space,
    MenuProps
} from 'antd';
import '../../Base.css'
import './API.css'
import {useTranslation} from "react-i18next";
import {ConverterIcon} from "../../util/icon";
import {useNavigate} from "react-router-dom";

const {Text, Title, Link} = Typography;
const {Header, Sider, Content} = Layout;

const App: React.FC = () => {
    const navigate = useNavigate()
    const {t, i18n} = useTranslation();
    const {
        token: {colorBgContainer, borderRadiusLG},
    } = theme.useToken();

    useEffect(() => {
    }, [])


    return (
        <Layout className="api">
            <Flex wrap className="api-container" gap={1}>
                <Flex className={"api-panel"} flex={1}>
                    left
                </Flex>
                <Flex className={"api-panel"} flex={1}>
                    right
                </Flex>
            </Flex>
        </Layout>
    )
};

export default App;