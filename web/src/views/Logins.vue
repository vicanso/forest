<template>
  <el-card class="logins">
    <div slot="header">
      <i class="el-icon-user-solid" />
      <span>用户登录查询</span>
    </div>
    <div v-loading="processing">
      <el-form label-width="80px">
        <el-row :gutter="15">
          <el-col :span="6">
            <el-form-item label="账号：">
              <el-input
                v-model="query.account"
                placeholder="请输入要查询的账号"
                @keyup.enter.native="handleSearch"
              />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="时间：">
              <el-date-picker
                class="dateRange"
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-button
              :loading="processing"
              icon="el-icon-search"
              class="submit"
              type="primary"
              @click="handleSearch"
              >查询</el-button
            >
          </el-col>
        </el-row>
      </el-form>
      <el-table :data="logins" row-key="id" stripe>
        <el-table-column
          prop="account"
          key="account"
          label="账户"
          width="120"
        />

        <el-table-column prop="ip" key="ip" label="IP" width="120" />
        <el-table-column
          prop="location"
          key="location"
          label="定位"
          width="180"
        />
        <el-table-column label="运营商" width="80">
          <template slot-scope="scope">
            {{ scope.row.isp || "--" }}
          </template>
        </el-table-column>
        <el-table-column prop="sessionId" key="sessionId" label="Session ID" />
        <el-table-column prop="trackId" key="trackId" label="Track ID" />
        <el-table-column label="Forwarded For" width="120">
          <template slot-scope="scope">
            <el-tooltip>
              <span slot="content">{{ scope.row.xForwardedFor || "" }}</span>
              <i class="el-icon-info" />
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="User Agent" width="120">
          <template slot-scope="scope">
            <el-tooltip>
              <span slot="content">{{ scope.row.userAgent || "" }}</span>
              <i class="el-icon-mobile-phone" />
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column
          prop="createdAtDesc"
          key="createdAtDesc"
          label="时间"
          width="180"
        />
      </el-table>
      <el-pagination
        class="pagination"
        layout="prev, pager, next, sizes"
        :current-page="currentPage"
        :page-size="query.limit"
        :total="loginCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import { today, yesterday } from "@/helpers/util";

function formatBegin(begin) {
  return begin.toISOString();
}
function formatEnd(end) {
  return new Date(end.getTime() + 24 * 3600 * 1000 - 1).toISOString();
}

export default {
  name: "Logins",
  data() {
    const pageSizes = [10, 20, 30, 50];
    const defaultDateRange = [yesterday(), today()];
    return {
      dateRange: defaultDateRange,
      query: {
        begin: formatBegin(defaultDateRange[0]),
        end: formatEnd(defaultDateRange[1]),
        offset: 0,
        limit: pageSizes[0],
        account: "",
        order: "-createdAt"
      }
    };
  },
  computed: {
    currentPage() {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    },
    ...mapState({
      loginCount: state => state.user.loginList.count,
      processing: state => state.user.loginListProcessing,
      logins: state => state.user.loginList.data || []
    })
  },
  watch: {
    dateRange(value) {
      if (!value) {
        this.query.begin = "";
        this.query.end = "";
        return;
      }
      this.query.begin = formatBegin(value[0]);
      this.query.end = formatEnd(value[1]);
    }
  },
  methods: {
    ...mapActions(["listUserLogins"]),
    handleCurrentChange(page) {
      this.query.offset = (page - 1) * this.query.limit;
      this.fetch();
    },
    handleSizeChange(pageSize) {
      this.query.limit = pageSize;
      this.query.offset = 0;
      this.fetch();
    },
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      try {
        await this.listUserLogins(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    handleSearch() {
      this.query.offset = 0;
      this.fetch();
    }
  },
  beforeMount() {
    this.fetch();
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.logins
  margin: $mainMargin
  i
    margin-right: 5px
  .dateRange, .submit
    width: 100%
.pagination
  text-align: right
</style>
