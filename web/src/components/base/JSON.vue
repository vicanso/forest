<template>
  <div class="jsonView">
    <pre v-if="json">{{ data }}</pre>
    <span v-else>{{ data }}</span>
  </div>
</template>
<script>
const formatJSON = data => {
  if (data && data[0] === "{") {
    return {
      json: true,
      data: JSON.stringify(JSON.parse(data), null, 2)
    };
  }
  return {
    json: false,
    data: data
  };
};
export default {
  name: "JSON",
  props: {
    content: String
  },
  data() {
    return formatJSON(this.$props.content);
  }
};
</script>

<style lang="sass" scoped>
@import "@/common.sass"
.jsonView
  overflow-x: scroll
</style>
