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
import API from "./components/api/API";

const router = createBrowserRouter([
    {
        path: "*",
        element: <NotFound/>
    },
    {
        path: "/",
        element: <Welcome/>
    },
    {
        path: "/",
        element: <App/>,
        children: [
            // {
            //     path: "home",
            //     element: <Home/>
            // },
            {
                path: "api",
                children: [
                    {
                        path: "builder",
                        element: <API/>
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
    <React.StrictMode>
        <RouterProvider router={router}/>
    </React.StrictMode>
);
