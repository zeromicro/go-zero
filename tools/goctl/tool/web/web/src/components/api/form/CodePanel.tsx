import React from "react";
import {Flex, Typography,} from "antd";
import {langs} from "@uiw/codemirror-extensions-langs";
import CodeMirror, {EditorView} from "@uiw/react-codemirror";
import {githubLight} from "@uiw/codemirror-theme-github";
import "./CodePanel.css"
import {useTranslation} from "react-i18next";

const {Title} = Typography;
const CodePanel: React.FC = () => {
    const {t, i18n} = useTranslation();
    return (
        <Flex vertical className={"code-panel"} flex={1}>
            <Title level={4}>{t("apiPanelTitle")}</Title>
            <div className={"code-container-divider"}/>
            <CodeMirror
                style={{overflowY:"auto"}}
                extensions={[langs.go(), EditorView.theme({
                    "&.cm-focused": {
                        outline: "none",
                    },
                })]}
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