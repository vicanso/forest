import React from "react";
import axios from "axios";
import { Form, Input, Icon, Card, Button, message } from "antd";

import { USERS_LOGIN, USERS_ME } from "../../urls";
import { sha256 } from "../../helpers/crypto";
import "./login_register.sass";
import * as router from "../../router";

class LoginRegister extends React.Component {
  loginMode = "login";
  registerMode = "register";
  state = {
    submitting: false,
    account: "",
    password: "",
    token: "",
    mode: ""
  };
  async handleSubmit(e) {
    const { setUserInfo } = this.props;
    e.preventDefault();
    const { account, password, mode, token } = this.state;

    if (!account || !password) {
      message.error("用户名与密码不能为空");
      return;
    }
    let url = "";
    const postData = {
      account
    };
    if (mode === this.loginMode) {
      if (!token) {
        message.error("Token为不能空");
        return;
      }
      url = USERS_LOGIN;
      postData.password = sha256(sha256(password) + token);
    } else {
      url = USERS_ME;
      postData.password = sha256(password);
    }
    this.setState({
      submitting: true
    });
    try {
      const { data } = await axios.post(url, postData);
      if (setUserInfo) {
        setUserInfo({
          account: data.account || "",
          roles: data.roles
        });
      }
      router.back();
      // TODO 触发reload
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        submitting: false
      });
    }
  }
  render() {
    const { mode } = this.state;
    const title = mode === this.loginMode ? "登录" : "注册";
    return (
      <div className="LoginRegister">
        <Card title={title}>
          <Form onSubmit={this.handleSubmit.bind(this)}>
            <Form.Item>
              <Input
                autoFocus
                prefix={
                  <Icon type="user" style={{ color: "rgba(0,0,0,.25)" }} />
                }
                onChange={e => {
                  this.setState({
                    account: e.target.value.trim()
                  });
                }}
                placeholder="用户名"
              />
            </Form.Item>
            <Form.Item>
              <Input
                prefix={
                  <Icon type="lock" style={{ color: "rgba(0,0,0,.25)" }} />
                }
                type="password"
                onChange={e => {
                  this.setState({
                    password: e.target.value.trim()
                  });
                }}
                placeholder="密码"
              />
            </Form.Item>
            <Button type="primary" htmlType="submit" className="submit">
              {title}
            </Button>
          </Form>
        </Card>
      </div>
    );
  }
}

export default LoginRegister;
