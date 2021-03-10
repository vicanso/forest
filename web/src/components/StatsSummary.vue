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
        @click="doFilter(item.name, subItem.key)"
        :class=`{
          selected: filters[item.name] === subItem.key
        }`
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
  emits: ["filter"],
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
      filters: {},
      stats,
    };
  },
  methods: {
    doFilter(key, value) {
      const { filters } = this;
      if (filters[key] === value) {
        delete filters[key];
      } else {
        filters[key] = value;
      }
      this.$emit("filter", Object.assign({}, filters));
    },
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
  padding 3px 5px
  cursor pointer
  &:hover
    background-color $lightBlue
  &.selected
    background-color $blue
    color $white
.count
  float right
</style>
