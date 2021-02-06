<template lang="pug">
config-editor(
  name="添加/更新MockTime配置"
  summary="针对应用时间Mock，用于测试环境中调整应用时间"
  :category="category"
  :defaultValue="defaultValue"
  :backDisabled="true"
  v-if="!processing"
  :id="currentID"
  :back="noop"
): template(
  #data="configProps"
): mock-time-data(
  :data="configProps.form.data"
  @change="configProps.form.data = $event"
)
</template>
<script lang="ts">
import { defineComponent } from "vue";

import ConfigEditor from "../../components/configs/Editor.vue";
import MockTimeData from "../../components/configs/MockTimeData.vue";
import { MOCK_TIME } from "../../constants/common";
import useConfigState, { configList } from "../../states/config";

export default defineComponent({
  name: "MockTime",
  components: {
    MockTimeData,
    ConfigEditor,
  },
  setup() {
    const configState = useConfigState();
    return {
      configs: configState.configs,
    };
  },
  data() {
    return {
      defaultValue: {
        name: MOCK_TIME,
        category: MOCK_TIME,
      },
      processing: true,
      currentID: 0,
      category: MOCK_TIME,
    };
  },
  async mounted() {
    const { $route, $router, configs } = this;
    this.processing = true;
    try {
      await configList({
        name: MOCK_TIME,
      });
      const configurations = configs.items;
      if (configurations && configurations.length !== 0) {
        let currentID = null;
        if ($route.query.id) {
          currentID = Number($route.query.id);
        }
        if (currentID !== configurations[0].id) {
          $router.replace({
            query: {
              id: configurations[0].id,
            },
          });
        }
      }
    } catch (err) {
      this.$error(err);
    } finally {
      this.processing = false;
    }
  },
  methods: {
    // 无操作
    noop() {
      return;
    },
  },
});
</script>
