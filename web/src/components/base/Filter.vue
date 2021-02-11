<template lang="pug">
//- 列表选择
mixin Select
  el-select.select(
    :placeholder="field.placeholder"
    v-model="current[field.key]"
    :multiple="field.multiple || false"
  )
    el-option(
      v-for="item in field.options"
      :key="item.key || item.value"
      :label="item.label || item.name"
      :value="item.value"
    )

//- 筛选按钮
mixin Filter
  ex-button(
    :onClick="doFilter"
    icon="el-icon-search"
  ) 筛选

//- 日期时间选择
mixin DateTimeRangePicker
  el-date-picker.dateRange.fullFill(
    v-model="current[field.key]"
    type="datetimerange"
    range-separator="至"
    start-placeholder="开始日期"
    end-placeholder="结束日期"
    :shortcuts="field.shortcuts"
  )

//- 日期选择
mixin DateRangePicker
  el-date-picker.dateRange.fullFill(
    v-model="current[field.key]"
    type="daterange"
    range-separator="至"
    start-placeholder="开始日期"
    end-placeholder="结束日期"
    :shortcuts="field.shortcuts"
  )

//- 数字输入
mixin NumberInput
  el-input(
    v-model="current[field.key]"
    type="number"
    :placeholder="field.placeholder"
    :default="field.defaultValue"
  )

//- 关键字输入
mixin KeywordInput
  el-input(
    @keyup.enter.native="doFilter"
    :clearable="field.clearable"
    v-model="current[field.key]"
    :disabled="field.disabled || false"
    :placeholder="field.placeholder"
  )

el-form.baseFilter(
  :label-width="$props.labelWidth"
): el-row(
  :gutter="15"
)
  el-col(
    v-for="field in $props.fields"
    :span="field.span || 8"
    :key="field.key"
  )
    el-form-item(
      :label="field.label"
      :label-width="field.labelWidth"
      :class="field.itemClass"
    )
      //- 列表选择
      template(
        v-if="field.type === 'select'"
      )
        +Select

      //- 点击筛选
      template(
        v-else-if="field.type === 'filter'"
      )
        +Filter
      
      //- 日期时间筛选
      template(
        v-else-if="field.type === 'dateTimeRange'"
      )
        +DateTimeRangePicker
      
      //- 日期筛选
      template(
        v-else-if="field.type === 'dateRange'"
      )
        +DateRangePicker

      //- 数字输入
      template(
        v-else-if="field.type === 'number'"
      )
        +NumberInput

      //- 关键字搜索
      template(
        v-else
      )
        +KeywordInput
</template>

<script lang="ts">
import { defineComponent } from "vue";

import ExButton from "../ExButton.vue";

export default defineComponent({
  name: "BaseFilter",
  components: {
    ExButton,
  },
  props: {
    labelWidth: {
      type: String,
      default: "90px",
    },
    fields: {
      type: Array,
      required: true,
    },
    filter: {
      type: Function,
      required: true,
    },
  },
  data() {
    const current = {};
    const { fields } = this.$props;
    fields.forEach((item) => {
      const { type, key, defaultValue } = item;
      if (type === "filter") {
        return;
      }
      current[key] = defaultValue || "";
    });
    return {
      processing: false,
      current,
    };
  },
  methods: {
    doFilter(): Promise<void> {
      return this.$props.filter(this.current);
    },
  },
});
</script>
<style lang="stylus" scoped>
.baseFilter
  .select, .btn, .dateRange
    width 100%
</style>
