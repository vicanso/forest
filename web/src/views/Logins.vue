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
        <el-table-column label="Session ID">
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.sessionID" />
          </template>
        </el-table-column>
        <el-table-column label="Track ID">
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.trackID" />
          </template>
        </el-table-column>
        <el-table-column label="Forwarded For" width="120">
          <template slot-scope="scope">
            <BaseToolTip
              icon="el-icon-info"
              :content="scope.row.xForwardedFor"
            />
          </template>
        </el-table-column>
        <el-table-column label="User Agent" width="120">
          <template slot-scope="scope">
            <BaseToolTip
              icon="el-icon-mobile-phone"
              :content="scope.row.userAgent"
            />
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
import { today, formatBegin, formatEnd } from "@/helpers/util";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import BaseToolTip from "@/components/base/ToolTip.vue";
import { PAGE_SIZES } from "@/constants/common";

const defaultDateRange = [today(), today()];
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
    BaseFilter,
    BaseToolTip
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
      const params = Object.assign({}, query);
      const value = params.dateRange;
      if (value) {
        params.begin = formatBegin(value[0]);
        params.end = formatEnd(value[1]);
      } else {
        params.begin = "";
        params.end = "";
      }
      delete params.dateRange;
      try {
        await this.listUserLogins(params);
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
