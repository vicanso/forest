<template>
  <div class="user">
    <BaseEditor
      v-if="!processing && fields"
      title="更新用户信息"
      icon="el-icon-user"
      :id="id"
      :findByID="getUserByID"
      :updateByID="updateUserByID"
      :fields="fields"
    />
  </div>
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";
const userRoles = [];
const userStatuses = [];
const fields = [
  {
    label: "账号：",
    key: "account",
    disabled: true
  },
  {
    label: "用户角色：",
    key: "roles",
    type: "select",
    placeholder: "请选择用户角色",
    labelWidth: "100px",
    multiple: true,
    options: userRoles,
    rules: [
      {
        required: true,
        message: "用户角色不能为空"
      }
    ]
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
    options: userStatuses,
    rules: [
      {
        required: true,
        message: "用户状态不能为空"
      }
    ]
  }
];

export default {
  name: "User",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      processing: false,
      id: 0
    };
  },
  methods: {
    ...mapActions([
      "getUserByID",
      "listUserRole",
      "listUserStatus",
      "updateUserByID"
    ])
  },
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      const { userRoles } = await this.listUserRole();
      const { statuses } = await this.listUserStatus();
      userRoles.length = 0;
      userRoles.push(...userRoles);
      userStatuses.length = 0;
      userStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
