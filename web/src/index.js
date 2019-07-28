import React from "react";
import ReactDOM from "react-dom";
import { Router } from "react-router";
import "antd/dist/antd.css";
import "./index.sass";
import * as serviceWorker from "./serviceWorker";

import App from "./App";
import "./request-interceptors";
// import * as serviceWorker from "./serviceWorker";

import { getHistory } from "./router";

// import * as serviceWorker from './serviceWorker';

ReactDOM.render(
  <Router history={getHistory()}>
    <App />
  </Router>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
