import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import "./index.css";
import Main from "./route/Main";
import AuthPage from "./route/Auth";
import Profile from "./route/Profile";
import Complain from "./route/Complain";
import Statistics from "./route/Statistics";
import Status from "./route/Status";
import Contacts from "./route/Contacts";
import Admin from "./route/Admin";

import "primereact/resources/themes/saga-orange/theme.css";

import "primereact/resources/primereact.min.css";
import "primeicons/primeicons.css";



const router = createBrowserRouter(
  [
    {
      path: "/Main",
      element: <Main />,
    },
    {
      path: "/Auth",
      element: <AuthPage />,
    },
    {
      path: "/Profile",
      element: <Profile />,
    },
    {
      path: "/Complain",
      element: <Complain />,
    },
    {
      path: "/Stat",
      element: <Statistics />,
    },
    {
      path: "/Status",
      element: <Status />,
    },
    {
      path: "/Contacts",
      element: <Contacts />,
    },
    {
      path: "/Admin",
      element: <Admin />,
    }
  ],
  {
    basename: process.env.PUBLIC_URL
  }
);

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    <RouterProvider router={router}>
    </RouterProvider>
  </React.StrictMode>
);
