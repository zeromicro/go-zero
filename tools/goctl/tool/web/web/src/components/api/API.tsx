import React, { useState } from "react";
import { Layout, Row, Col, Flex } from "antd";
import "../../Base.css";
import "./API.css";
import FormPanel from "./form/FormPanel";
import CodePanel from "./form/CodePanel";

const App: React.FC = () => {
  const [code, setCode] = useState(`{
  "code": 0,
  "msg": "This is a json placeholder",
  "data": {
    "author": {
      "first": "an",
      "last": "keson"
    },
    "github": "https://github.com/kesonan/converter",
    "star": 1800,
    "nickname": "kesonan",
    "description": "A Golang Software Developer engineer from ShangHai.",
    "active": true,
    "tags": ["goctl","go-zero","go","java","android"],
    "projects": [
      {
        "organization": "zeromicro",
        "name": "go-zero",
        "github": "https://github.com/zeromicro/go-zero",
        "stars": 28200,
        "description": "A cloud-native Go microservices framework with cli tool for productivityA cloud-native Go microservices framework with cli tool for productivity."
      },
      {
        "organization": "kesonan",
        "name": "github-compare",
        "github": "https://github.com/kesonan/github-compare",
        "stars": 149,
        "description": "A GitHub repositories statistics command-line tool for the terminal"
      },
      {
        "organization": "kesonan",
        "name": "goimportx",
        "github": "https://github.com/kesonan/goimportx",
        "stars": 10,
        "description": "A tool to help you manage your go imports."
      }
    ],
    "others": [],
    "job": {}
  }
}`);
  return (
    <Layout className="api">
      <Flex vertical={false} wrap className="api-container" gap={1}>
        <div className={"api-form-panel"}>
          <FormPanel
            onBuild={(data) => {
              const js = JSON.stringify(data);
              setCode(js);
            }}
          />
        </div>
        <div className={"api-code-panel"}>
          <CodePanel
            onChange={(code) => {
              setCode(code);
            }}
            value={code}
          />
        </div>
      </Flex>
    </Layout>
  );
};

export default App;
