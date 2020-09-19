<template>
  <ConfigEditor
    name="设置session的拦截提示信息"
    summary="注意：针对session拦截，用于将所有用户相关接口拦截处理（如系统维护等），配置时需要确保配置正确"
    :category="category"
    :defaultValue="defaultValue"
    v-if="!processing"
    :id="currentID"
    :back="back"
  >
    <template v-slot:data="configProps">
      <SessionInterceptorData
        :data="configProps.form.data"
        @change="configProps.form.data = $event"
      />
    </template>
  </ConfigEditor>
</template>
<script>
import ConfigEditor from "@/components/configs/Editor.vue";
import SessionInterceptorData from "@/components/configs/SessionInterceptorData.vue";
import { SESSION_INTERCEPTOR } from "@/constants/config";
import { mapActions } from "vuex";
export default {
  name: "SessionInterceptor",
  components: {
    SessionInterceptorData,
    ConfigEditor
  },
  data() {
    return {
      defaultValue: {
        name: SESSION_INTERCEPTOR,
        category: SESSION_INTERCEPTOR
      },
      processing: false,
      currentID: 0,
      category: SESSION_INTERCEPTOR
    };
  },
  methods: {
    ...mapActions(["listConfig"]),
    back() {}
  },
  async beforeMount() {
    this.processing = true;
    try {
      const { configurations } = await this.listConfig({
        name: SESSION_INTERCEPTOR
      });
      if (configurations.length !== 0) {
        let currentID = null;
        if (this.$route.query.id) {
          currentID = Number(this.$route.query.id);
        }
        if (currentID !== configurations[0].id) {
          this.$router.replace({
            query: {
              id: configurations[0].id
            }
          });
        }
      }
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
