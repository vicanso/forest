import React from "react";
import {
  Button,
  Card,
  Typography,
  Select,
  Form,
  Col,
  notification
} from "antd";
import PropTypes from "prop-types";

import ConfigEditor from "../config_editor";
import ConfigTable from "../config_table";
import "./signed_keys_config.sass";

const { Paragraph } = Typography;
const signedKeyCategory = "signedKey";
const editMode = "edit";

class SignedKeysConfig extends React.Component {
  state = {
    mode: "",
    currentData: null,
    currentKeys: null
  };
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
    const { mode, currentData, currentKeys } = this.state;
    if (mode !== editMode) {
      return;
    }
    const originalData = currentData || {
      category: signedKeyCategory
    };
    const content = (
      <Col span={8}>
        <Form.Item label="Key列表">
          <Select
            defaultValue={currentKeys || []}
            mode="tags"
            placeholder="请输入需要配置的key"
            onChange={value => {
              this.setState({
                currentKeys: value
              });
            }}
          ></Select>
        </Form.Item>
      </Col>
    );
    return (
      <Card size="small" title="添加/更新签名配置">
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
          category: signedKeyCategory
        }}
        formatData={data => {
          if (!data) {
            return "[]";
          }
          return JSON.stringify(data.split(","));
        }}
        onUpdate={data => {
          this.setState({
            currentKeys: data.data.split(","),
            mode: editMode,
            currentData: data
          });
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
    return <div className="SignedKeysConfig">{this.renderContent()}</div>;
  }
}

SignedKeysConfig.propTypes = {
  account: PropTypes.string.isRequired,
  isAdmin: PropTypes.bool.isRequired
};

export default SignedKeysConfig;
