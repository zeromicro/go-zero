import React, { useState, useEffect } from "react";
import {
  Button,
  Layout,
  ConfigProvider,
  Flex,
  Typography,
  Tag,
  Dropdown,
  Space,
  MenuProps,
} from "antd";
import "../../Base.css";
import "./Welcome.css";
import { useTranslation } from "react-i18next";
import { ConverterIcon } from "../../util/icon";
import { useNavigate } from "react-router-dom";
import { DownOutlined } from "@ant-design/icons";
import zhCN from "antd/locale/zh_CN";
import enUS from "antd/locale/en_US";

const { Title } = Typography;
const { Header } = Layout;
const items: MenuProps["items"] = [
  {
    key: "en",
    label: "EN",
  },
  {
    key: "zh",
    label: "中",
  },
];
const App: React.FC = () => {
  const [locale, setLocale] = useState(zhCN);
  const [localeZH, setLocaleZh] = useState(true);
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();

  useEffect(() => {
    setLocaleZh(i18n.language === "zh");
  }, []);

  const onClick: MenuProps["onClick"] = ({ key }) => {
    if (key === "zh") {
      i18n.changeLanguage("zh");
      setLocale(zhCN);
      setLocaleZh(true);
    } else {
      i18n.changeLanguage("en");
      setLocale(enUS);
      setLocaleZh(false);
    }
  };
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#333333",
          colorInfo: "#333333",
        },
      }}
    >
      <Layout className="welcome">
        <Header style={{ background: "transparent" }}>
          <Dropdown menu={{ items, onClick }}>
            <a onClick={(e) => e.preventDefault()}>
              <Space className={"welcome-language"}>
                Language: {localeZH ? "中" : "EN"}
                <DownOutlined />
              </Space>
            </a>
          </Dropdown>
        </Header>
        <Flex
          gap={30}
          justify={"space-evenly"}
          style={{ height: "100%" }}
          align={"center"}
        >
          <Flex vertical style={{ marginLeft: "10%" }} gap={10}>
            <span className="welcome-text-gradient">欢迎来到 goctl 网页端</span>
            <span className="welcome-text">Welcome to goctl web tool</span>
            <Flex
              gap={0}
              style={{ fontSize: 20, color: "#525252" }}
              align={"center"}
            >
              <Tag color="#f50" className="welcome-tag">
                # go-zero{" "}
              </Tag>
              <Tag color="#2db7f5" className="welcome-tag">
                # goctl{" "}
              </Tag>
              <Tag color="#87d068" className="welcome-tag">
                # api{" "}
              </Tag>
              <Tag color="#108ee9" className="welcome-tag">
                # generator{" "}
              </Tag>
            </Flex>
            <Flex style={{ marginTop: 50, height: 50 }} gap={30}>
              <Button
                type={"primary"}
                style={{ height: "100%", flex: 1 }}
                onClick={() => {
                  navigate("/api/builder");
                }}
              >
                <ConverterIcon
                  type={"icon-terminal"}
                  className="welcome-start-icon"
                />
                {t("welcomeStart")}
              </Button>

              <Button
                style={{ height: "100%", flex: 1 }}
                onClick={() => {
                  window.open("https://go-zero.dev", "_blank");
                }}
              >
                <ConverterIcon
                  type={"icon-document"}
                  className="welcome-start-icon"
                />
                {t("welcomeDoc")}
              </Button>
            </Flex>
          </Flex>

          <Flex flex={1}>
            <Title></Title>
          </Flex>
        </Flex>
      </Layout>
    </ConfigProvider>
  );
};

export default App;
