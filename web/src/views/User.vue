<template>
  <div class="user">
    <el-card v-if="mode === modifyMode">
      <div slot="header">
        <i class="el-icon-user" />
        <span>更新用户信息</span>
      </div>
      <el-form
        :model="currentUser"
        label-width="120px"
        v-loading="updateProcessing"
      >
        <el-row :gutter="15">
          <el-col :span="8">
            <el-form-item label="账号：">
              <el-input v-model="currentUser.account" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="用户角色：">
              <el-select
                class="selector"
                v-model="currentUser.roles"
                placeholder="请选择用户角色"
                multiple
              >
                <el-option
                  v-for="item in userRoles"
                  :key="item.value"
                  :label="item.name"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="用户状态：">
              <el-select
                class="selector"
                v-model="currentUser.status"
                placeholder="请选择用户状态"
              >
                <el-option
                  v-for="item in userStatuses"
                  :key="item.value"
                  :label="item.name"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item>
              <el-button class="submit" type="primary" @click="update"
                >更新</el-button
              >
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item>
              <el-button class="submit" @click="goBack">返回</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
    <el-card v-else>
      <div slot="header">
        <i class="el-icon-user-solid" />
        <span>用户列表</span>
      </div>
      <div class="pagination" v-loading="processing">
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
                <li v-for="role in scope.row.roleDescList" :key="role">
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
          :total="userCount"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import { diff } from "@/helpers/util";

const modifyMode = "modify";

export default {
  name: "User",
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      mode: "",
      modifyMode,
      currentUser: null,
      count: 0,
      pageSizes,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-updatedAt"
      }
    };
  },
  computed: {
    currentPage() {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    },
    ...mapState({
      userCount: state => state.user.list.count,
      users: state => state.user.list.data || [],
      userRoles: state => state.user.roles || [],
      userStatuses: state => state.user.statuses || [],
      processing: state => state.user.userListProcessing,
      updateProcessing: state => state.user.updateProcessing
    })
  },
  methods: {
    ...mapActions([
      "listUser",
      "listUserRole",
      "listUserStatus",
      "updateUserByID"
    ]),
    async fetch() {
      const { query } = this;
      try {
        await this.listUser(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    handleCurrentChange(page) {
      this.query.offset = (page - 1) * this.query.limit;
      this.fetch();
    },
    handleSizeChange(pageSize) {
      this.query.limit = pageSize;
      this.query.offset = 0;
      this.fetch();
    },
    handleSortChange({ prop, order }) {
      let key = prop.replace("Desc", "");
      if (order === "descending") {
        key = `-${key}`;
      }
      this.query.order = key;
      this.query.offset = 0;
      this.fetch();
    },
    async modify(data) {
      // 使用简单的方式，不修改路由参数
      this.mode = modifyMode;
      this.currentUser = Object.assign({}, data);
      try {
        await this.listUserRole();
        await this.listUserStatus();
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    async update() {
      const { users, currentUser } = this;
      let updateInfo = null;
      if (currentUser.roles.length === 0) {
        this.$message.warning("用户角色不能为空");
        return;
      }
      users.forEach(item => {
        if (item.id === currentUser.id) {
          updateInfo = diff(currentUser, item);
        }
      });
      if (!updateInfo || updateInfo.modifiedCount === 0) {
        this.$message.warning("请选择要更新的信息");
        return;
      }
      try {
        await this.updateUserByID({
          id: currentUser.id,
          data: updateInfo.data
        });
        this.goBack();
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    goBack() {
      this.mode = "";
    }
  },
  beforeMount() {
    this.fetch();
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
