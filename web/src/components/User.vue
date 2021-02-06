<template lang="pug">
.user(
  v-loading="processing"
): base-editor(
  v-if="!processing && fields"
  title="更新用户信息"
  icon="el-icon-user"
  :id="id"
  :findByID="findByID"
  :updateByID="updateByID"
  :fields="fields"
)

</template>

<script lang="ts">
import { defineComponent } from "vue";

import useUserState, {
  userFindByID,
  userUpdateByID,
  userListRole,
} from "../states/user";
import useCommonState, { commonListStatus } from "../states/common";
import BaseEditor from "./base/Editor.vue";

const roleSelectList = [];
const statusSelectList = [];
const fields = [
  {
    label: "账号：",
    key: "account",
    disabled: true,
  },
  {
    label: "用户角色：",
    key: "roles",
    type: "select",
    placeholder: "请选择用户角色",
    labelWidth: "100px",
    multiple: true,
    options: roleSelectList,
    rules: [
      {
        required: true,
        message: "用户角色不能为空",
      },
    ],
  },
  // {
  //   label: "用户组：",
  //   key: "groups",
  //   type: "select",
  //   placeholder: "请选择用户分组",
  //   multiple: true,
  //   options: userGroups
  //   // rules: [
  //   //   {
  //   //     required: true,
  //   //     message: "用户分组不能为空"
  //   //   }
  //   // ]
  // },
  {
    label: "用户状态：",
    key: "status",
    type: "select",
    placeholder: "请选择用户状态",
    labelWidth: "100px",
    options: statusSelectList,
    rules: [
      {
        required: true,
        message: "用户状态不能为空",
      },
    ],
  },
];

export default defineComponent({
  name: "User",
  components: {
    BaseEditor,
  },
  setup() {
    const userState = useUserState();
    const commonState = useCommonState();
    return {
      findByID: userFindByID,
      updateByID: userUpdateByID,
      userRoles: userState.roles,
      statuses: commonState.statuses,
      getStatusDesc: (status) => {
        let desc = "";
        commonState.statuses.items.forEach((item) => {
          if (item.value === status) {
            desc = item.name;
          }
        });
        return desc;
      },
    };
  },
  data() {
    return {
      fields: null,
      processing: false,
      id: 0,
    };
  },
  async beforeMount() {
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      this.processing = true;
      await userListRole();
      await commonListStatus();

      // 重置
      roleSelectList.length = 0;
      roleSelectList.push(...this.userRoles.items);

      // 重置
      statusSelectList.length = 0;
      statusSelectList.push(...this.statuses.items);

      this.fields = fields;
    } catch (err) {
      this.$error(err);
    } finally {
      this.processing = false;
    }
  },
});
</script>
