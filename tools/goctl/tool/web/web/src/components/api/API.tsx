import React, {useState, useEffect} from 'react';
import {
    Layout,
    Form, Row, Col
} from 'antd';
import '../../Base.css'
import './API.css'
import FormPanel from "./form/FormPanel";
import CodePanel from "./form/CodePanel";

const App: React.FC = () => {
    const [code, setCode] = useState("")
    return (
        <Layout className="api">
            <Row wrap className="api-container" gutter={1}>
                <Col span={12} className={"api-form-panel"}>
                    <FormPanel onBuild={(data) => {
                        const js = JSON.stringify(data)
                        setCode(js)
                    }}/>
                </Col>
                <Col span={12} className={"api-code-panel"}>
                    <CodePanel
                        onChange={(code) => {
                            setCode(code)
                        }}
                        value={code}/>
                </Col>
            </Row>
        </Layout>
    )
};

export default App;