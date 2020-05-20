<template>
  <div class="routerData">
    <el-col :span="8">
      <el-form-item class="hidden">
        <el-input />
      </el-form-item>
    </el-col>
    <el-col :span="8">
      <el-form-item label="路由选择：">
        <el-select
          class="selector"
          v-model="form.router"
          placeholder="请选择路由"
          v-loading="processing"
        >
          <el-option
            v-for="item in routers"
            :key="item.key"
            :label="item.key"
            :value="item.key"
          >
          </el-option>
        </el-select>
      </el-form-item>
    </el-col>
    <el-col :span="8">
      <el-form-item label="状态码：">
        <el-input
          v-model="form.status"
          type="number"
          placeholder="请输入响应状态码"
        />
      </el-form-item>
    </el-col>
    <el-col :span="8">
      <el-form-item label="响应类型：">
        <el-select
          class="selector"
          v-model="form.contentType"
          placeholder="请选择响应类型"
        >
          <el-option
            v-for="item in contentTypeList"
            :key="item"
            :label="item"
            :value="item"
          >
          </el-option>
        </el-select>
      </el-form-item>
    </el-col>
    <el-col :span="24">
      <el-form-item label="响应数据：">
        <el-input
          v-model="form.response"
          type="textarea"
          :autosize="{ minRows: 5, maxRows: 10 }"
          placeholder="请按选择的响应类型输入对应的响应数据"
        />
      </el-form-item>
    </el-col>
  </div>
</template>
<script>
import { mapState, mapActions } from "vuex";
export default {
  name: "RouterData",
  props: {
    data: String
  },
  computed: mapState({
    processing: state => state.common.processing,
    routers: state => state.common.routers || []
  }),
  data() {
    const form = {
      router: "",
      status: null,
      contentType: "",
      response: ""
    };
    if (this.$props.data) {
      const data = JSON.parse(this.$props.data);
      data.router = `${data.method} ${data.route}`;
      if (data.response && data.response[0] === "{") {
        data.response = JSON.stringify(JSON.parse(data.response), null, 2);
      }
      Object.assign(form, data);
    }
    return {
      contentTypeList: [
        "application/json; charset=UTF-8",
        "text/plain; charset=UTF-8",
        "text/html; charset=UTF-8"
      ],
      form
    };
  },
  watch: {
    "form.router": function() {
      this.handleChange();
    },
    "form.status": function() {
      this.handleChange();
    },
    "form.contentType": function() {
      this.handleChange();
    },
    "form.response": function() {
      this.handleChange();
    }
  },
  methods: {
    ...mapActions(["listRouter"]),
    handleChange() {
      const { router, status, contentType, response } = this.form;
      let value = "";
      if (router && status && contentType && response) {
        const [method, route] = router.split(" ");
        value = JSON.stringify({
          route,
          method,
          status: Number(status),
          contentType,
          response: response.trim()
        });
      }
      this.$emit("change", value);
    }
  },
  async beforeMount() {
    try {
      await this.listRouter();
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
.selector
  width: 100%
.hidden
  visibility: hidden
</style>
