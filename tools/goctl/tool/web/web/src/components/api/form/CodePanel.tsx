import React, {useState} from "react";
import {Button, Flex, Typography, message} from "antd";
import {langs} from "@uiw/codemirror-extensions-langs";
import CodeMirror, {EditorView} from "@uiw/react-codemirror";
import {githubLight} from "@uiw/codemirror-theme-github";
import "./CodePanel.css"
import {useTranslation} from "react-i18next";
import {CopyOutlined} from '@ant-design/icons';
import {CopyToClipboard} from 'react-copy-to-clipboard';

const {Title} = Typography;

interface CodePanelProps {
    value: string
}


const CodePanel: React.FC<CodePanelProps & React.RefAttributes<HTMLDivElement>> = (props) => {
    const {t, i18n} = useTranslation();
    const [api, contextHolder] = message.useMessage();

    const onCopy = () => {
        api.open({
            type: 'success',
            content: 'Copied to clipboard',
        });
    }

    return (
        <Flex vertical className={"code-panel"} flex={1}>
            {contextHolder}
            <Flex justify={"space-between"} align={"center"}>
                <Title level={4}>{t("apiPanelTitle")}</Title>
                <CopyToClipboard text={props.value} onCopy={onCopy}>
                    <Button size={"middle"}>
                        <CopyOutlined/> {t("btnCopy")}
                    </Button>
                </CopyToClipboard>
            </Flex>
            <div className={"code-container-divider"}/>
            <CodeMirror
                style={{overflowY: "auto"}}
                extensions={[langs.go(), EditorView.theme({
                    "&.cm-focused": {
                        outline: "none",
                    },
                })]}
                value={props.value}
                editable={false}
                readOnly
                theme={githubLight}
                onChange={(code) => {

                }}
            />
        </Flex>
    )
}
export default CodePanel;