import {UploadOutlined, UserOutlined, SmileFilled, VideoCameraOutlined} from "@ant-design/icons";
import React, {CSSProperties} from "react";
import type {MenuProps, MenuTheme} from 'antd';
import {ConverterIcon} from "../util/icon";

type MenuItem = Required<MenuProps>['items'][number];

const iconStyle: CSSProperties = {
    color: "#575757",
    fontSize: 24,
    textAlign: "center",
    marginLeft: "-4px"
}

export const menuItems: MenuItem[] = [
    {
        key: '/welcome',
        icon: <ConverterIcon style={iconStyle} type={"icon-welcome"}/>,
        label: '欢迎来到 goctl web'
    },
    {
        key: 'api',
        icon: <ConverterIcon style={iconStyle} type={"icon-code"}/>,
        label: 'API',
        children: [
            {
                key: '/api/builder',
                icon: <UserOutlined/>,
                label: '构造器',
            },
            {
                key: '2',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '3',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
            {
                key: '4',
                icon: <UserOutlined/>,
                label: 'nav 1',
            },
            {
                key: '5',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '6',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
        ]
    },
    {
        key: '2',
        icon: <VideoCameraOutlined/>,
        label: 'nav 2',
        children: [
            {
                key: '1',
                icon: <UserOutlined/>,
                label: 'nav 1',
            },
            {
                key: '2',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '3',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
            {
                key: '4',
                icon: <UserOutlined/>,
                label: 'nav 1',
            },
            {
                key: '5',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '6',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
        ]
    },
    {
        key: '3',
        icon: <UploadOutlined/>,
        label: 'nav 3',
        children: [
            {
                key: '1',
                icon: <UserOutlined/>,
                label: 'nav 1',
            },
            {
                key: '2',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '3',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
            {
                key: '4',
                icon: <UserOutlined/>,
                label: 'nav 1',
            },
            {
                key: '5',
                icon: <VideoCameraOutlined/>,
                label: 'nav 2',
            },
            {
                key: '6',
                icon: <UploadOutlined/>,
                label: 'nav 3',
            },
        ]
    },
]