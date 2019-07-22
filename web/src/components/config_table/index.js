import React from "react";
import PropTypes from "prop-types";
import { message, Table, Icon, Divider, Spin } from "antd";
import axios from "axios";
import moment from "moment";

import {
  CONFIGURATIONS_LIST,
  CONFIGURATIONS_DELETE,
  CONFIGURATIONS_LIST_AVAILABLE
} from "../../urls";
import { TIME_FORMAT } from "../../vars";
import "./config_table.sass";

class ConfigTable extends React.Component {
  state = {
    configs: null,
    loading: false
  };
  async componentWillMount() {
    const { params } = this.props;
    this.setState({
      loading: true
    });
    try {
      let url = CONFIGURATIONS_LIST;
      if (params.available) {
        url = CONFIGURATIONS_LIST_AVAILABLE;
        delete params.available;
      }
      const { data } = await axios.get(url, {
        params: params
      });
      const { configs } = data;
      configs.forEach(item => {
        item.key = `${item.id}`;
      });
      this.setState({
        configs
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  async deleteConfig(id) {
    const { configs } = this.state;
    try {
      const url = CONFIGURATIONS_DELETE.replace(":id", id);
      await axios.delete(url);
      const result = [];
      configs.forEach(item => {
        if (item.id !== id) {
          result.push(item);
        }
      });
      this.setState({
        configs: result
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  render() {
    const { configs, loading } = this.state;
    const { onUpdate, formatData } = this.props;
    const columns = [
      {
        title: "名称",
        dataIndex: "name",
        key: "name"
      },
      {
        title: "分类",
        dataIndex: "category",
        key: "category"
      },
      {
        title: "是否启用",
        dataIndex: "enabled",
        key: "enabled",
        width: "100px",
        render: value => {
          if (value) {
            return <Icon type="check-circle" theme="twoTone" />;
          }
          return <Icon type="close-circle" />;
        }
      },
      {
        title: "配置",
        dataIndex: "data",
        key: "data",
        render: value => {
          if (formatData) {
            return formatData(value);
          }
          return value;
        }
      },
      {
        title: "开始时间",
        dataIndex: "beginDate",
        key: "beginDate",
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      },
      {
        title: "结束时间",
        dataIndex: "endDate",
        key: "endDate",
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
        width: "120px",
        render: (text, record) => {
          return (
            <span>
              {onUpdate && (
                <a
                  href="/update"
                  onClick={e => {
                    e.preventDefault();
                    onUpdate(record);
                  }}
                >
                  更新
                </a>
              )}
              {onUpdate && <Divider type="vertical" />}
              <a
                href="/delete"
                onClick={e => {
                  e.preventDefault();
                  this.deleteConfig(record.id);
                }}
              >
                删除
              </a>
            </span>
          );
        }
      }
    ];
    return (
      <div className="ConfigTable">
        {loading && (
          <div className="loadingWrapper">
            <Spin tip="加载中..." />
          </div>
        )}
        {!loading && <Table dataSource={configs} columns={columns} />}
      </div>
    );
  }
}

ConfigTable.propTypes = {
  params: PropTypes.object,
  onUpdate: PropTypes.func,
  formatData: PropTypes.func
};

export default ConfigTable;
