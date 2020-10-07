<template>
  <div class="user" v-loading="!inited">
    <!-- 用户列表 -->
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-user-solid" />
        <span>用户列表</span>
      </div>
      <!-- 搜索条件 -->
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <div v-loading="processing">
        <el-table
          :data="users"
          row-key="id"
          stripe
          @sort-change="handleSortChange"
        >
          <el-table-column prop="id" key="id" label="ID" width="80" sortable />
          <el-table-column
            prop="account"
            key="account"
            label="账户"
            width="120"
          />
          <el-table-column
            prop="statusDesc"
            key="statusDesc"
            label="状态"
            width="80"
          />
          <el-table-column label="用户角色">
            <template slot-scope="scope">
              <ul>
                <li v-for="role in scope.row.rolesDesc" :key="role">
                  {{ role }}
                </li>
              </ul>
            </template>
          </el-table-column>
          <el-table-column
            prop="updatedAtDesc"
            key="updatedAtDesc"
            label="更新于"
            width="180"
            sortable
          />
          <el-table-column label="操作" width="120">
            <template slot-scope="scope">
              <el-button
                class="op"
                type="text"
                size="small"
                @click="modify(scope.row)"
                >编辑</el-button
              >
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          class="pagination"
          layout="prev, pager, next, sizes"
          :current-page="currentPage"
          :page-size="query.limit"
          :page-sizes="pageSizes"
          :total="userCount"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    <User v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import User from "@/components/User.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { PAGE_SIZES } from "@/constants/common";

const userRoles = [];
const userStatuses = [];
const filterFields = [
  {
    label: "用户角色：",
    key: "role",
    type: "select",
    options: userRoles,
    span: 6
  },
  {
    label: "用户状态：",
    key: "status",
    type: "select",
    options: userStatuses,
    span: 6
  },
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入关键字",
    clearable: true,
    span: 6
  },
  {
    label: "",
    type: "filter",
    span: 6,
    labelWidth: "0px"
  }
];

export default {
  name: "Users",
  extends: BaseTable,
  components: {
    User,
    BaseFilter
  },
  data() {
    return {
      inited: false,
      filterFields: null,
      pageSizes: PAGE_SIZES,
      query: {
        offset: 0,
        limit: PAGE_SIZES[0],
        order: "-updatedAt"
      }
    };
  },
  computed: {
    ...mapState({
      userCount: state => state.user.list.count,
      users: state => state.user.list.data || [],
      userRoles: state => state.user.roles || [],
      userStatuses: state => state.user.statuses || [],
      processing: state => state.user.listProcessing,
      updateProcessing: state => state.user.updateProcessing
    })
  },
  methods: {
    ...mapActions(["listUser", "listUserRole", "listUserStatus"]),
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      try {
        await this.listUser(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      const { userRoles } = await this.listUserRole();
      const { statuses } = await this.listUserStatus();
      userRoles.length = 0;
      userRoles.push({
        name: "所有",
        value: ""
      });
      userRoles.push(...userRoles);

      userStatuses.length = 0;
      userStatuses.push({
        name: "所有",
        value: ""
      });
      userStatuses.push(...statuses);
      this.filterFields = filterFields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.inited = true;
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.user
  margin: $mainMargin
  i
    margin-right: 5px
  ul
    li
      list-style: inside
.selector, .submit
  width: 100%
.pagination
  text-align: right
  margin-top: 15px
</style>
