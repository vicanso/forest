import React from "react";
import { Route, HashRouter } from "react-router-dom";
import axios from "axios";
import { message, Spin } from "antd";

import "./app.sass";
import {
  ALL_CONFIG_PATH,
  BASIC_CONFIG_PATH,
  SIGNED_KEYS_CONFIG_PATH,
  ROUTER_CONFIG_PATH,
  REGISTER_PATH,
  LOGIN_PATH
} from "./paths";
import { USERS_ME } from "./urls";
import AppMenu from "./components/app_menu";
import AppHeader from "./components/app_header";
import BasicConfig from "./components/basic_config";
import SignedKeysConfig from "./components/signed_keys_config";
import Login from "./components/login";
import Register from "./components/register";
import RouterConfig from "./components/router_config";
import ConfigList from "./components/config_list";

class App extends React.Component {
  state = {
    loading: false,
    account: "",
    isAdmin: false,
  };
  async componentWillMount() {
    this.setState({
      loading: true
    });
    try {
      const { data } = await axios.get(USERS_ME);
      this.setUserInfo(data);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
    // 更新session与cookie有效期
    setTimeout(() => {
      axios.patch(USERS_ME)
    }, 5 * 1000);
  }
  setUserInfo(data)  {
    let isAdmin = false;
    (data.roles || []).forEach((item) => {
      if (item === "su" || item === "admin") {
        isAdmin = true;
      } 
    });
    this.setState({
      account: data.account || "",
      isAdmin,
    });
  }
  render() {
    const {
      account,
      isAdmin,
      loading,
    } = this.state;
    return (
      <div className="App">
        <HashRouter>
          <AppMenu />
          {loading && <div className="loadingWrapper">
            <Spin tip="加载中..." />
          </div>}
          {!loading && 
            <div className="contentWrapper">
              <AppHeader
                account={account}
                setUserInfo={this.setUserInfo.bind(this)}
              />

              <Route
                path={LOGIN_PATH}
                render={(props) => <Login
                  {...props}
                  setUserInfo={this.setUserInfo.bind(this)}
                />}
              />
              <Route path={REGISTER_PATH} component={Register} />
              <Route
                path={ALL_CONFIG_PATH}
                render={(props) => <ConfigList
                  {...props}
                  account={account}
                  isAdmin={isAdmin}
                />}
              />
              <Route
                path={BASIC_CONFIG_PATH}
                render={(props) => <BasicConfig
                  {...props}
                  account={account}
                  isAdmin={isAdmin}
                />}
              />
              <Route
                path={SIGNED_KEYS_CONFIG_PATH}
                render={(props) => <SignedKeysConfig
                  {...props}
                  account={account}
                  isAdmin={isAdmin}
                />}
              />
              <Route
                path={ROUTER_CONFIG_PATH}
                render={(props) => <RouterConfig
                  {...props}
                  account={account}
                  isAdmin={isAdmin}
                />}
              />
            </div>
          } 
        </HashRouter>
      </div>
    );
  }
}

export default App;
