import React from "react";

import { Menu, Icon } from "antd";
import { Link } from "react-router-dom";

import {
  BASIC_CONFIG_PATH,
  ROUTER_CONFIG_PATH,
  SIGNED_KEYS_CONFIG_PATH
} from "../../paths";
import "./app_menu.sass";

const { SubMenu } = Menu;

class AppMenu extends React.Component {
  render() {
    return (
      <div className="AppMenu">
        <Menu mode="inline" theme="dark" defaultOpenKeys={["configuration"]}>
          <SubMenu
            key="configuration"
            title={
              <span>
                <Icon type="setting" />
                <span>配置</span>
              </span>
            }
          >
            <Menu.Item key="basic-config">
              <Link to={BASIC_CONFIG_PATH}>基本配置</Link>
            </Menu.Item>
            <Menu.Item key="router-config">
              <Link to={ROUTER_CONFIG_PATH}>路由配置</Link>
            </Menu.Item>
            <Menu.Item>
              <Link to={SIGNED_KEYS_CONFIG_PATH}>签名配置</Link>
            </Menu.Item>
          </SubMenu>
        </Menu>
      </div>
    );
  }
}

export default AppMenu;
