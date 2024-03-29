import { defineComponent, onUnmounted } from "vue";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";

import ExTable from "../../components/ExTable";
import { today } from "../../helpers/util";
import { FormItemTypes } from "../../components/ExForm";
import { useUserLoginsStore } from "../../stores/user-logins";
import { storeToRefs } from "pinia";

function getColumns(): TableColumn[] {
  return [
    {
      title: "账户",
      key: "account",
      width: 120,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "IP",
      key: "ip",
      width: 120,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "定位",
      key: "location",
      width: 150,
    },
    {
      title: "运营商",
      key: "isp",
      width: 80,
    },
    {
      title: "Track ID",
      key: "trackID",
    },
    {
      title: "Session ID",
      key: "sessionID",
    },
    {
      title: "Forwarded For",
      key: "xForwardedFor",
      width: 140,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "User Agent",
      key: "userAgent",
      width: 120,
      ellipsis: {
        tooltip: true,
      },
    },
  ];
}

function getFilters() {
  return [
    {
      key: "account",
      name: "账户：",
      placeholder: "请输入要筛选的账号",
    },
    {
      key: "begin:end",
      name: "开始结束时间：",
      type: FormItemTypes.DateRange,
      placeholder: "请选择开始时间:请选择结束时间",
      span: 12,
      defaultValue: [today().toISOString(), new Date().toISOString()],
    },
  ];
}

export default defineComponent({
  name: "LoginStats",
  setup() {
    const userLoginsStore = useUserLoginsStore();
    const { processing, logins, count } = storeToRefs(userLoginsStore);

    onUnmounted(() => {
      userLoginsStore.$reset();
    });

    return {
      processing,
      logins,
      count,
      fetchLogins: userLoginsStore.list,
    };
  },
  render() {
    const { processing, logins, count, fetchLogins } = this;
    return (
      <ExTable
        title={"登录查询"}
        filters={getFilters()}
        columns={getColumns()}
        data={{
          count,
          items: logins,
          processing,
        }}
        fetch={fetchLogins}
      />
    );
  },
});
