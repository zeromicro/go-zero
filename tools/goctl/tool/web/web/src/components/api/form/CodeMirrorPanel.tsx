import React, { useEffect, useState } from "react";
import { Modal, Tooltip } from "antd";
import { FullscreenOutlined } from "@ant-design/icons";
import { useTranslation } from "react-i18next";
import CodeMirror, { EditorView } from "@uiw/react-codemirror";
import { githubLight } from "@uiw/codemirror-theme-github";
import { langs } from "@uiw/codemirror-extensions-langs";
import type { ViewUpdate } from "@codemirror/view";

interface CodeMirrorPanelProps {
  value: string;
  placeholder?: string;

  onChange?(value: string, viewUpdate?: ViewUpdate): void;
}

const CodeMirrorPanel: React.FC<
  CodeMirrorPanelProps & React.RefAttributes<HTMLDivElement>
> = (props) => {
  const { t } = useTranslation();
  const [modalOpen, setModalOpen] = useState(false);
  const [code, setCode] = useState("");
  const [modalCode, setModalCode] = useState("");

  useEffect(() => {
    setCode(props.value);
  });

  return (
    <div>
      <Modal
        title={t("formResponseBodyModelTitle")}
        centered
        open={modalOpen}
        maskClosable={false}
        keyboard={false}
        closable={false}
        destroyOnClose
        onOk={() => {
          setModalOpen(false);
          if (props.onChange) {
            props.onChange(modalCode);
          }
        }}
        onCancel={() => {
          setModalOpen(false);
        }}
        width={1000}
        cancelText={t("formResponseBodyModalCancel")}
        okText={t("formResponseBodyModalConfirm")}
      >
        <CodeMirror
          style={{ marginTop: 10, overflow: "auto" }}
          extensions={[
            langs.json(),
            EditorView.theme({
              "&.cm-focused": {
                outline: "none",
              },
            }),
          ]}
          theme={githubLight}
          height={"70vh"}
          value={code}
          onChange={(code) => {
            setModalCode(code);
          }}
        />
      </Modal>

      <span
        style={{
          position: "absolute",
          top: -30,
          right: 0,
          zIndex: 1000,
        }}
      >
        <Tooltip title={t("tooltipFullScreen")}>
          <FullscreenOutlined
            style={{ cursor: "pointer" }}
            onClick={() => {
              setModalCode(code);
              setModalOpen(true);
            }}
          />
        </Tooltip>
      </span>

      <CodeMirror
        style={{ overflow: "scroll", minHeight: 100, maxHeight: 200 }}
        extensions={[
          langs.json(),
          EditorView.theme({
            "&.cm-focused": {
              outline: "none",
            },
          }),
        ]}
        value={code}
        placeholder={props.placeholder ? props.placeholder : ""}
        theme={githubLight}
        onChange={(value, viewUpdate) => {
          if (props.onChange) {
            props.onChange(value, viewUpdate);
          }
        }}
      />
    </div>
  );
};

export default CodeMirrorPanel;
