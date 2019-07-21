import React from "react";
import PropTypes from "prop-types";
import {
  Button,
  Card,
  Typography,
  Select,
  Form,
  Col,
  notification,
  Input,
  message
} from "antd";
import axios from "axios";

import ConfigEditor from "../config_editor";
import ConfigTable from "../config_table";
import "./router_config.sass";
import { ROUTERS } from "../../urls";

const { Paragraph } = Typography;
const Option = Select.Option;
const { TextArea } = Input;
const editMode = "edit";

const routerConfigCategory = "router-config";

class RouterConfig extends React.Component {
  state = {
    routers: null,
    mode: "",
    currentData: null,
    currentKeys: null
  };
  async componentWillMount() {
    try {
      const { data } = await axios.get(ROUTERS);
      this.setState({
        routers: data.routers
      });
    } catch (err) {
      message.error(err);
    }
  }
  componentDidMount() {
    const { isAdmin } = this.props;
    if (!isAdmin) {
      notification.open({
        message: "请使用管理员登录",
        description: "此功能需要先登录并有管理权限才可使用"
      });
    }
  }
  reset() {
    this.setState({
      mode: "",
      currentKeys: null,
      currentData: null
    });
  }
  renderConfigEditor() {
    const { mode, currentData, routers } = this.state;
    if (mode !== editMode) {
      return;
    }
    const originalData = currentData || {
      category: routerConfigCategory
    };
    const opts = routers.map(item => {
      const { method, path } = item;
      const key = `${method} ${path}`;
      return (
        <Option key={key} value={key}>
          {key}
        </Option>
      );
    });
    const contentTypes = [
      "application/json; charset=UTF-8",
      "text/plain; charset=UTF-8",
      "text/html; charset=UTF-8"
    ].map(item => {
      return (
        <Option key={item} value={item}>
          {item}
        </Option>
      );
    });
    const colSpan = 8;
    const content = (
      <div>
        <Col span={colSpan}>
          <Form.Item label="路由选择">
            <Select placeholder="请选择要配置的路由">{opts}</Select>
          </Form.Item>
        </Col>
        <Col span={colSpan}>
          <Form.Item label="状态码">
            <Input type="number" placeholder="请输入响应状态码" />
          </Form.Item>
        </Col>
        <Col span={colSpan}>
          <Form.Item label="响应类型">
            <Select placeholder="请选择响应数据类型">{contentTypes}</Select>
          </Form.Item>
        </Col>
        <Col>
          <Form.Item label="响应数据">
            <TextArea
              autosize={{
                minRows: 4,
              }}
            />
          </Form.Item>
        </Col>
      </div>
    );
    return (
      <Card size="small" title="添加/更新路由配置">
        <Paragraph>用于生成session的cookie认证</Paragraph>
        <ConfigEditor
          originalData={originalData}
          content={content}
          getConfigData={() => {
            const { currentKeys } = this.state;
            if (!currentKeys) {
              return "";
            }
            return currentKeys.join(",");
          }}
          onSuccess={this.reset.bind(this)}
        />
        <Button className="back" type="primary" onClick={this.reset.bind(this)}>
          返回
        </Button>
      </Card>
    );
  }
  renderTable() {
    const { mode } = this.state;
    if (mode) {
      return;
    }
    return (
      <ConfigTable
        params={{
          category: routerConfigCategory
        }}
        formatData={data => {
          if (!data) {
            return "";
          }
          return JSON.stringify(data);
        }}
        onUpdate={data => {
          console.dir(data);
          // this.setState({
          //   currentKeys: data.data.split(","),
          //   mode: editMode,
          //   currentData: data
          // });
        }}
      />
    );
  }
  renderContent() {
    const { mode } = this.state;
    const { isAdmin } = this.props;
    if (!isAdmin) {
      return;
    }
    return (
      <div>
        {this.renderTable()}
        {this.renderConfigEditor()}
        {!mode && (
          <Button
            onClick={() => {
              this.setState({
                mode: editMode
              });
            }}
            type="primary"
            className="add"
          >
            添加
          </Button>
        )}
      </div>
    );
  }
  render() {
    return <div className="RouterConfig">{this.renderContent()}</div>;
  }
}

RouterConfig.propTypes = {
  isAdmin: PropTypes.bool.isRequired
};

export default RouterConfig;
