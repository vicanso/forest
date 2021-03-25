<template lang="pug">
el-tooltip(
  v-if="value"
  placement="bottom"
)
  template(
    #content
  )
    ul
      li(
        v-if="file"
      ) file: {{file}}
      li(
        v-if="line"
      ) line: {{line}}
      li(
        v-if="status"
      ) status: {{status}}
      li(
        v-if="category"
      ) category: {{category}}
  span
    | {{value}}
    i.el-icon-info
span(
  v-else
) --
</template>

<script lang="ts">
import { defineComponent } from "vue";

function getValue(reg, message) {
  const arr = reg.exec(message);
  if (!arr || !arr[1]) {
    return "";
  }
  return arr[1];
}

export default defineComponent({
  name: "HTTPErrorFormater",
  props: {
    message: {
      type: String,
      default: "",
    },
  },
  data() {
    const { message } = this.$props;
    let value = message || "";
    const keyword = "message=";
    const index = value.indexOf(keyword);
    if (index !== -1) {
      value = value.substring(index + keyword.length);
    }
    return {
      file: getValue(/file=(\S+),/, message),
      line: getValue(/line=(\S+),/, message),
      status: getValue(/statusCode=(\S+),/, message),
      category: getValue(/category=(\S+),/, message),
      value: value,
    };
  },
});
</script>

<style lang="stylus" scoped>
i
  margin-left: 3px
</style>
