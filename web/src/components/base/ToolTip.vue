<template>
  <el-tooltip v-if="$props.content">
    <span slot="content">{{ $props.content }}</span>
    <i v-if="$props.icon" :class="$props.icon" />
    <span v-else>{{
      $props.content.substring($props.content.length - $props.viewSize)
    }}</span>
  </el-tooltip>
  <span v-else>--</span>
</template>
<script>
export default {
  name: "BaseToolTip",
  props: {
    icon: String,
    content: String,
    viewSize: {
      type: Number,
      default: 5
    }
  }
};
</script>
