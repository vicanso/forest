<template>
  <el-col :span="8">
    <el-form-item label="时间配置：" class="mockTimeData">
      <el-date-picker
        class="datePicker"
        type="datetime"
        placeholder="选择要设置的时间"
        v-model="time"
      >
      </el-date-picker>
    </el-form-item>
  </el-col>
</template>
<script>
// mock日期的数据组件
export default {
  name: "MockTimeData",
  props: {
    data: String
  },
  data() {
    return {
      time: this.$props.data || ""
    };
  },
  watch: {
    time: function(value) {
      if (value) {
        value = value.toISOString();
      }
      this.$emit("change", value);
    }
  }
};
</script>
<style lang="sass" scoped>
.mockTimeData .datePicker
  width: 100%
</style>
