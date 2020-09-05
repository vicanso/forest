<template>
  <el-form class="baseFilter" :label-width="$props.labelWidth">
    <el-row :gutter="15">
      <el-col
        v-for="field in $props.fields"
        :span="field.span || 8"
        :key="field.key"
      >
        <el-form-item
          :label="field.label"
          :label-width="field.labelWidth"
          :class="field.itemClass"
        >
          <el-select
            class="select"
            v-if="field.type === 'select'"
            :placeholder="field.placeholder"
            v-model="current[field.key]"
            :multiple="field.multiple || false"
          >
            <el-option
              v-for="item in field.options"
              :key="item.key || item.value"
              :label="item.label || item.name"
              :value="item.value"
            />
          </el-select>
          <el-button
            v-else-if="field.type === 'filter'"
            :loading="processing"
            icon="el-icon-search"
            class="btn"
            type="primary"
            @click="filter"
            >筛选</el-button
          >
          <el-date-picker
            v-else-if="field.type === 'dateRange'"
            class="dateRange"
            v-model="current[field.key]"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          >
          </el-date-picker>
          <el-input
            v-else
            @keyup.enter.native="filter"
            :clearable="field.clearable"
            v-model="current[field.key]"
            :disabled="field.disabled || false"
            :placeholder="field.placeholder"
          />
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>
<script>
export default {
  name: "BaseFilter",
  props: {
    labelWidth: {
      type: String,
      default: "90px"
    },
    fields: {
      type: Array,
      required: true
    }
  },
  data() {
    const current = {};
    const { fields } = this.$props;
    fields.forEach(item => {
      if (item.type === "filter") {
        return;
      }
      current[item.key] = item.defaultValue || "";
    });
    return {
      processing: false,
      current
    };
  },
  methods: {
    filter() {
      this.$emit("filter", this.current);
    }
  }
};
</script>
<style lang="sass" scoped>
.baseFilter
  .select, .btn, .dateRange
    width: 100%
</style>
