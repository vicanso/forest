<template>
  <ConfigEditor :category="category" v-if="!processing" :id="currentID" />
</template>
<script>
import ConfigEditor from "@/components/configs/Editor.vue";
import { MOCK_TIME } from "@/constants/config";
import { mapActions } from "vuex";

export default {
  name: "MockTime",
  components: {
    ConfigEditor
  },
  data() {
    return {
      processing: false,
      currentID: 0,
      category: MOCK_TIME
    };
  },
  methods: {
    ...mapActions(["listConfig"])
  },
  async beforeMount() {
    this.processing = true;
    try {
      const { configs } = await this.listConfig({
        name: MOCK_TIME
      });
      if (configs.length !== 0) {
        this.currentID = configs[0].id;
      }
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
