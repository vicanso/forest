<template lang="pug">
config-editor(
  name="设置session的拦截提示信息"
  summary="注意：针对session拦截，用于将所有用户相关接口拦截处理（如系统维护等），配置时需要确保配置正确"
  :category="category"
  :defaultValue="defaultValue"
  :backDisabled="true"
  v-if="!processing"
  :id="currentID"
  :back="noop"
): template(
  #data="configProps"
): session-interceptor-data(
  :data="configProps.form.data"
  @change.self="configProps.form.data = $event"
)
</template>

<script lang="ts">
import { defineComponent } from "vue";

import ConfigEditor from "../../components/configs/Editor.vue";
import SessionInterceptorData from "../../components/configs/SessionInterceptorData.vue";
import { SESSION_INTERCEPTOR } from "../../constants/common";
import useConfigState, { configList } from "../../store/config";

export default defineComponent({
  name: "BlockIP",
  components: {
    SessionInterceptorData,
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
        name: SESSION_INTERCEPTOR,
        category: SESSION_INTERCEPTOR,
      },
      category: SESSION_INTERCEPTOR,
      processing: true,
      currentID: 0,
    };
  },
  async mounted() {
    const { $route, $router, configs } = this;
    this.processing = true;
    try {
      await configList({
        name: SESSION_INTERCEPTOR,
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
    // eslint-disable-next-line
    noop() {
      return;
    },
  },
});
</script>
<style lang="stylus" scoped>
@import "../../common";

.add
  margin: $mainMargin
.addBtn
  width: 100%
</style>
