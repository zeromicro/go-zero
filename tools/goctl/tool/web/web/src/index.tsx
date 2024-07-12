import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './components/app/App';
import {
    createBrowserRouter,
    RouterProvider,
} from "react-router-dom";
import "./i18n";
import Welcome from "./components/welcome/Welcome";
import NotFound from "./components/notfound/NotFound";
import {ConfigProvider} from 'antd';

const router = createBrowserRouter([
    {
        path: "*",
        element: <NotFound/>
    },
    {
        path: "/welcome",
        element: <Welcome/>
    },
    {
        path: "/",
        element: <Welcome/>
    },
    {
        path: "/",
        element: <App/>,
        children: [
            {
                path: "api",
                children: [
                    {
                        path: "builder",
                        element: <Welcome/>
                    },
                ]
            }
        ]
    },
]);

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);
root.render(
    <ConfigProvider
        theme={{
            token: {
                colorPrimary: "#333333",
                colorInfo: "#333333",
            },
        }}
    >
        <React.StrictMode>
            <RouterProvider router={router}/>
        </React.StrictMode>
    </ConfigProvider>
);
