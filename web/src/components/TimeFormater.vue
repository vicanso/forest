<template lang="pug">
span {{value}}
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { formatDate } from "../helpers/util";

export default defineComponent({
  name: "TimeFormater",
  props: {
    time: String,
  },
  data() {
    let value = "--";
    if (this.$props.time) {
      value = formatDate(this.$props.time);
    }
    return {
      value,
    };
  },
});
</script>
