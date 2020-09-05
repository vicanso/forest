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
const userGroups = [];
const userStatuses = [];
const userMarketingGroups = [];
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
  {
    label: "用户组：",
    key: "groups",
    type: "select",
    placeholder: "请选择用户分组",
    multiple: true,
    options: userGroups
    // rules: [
    //   {
    //     required: true,
    //     message: "用户分组不能为空"
    //   }
    // ]
  },
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
  },
  {
    label: "销售分组：",
    key: "marketingGroup",
    type: "select",
    placeholder: "请选择用户销售分组",
    // multiple: true,
    labelWidth: "100px",
    options: userMarketingGroups
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
      "listUserGroup",
      "listUserMarketingGroup",
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
      const { roles } = await this.listUserRole();
      const { groups } = await this.listUserGroup();
      const { statuses } = await this.listUserStatus();
      const { marketingGroups } = await this.listUserMarketingGroup();
      userRoles.length = 0;
      userRoles.push(...roles);
      userGroups.length = 0;
      userGroups.push(...groups);
      userStatuses.length = 0;
      userStatuses.push(...statuses);
      userMarketingGroups.length = 0;
      userMarketingGroups.push({
        name: "NULL",
        value: "NULL"
      });
      userMarketingGroups.push(...marketingGroups);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
