<template lang="pug">
el-popover(
  placement="bottom-end"
  trigger="click"
  :width="350"
)
  .statsInfo(
    v-if="stats && stats.length"
    v-for="item in stats"
    :key="item.name"
  )
    h5 {{item.name.toUpperCase()}}
    ul
      li(
        v-for="subItem in item.items"
      )
        span.count {{subItem.value}}
        span {{subItem.key}}
  span(
    v-else
  ) 无汇总数据
  template(
    #reference
  ): el-button(
    type="text"
    size="small"
    icon="el-icon-s-grid"
  ) 查看汇总(TOP 5)
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "StatsSummary",
  props: {
    data: {
      type: Array,
      default: () => [],
    },
    fields: {
      type: Array,
      required: true,
    },
  },
  data() {
    const { data, fields } = this.$props;
    const stats = [];
    const result = {};
    const topCount = 5;
    (data || []).forEach((item: { field: string }) => {
      fields.forEach((field: string) => {
        if (!result[field]) {
          result[field] = {};
        }
        const subResult = result[field];
        const value = item[field] || "-";
        const count = subResult[value] || 0;
        subResult[value] = count + 1;
      });
    });
    Object.keys(result).forEach((key) => {
      const tmp = result[key];
      const arr = [];
      Object.keys(tmp).forEach((subKey) => {
        arr.push({
          key: subKey,
          value: tmp[subKey],
        });
      });
      // 如果无数据
      if (!arr.length) {
        return;
      }

      // 排序
      arr.sort((item1, item2) => item2.value - item1.value);

      // 仅显示top N的
      if (arr.length > topCount) {
        arr.length = topCount;
      }

      stats.push({
        name: key,
        items: arr,
      });
    });
    return {
      stats,
    };
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";

.statsInfo
  margin-top 15px
  &:first-child
    margin-top 0
ul
  list-style-position inside
li
  padding 3px 5px 3px 0
  &:hover
    background-color $lightBlue
.count
  float right
</style>
