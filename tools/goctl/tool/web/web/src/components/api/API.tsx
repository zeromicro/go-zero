import React, {useState, useEffect} from 'react';
import {
    Layout,
    Flex,
    Form,
} from 'antd';
import '../../Base.css'
import './API.css'
import {useTranslation} from "react-i18next";
import {useNavigate} from "react-router-dom";
import FormPanel from "./form/FormPanel";
import CodePanel from "./form/CodePanel";
import {codePlaceholder} from './_defaultProps'

const App: React.FC = () => {
    const navigate = useNavigate()
    const {t, i18n} = useTranslation();
    const [form] = Form.useForm();
    const [code, setCode] = useState(codePlaceholder())

    useEffect(() => {
    }, [])

    return (
        <Layout className="api">
            <Flex wrap className="api-container" gap={1}>
                <FormPanel/>
                <CodePanel value={code}/>
            </Flex>
        </Layout>
    )
};

export default App;