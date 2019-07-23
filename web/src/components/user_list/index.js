import React from "react";
import {
  Row,
  Col,
  Input,
  Card,
  Select,
  Table,
  message,
  Spin,
  Button,
  Form
} from "antd";
import moment from "moment";
import axios from "axios";

import "./user_list.sass";
import { TIME_FORMAT } from "../../vars";
import { USERS } from "../../urls";

const { Search } = Input;
const { Option } = Select;
const editMode = "edit";

const roles = ["su", "admin"];

class UserList extends React.Component {
  state = {
    mode: "",
    keyword: "",
    current: null,
    loading: false,
    users: null
  };
  async search() {
    const { loading, keyword } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      const { data } = await axios.get(USERS, {
        params: {
          keyword
        }
      });
      this.setState({
        users: data.users
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  handleSubmit(e) {
    e.preventDefault();
  }
  renderTable() {
    const { users } = this.state;
    const columns = [
      {
        title: "用户名",
        dataIndex: "account",
        key: "account"
      },
      {
        title: "角色",
        dataIndex: "roles",
        key: "roles"
      },
      {
        title: "创建于",
        dataIndex: "createdAt",
        key: "createdAt",
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      },
      {
        title: "操作",
        key: "op",
        width: "100px",
        render: (text, record) => {
          return (
            <a
              href="/update"
              onClick={e => {
                e.preventDefault();
                this.setState({
                  current: record,
                  mode: editMode
                });
              }}
            >
              更新
            </a>
          );
        }
      }
    ];
    return <Table className="users" dataSource={users} columns={columns} />;
  }
  renderRoles() {
    return roles.map(item => (
      <Option key={item} value={item}>
        {item}
      </Option>
    ));
  }
  renderUserList() {
    const { loading, mode } = this.state;
    if (mode === editMode) {
      return;
    }
    return (
      <div>
        <Card title="用户搜索" size="small">
          <Spin spinning={loading}>
            <div className="filter">
              <Select className="roles" placeholder="请选择用户角色">
                {this.renderRoles()}
              </Select>
              <Search
                className="keyword"
                placeholder="请输入关键字"
                onKeyDown={e => {
                  if (e.keyCode === 0x0d) {
                    this.search();
                  }
                }}
                onSearch={keyword => {
                  this.setState({
                    keyword
                  });
                  this.search();
                }}
                enterButton
              />
            </div>
          </Spin>
        </Card>
        {this.renderTable()}
      </div>
    );
  }
  renderEditor() {
    const { mode, current } = this.state;
    if (mode !== editMode) {
      return;
    }
    const colSpan = 12;
    return (
      <Card title="更新用户信息" size="small">
        <Form onSubmit={this.handleSubmit.bind(this)}>
          <Row gutter={24}>
            <Col span={colSpan}>
              <Form.Item label="用户名">
                <Input disabled defaultValue={current.account} />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="用户角色">
                <Select mode="multiple" placeholder="请选择要添加的角色">
                  {this.renderRoles()}
                </Select>
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Button className="submit" type="primary" htmlType="submit">
                更新
              </Button>
            </Col>
            <Col span={colSpan}>
              <Button
                className="back"
                onClick={() => {
                  this.setState({
                    mode: ""
                  });
                }}
              >
                返回
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    );
  }
  render() {
    return (
      <div className="UserList">
        {this.renderUserList()}
        {this.renderEditor()}
      </div>
    );
  }
}

export default UserList;
