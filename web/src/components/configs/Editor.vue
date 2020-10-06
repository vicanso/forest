<template>
  <div class="configurationEditor">
    <el-card>
      <div slot="header">
        <i class="el-icon-s-tools" />
        <span>{{ $props.name || "添加/更新配置" }}</span>
      </div>
      <el-form label-width="90px" v-loading="processing" v-if="!fetching">
        <p>
          <i class="el-icon-info" />
          {{ $props.summary || "添加或更新系统配置信息" }}
        </p>
        <el-row :gutter="15">
          <el-col :span="8">
            <el-form-item label="名称：">
              <el-input
                placeholder="请输入配置名称"
                v-model="form.name"
                clearable
                :disabled="!!$props.defaultValue.name"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="分类：">
              <el-input
                placeholder="请输入配置分类（可选）"
                v-model="form.category"
                clearable
                :disabled="!!$props.defaultValue.category"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="是否启用：">
              <el-select
                class="selector"
                v-model="form.status"
                placeholder="请选择配置状态"
              >
                <el-option
                  v-for="item in status"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="开始时间：">
              <el-date-picker
                class="datePicker"
                v-model="form.startedAt"
                type="datetime"
                placeholder="选择日期时间"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="结束时间：">
              <el-date-picker
                class="datePicker"
                v-model="form.endedAt"
                type="datetime"
                placeholder="选择日期时间"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
          <slot :form="form" name="data"></slot>
          <el-col :span="primarySpan">
            <el-form-item>
              <el-button class="submit" type="primary" @click="submit">{{
                submitText
              }}</el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="!$props.backDisabled">
            <el-form-item>
              <el-button class="submit" @click="goBack">返回</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { diff } from "@/helpers/util";

export default {
  name: "ConfigEditor",
  props: {
    defaultValue: {
      type: Object,
      default: () => {
        return {};
      }
    },
    category: {
      type: String,
      required: true
    },
    name: String,
    summary: String,
    // 返回函数
    back: Function,
    backDisabled: Boolean
  },
  computed: {
    ...mapState({
      processing: state => state.config.processing,
      status: state => state.config.status
    }),
    id() {
      const { id } = this.$route.query;
      if (!id) {
        return 0;
      }
      return Number(id);
    }
  },
  data() {
    const { $props, $route } = this;
    const { defaultValue, backDisabled } = $props;
    const submitText = $route.query.id ? "更新" : "提交";
    const primarySpan = backDisabled ? 24 : 12;
    return {
      primarySpan,
      originalValue: null,
      fetching: false,
      submitText,
      form: {
        name: defaultValue.name || "",
        category: defaultValue.category || "",
        status: null,
        startedAt: "",
        endedAt: "",
        data: ""
      }
    };
  },
  methods: {
    ...mapActions(["addConfig", "getConfigByID", "updateConfigByID"]),
    async submit() {
      const { name, category, status, startedAt, endedAt, data } = this.form;
      if (!name || !status || !startedAt || !endedAt || !data) {
        this.$message.warning("名称、状态、开始结束日期以及配置数据均不能为空");
        return;
      }
      const { id } = this;
      try {
        const config = {
          name,
          status,
          category,
          startedAt,
          endedAt,
          data
        };
        if (startedAt.toISOString) {
          config.startedAt = startedAt.toISOString();
        }
        if (endedAt.toISOString) {
          config.endedAt = endedAt.toISOString();
        }
        // 更新
        if (id) {
          const info = diff(config, this.originalValue);
          if (!info.modifiedCount) {
            this.$message.warning("未修改配置无法更新");
            return;
          }
          await this.updateConfigByID({
            id,
            data: info.data
          });
          this.$message.info("修改配置成功");
        } else {
          await this.addConfig(config);
          this.$message.info("添加配置成功");
        }
        this.goBack();
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    goBack() {
      if (this.$props.back) {
        this.$props.back();
        return;
      }
      this.$router.back();
    }
  },
  async beforeMount() {
    const { id } = this;
    if (!id) {
      return;
    }
    // 如果指定了ID则为更新
    this.fetching = true;
    try {
      const data = await this.getConfigByID(id);
      const config = {};
      Object.keys(this.form).forEach(key => {
        config[key] = data[key];
      });
      this.originalValue = config;
      Object.assign(this.form, config);
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.fetching = false;
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.configurationEditor
  margin: $mainMargin
  i
    margin-right: 3px
  p
    color: $darkGray
    font-size: 13px
    margin: 0 0 15px 0
  .selector, .datePicker, .submit
    width: 100%
</style>
