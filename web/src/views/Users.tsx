import { defineComponent, onUnmounted } from "vue";
import ExTable from "../components/ExTable";
import useUserState, { userList, userListClear } from "../states/user";

function getColumns() {
  return [
    {
      title: "账户",
      key: "account",
    },
    {
      title: "角色",
      key: "roles",
    },
    {
      title: "状态",
      key: "statusDesc",
    },
    {
      title: "邮箱",
      key: "email",
    },
  ];
}

function getFilters() {
  return [
    {
      name: "角色：",
      key: "role",
      placeholder: "请选择要筛选的用户角色",
      type: "select",
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "超级用户",
          value: "su",
        },
        {
          label: "管理员",
          value: "admin",
        },
      ],
    },
    {
      name: "状态：",
      key: "status",
      placeholder: "请选择要筛选的账户状态",
      type: "select",
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "启用",
          value: "1",
        },
        {
          label: "禁用",
          value: "2",
        },
      ],
    },
    {
      name: "关键字：",
      key: "keyword",
      placeholder: "请输入搜索关键字",
    },
  ];
}
export default defineComponent({
  name: "Users",
  setup() {
    const { users } = useUserState();

    const fetchUsers = async (params: {
      limit: number;
      offset: number;
      keyword?: string;
      status?: string;
      role?: string;
    }) =>
      userList({
        limit: params.limit,
        offset: params.offset,
        keyword: params.keyword || "",
        status: params.status || "",
        role: params.role || "",
      });

    onUnmounted(() => {
      userListClear();
    });

    return {
      users,
      fetchUsers,
    };
  },
  render() {
    const { users, fetchUsers } = this;
    return (
      <ExTable
        title={"用户查询"}
        columns={getColumns()}
        filters={getFilters()}
        data={users}
        fetch={fetchUsers}
      />
    );
  },
});
