import React, {useState} from 'react';
import '../../Base.css'
import {useTranslation} from "react-i18next";
import {Button, Result} from 'antd';
import {useNavigate} from "react-router-dom";

const App: React.FC = () => {
    const {t, i18n} = useTranslation();
    const navigate = useNavigate()
    return (
        <Result
            status="404"
            title="404"
            subTitle={t("notFoundDescription")}
            extra={<Button type="primary" onClick={()=>{navigate("/")}}>{t("backHome")}</Button>}
        />
    )
};

export default App;