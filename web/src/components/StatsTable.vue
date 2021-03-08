<template lang="pug">
.statsTable
  el-table(
    :data="items"
    row-key="_time"
    stripe
    :default-sort="{ prop: '_time', order: 'descending' }"
  )
    slot
  el-pagination.pagination(
    v-if="$props.count > 0"
    layout="prev, pager, next, sizes"
    :current-page="Math.floor(offset / limit) + 1"
    :page-size="limit"
    :page-sizes="pageSizes"
    :total="$props.count"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  )
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { PAGE_SIZES } from "../constants/common";

export default defineComponent({
  name: "StatsTable",
  props: {
    data: {
      type: Array,
      default: () => [],
    },
    count: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      offset: 0,
      limit: 10,
      pageSizes: PAGE_SIZES,
    };
  },
  computed: {
    items() {
      const { data } = this.$props;
      const { offset, limit } = this;
      const originalItems = data.slice(0);
      originalItems.sort(
        (item1: { _time: string }, item2: { _time: string }) => {
          const date1 = new Date(item1._time);
          const date2 = new Date(item2._time);
          return date2.getTime() - date1.getTime();
        }
      );
      return originalItems.slice(offset, offset + limit);
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
.pagination
  text-align right
  margin-top 15px
</style>
