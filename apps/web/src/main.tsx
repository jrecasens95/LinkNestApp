import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./app";
import { AppProviders } from "./app/providers";
import { BrowserRouter } from "react-router-dom";
import "@radix-ui/themes/styles.css";
import { Theme } from "@radix-ui/themes";
import "./styles/index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Theme appearance="dark" accentColor="crimson" grayColor="sand" radius="large">
    <AppProviders>
    <BrowserRouter>
    <App />
    </BrowserRouter>
    </AppProviders>
    </Theme>
  </React.StrictMode>
);
