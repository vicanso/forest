<template lang="pug">
.statsTable

  el-divider(
    v-if="$props.flux"
  )
    el-tooltip(
      placement="bottom"
    )
      template(
        #content
      )
        pre {{$props.flux}}
      span.flux 
        i.el-icon-connection
        span FLUX SCRIPT

  el-table(
    :data="items"
    row-key="_time"
    stripe
    :default-sort="{ prop: '_time', order: 'descending' }"
  )
    slot
  el-pagination.pagination(
    v-if="filterItems && filterItems.length > 0"
    layout="prev, pager, next, sizes"
    :current-page="Math.floor(offset / limit) + 1"
    :page-size="limit"
    :page-sizes="pageSizes"
    :total="filterItems.length"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  )
</template>

<script lang="ts">
import { defineComponent, inject } from "vue";
import { PAGE_SIZES } from "../constants/common";

export default defineComponent({
  name: "StatsTable",
  props: {
    data: {
      type: Array,
      default: () => [],
    },
    flux: {
      type: String,
      default: "",
    },
  },
  setup() {
    return {
      statsParams: inject("statsParams"),
    };
  },
  data() {
    const { data } = this.$props;
    const originalItems = data.slice(0);
    originalItems.sort((item1: { _time: string }, item2: { _time: string }) => {
      const date1 = new Date(item1._time);
      const date2 = new Date(item2._time);
      return date2.getTime() - date1.getTime();
    });
    return {
      originalItems: originalItems,
      filterItems: originalItems.slice(0),
      offset: 0,
      limit: 10,
      pageSizes: PAGE_SIZES,
    };
  },
  computed: {
    items() {
      const { offset, limit, filterItems } = this;
      return filterItems.slice(offset, offset + limit);
    },
  },
  watch: {
    "statsParams.filters": function (filters) {
      this.offset = 0;
      const keys = Object.keys(filters);
      if (keys.length === 0) {
        this.filterItems = this.originalItems.slice(0);
        return;
      }
      // 根据filter的字段筛选
      this.filterItems = this.originalItems.filter((item) => {
        let matched = true;
        keys.forEach((key) => {
          if (!matched) {
            return;
          }
          // - 表示为空
          if (filters[key] === "-") {
            if (item[key]) {
              matched = false;
            }
            return;
          }
          if (
            item[key] !== filters[key] &&
            item[key] !== Number(filters[key])
          ) {
            matched = false;
          }
        });
        return matched;
      });
    },
  },
  methods: {
    handleCurrentChange(page: number): void {
      this.offset = (page - 1) * this.limit;
    },
    handleSizeChange(pageSize: number): void {
      this.offset = 0;
      this.limit = pageSize;
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";

.pagination
  text-align right
  margin-top 15px
.flux
  color $darkGray
  i
    margin-right 5px
    font-weight bold
</style>
