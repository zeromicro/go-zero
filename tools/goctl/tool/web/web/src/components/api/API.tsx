import React, {useState, useEffect} from 'react';
import {
    Layout,
    Flex,
    Typography,
    Form,
    Input,
} from 'antd';
import '../../Base.css'
import './API.css'
import {useTranslation} from "react-i18next";
import {useNavigate} from "react-router-dom";
import FormPanel from "./form/FormPanel";
import {langs} from "@uiw/codemirror-extensions-langs";
import CodeMirror, {EditorView} from "@uiw/react-codemirror";
import {githubDark} from "@uiw/codemirror-theme-github";
import CodePanel from "./form/CodePanel";

const {Text, Title, Link} = Typography;
const {Header, Sider, Content} = Layout;
const {TextArea} = Input;

const App: React.FC = () => {
    const navigate = useNavigate()
    const {t, i18n} = useTranslation();
    const [form] = Form.useForm();

    useEffect(() => {
    }, [])

    return (
        <Layout className="api">
            <Flex wrap className="api-container" gap={1}>
                <FormPanel/>
                <CodePanel/>
            </Flex>
        </Layout>
    )
};

export default App;