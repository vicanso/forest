<template>
  <el-card class="logins">
    <div slot="header">
      <i class="el-icon-user-solid" />
      <span>用户登录查询</span>
    </div>
    <div v-loading="processing">
      <BaseFilter :fields="filterFields" @filter="filter" />
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
        <el-table-column prop="sessionID" key="sessionID" label="Session ID" />
        <el-table-column prop="trackID" key="trackID" label="Track ID" />
        <el-table-column label="Forwarded For" width="120">
          <template slot-scope="scope">
            <el-tooltip v-if="scope.row.xForwardedFor">
              <span slot="content">{{ scope.row.xForwardedFor }}</span>
              <i class="el-icon-info" />
            </el-tooltip>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column label="User Agent" width="120">
          <template slot-scope="scope">
            <el-tooltip v-if="scope.row.userAgent">
              <span slot="content">{{ scope.row.userAgent }}</span>
              <i class="el-icon-mobile-phone" />
            </el-tooltip>
            <span v-else>--</span>
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
        :page-sizes="pageSizes"
        :total="loginCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import { today, yesterday, formatBegin, formatEnd } from "@/helpers/util";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { PAGE_SIZES } from "@/constants/common";

const defaultDateRange = [yesterday(), today()];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
    span: 6
  },
  {
    label: "时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    defaultValue: defaultDateRange,
    span: 12
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 6
  }
];

export default {
  name: "Logins",
  extends: BaseTable,
  components: {
    BaseFilter
  },
  data() {
    return {
      filterFields,
      pageSizes: PAGE_SIZES,
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: PAGE_SIZES[0],
        account: "",
        order: "-createdAt"
      }
    };
  },
  computed: {
    ...mapState({
      loginCount: state => state.user.userLogins.count,
      processing: state => state.user.userLoginListProcessing,
      logins: state => state.user.userLogins.data || []
    })
  },
  methods: {
    ...mapActions(["listUserLogins"]),
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      const value = query.dateRange;
      if (value) {
        query.begin = formatBegin(value[0]);
        query.end = formatEnd(value[1]);
      } else {
        query.begin = "";
        query.end = "";
      }
      delete query.dateRange;
      try {
        await this.listUserLogins(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.logins
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
