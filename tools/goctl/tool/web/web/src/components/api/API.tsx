import React, { useState } from "react";
import { Layout, Flex } from "antd";
import "../../Base.css";
import "./API.css";
import FormPanel from "./form/FormPanel";
import CodePanel from "./form/CodePanel";

const App: React.FC = () => {
  const [code, setCode] = useState("");
  return (
    <Layout className="api">
      <Flex vertical={false} wrap className="api-container" gap={1}>
        <div className={"api-form-panel"}>
          <FormPanel
            onBuild={(data) => {
              setCode(data);
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
