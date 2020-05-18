<template>
  <div class="configurationEditor">
    <el-card>
      <div slot="header">
        <i class="el-icon-s-tools" />
        <span>{{ $props.title || "添加/更新配置" }}</span>
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
                :disabled="!!defaultValue.name"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="分类：">
              <el-input
                placeholder="请输入配置分类（可选）"
                v-model="form.category"
                clearable
                :disabled="!!defaultValue.category"
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
                v-model="form.beginDate"
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
                v-model="form.endDate"
                type="datetime"
                placeholder="选择日期时间"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
          <MockTimeData
            :data="form.data"
            v-if="$props.category === catMockTime"
            @change="handleChange"
          />
          <el-col :span="24">
            <el-form-item>
              <el-button class="submit" type="primary" @click="submit">{{
                submitText
              }}</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { MOCK_TIME } from "@/constants/config";
import MockTimeData from "@/components/configs/MockTimeData.vue";
import { diff } from "@/helpers/util";

export default {
  name: "ConfigEditor",
  props: {
    category: {
      type: String,
      required: true
    },
    title: String,
    summary: String,
    id: Number
  },
  components: {
    MockTimeData
  },
  computed: mapState({
    processing: state => state.config.processing
  }),
  data() {
    const defaultValue = {};
    const { $props } = this;
    // 如果是mock time，则固定名字与分类
    switch ($props.category) {
      case MOCK_TIME:
        defaultValue.name = MOCK_TIME;
        defaultValue.category = MOCK_TIME;
        break;
      default:
        break;
    }
    const submitText = $props.id ? "更新" : "提交";
    return {
      originalValue: null,
      fetching: false,
      catMockTime: MOCK_TIME,
      defaultValue,
      submitText,
      status: [
        {
          label: "启用",
          value: 1
        },
        {
          label: "禁用",
          value: 2
        }
      ],
      form: {
        name: defaultValue.name || "",
        category: defaultValue.category || "",
        status: null,
        beginDate: "",
        endDate: "",
        data: ""
      }
    };
  },
  methods: {
    ...mapActions(["addConfig", "getConfigByID", "updateConfigByID"]),
    handleChange(data) {
      this.form.data = data;
    },
    async submit() {
      const { name, category, status, beginDate, endDate, data } = this.form;
      if (!name || !status || !beginDate || !endDate || !data) {
        this.$message.warning("名称、状态、开始结束日期以及配置数据均不能为空");
        return;
      }
      const { id } = this.$props;
      try {
        const config = {
          name,
          status,
          category,
          beginDate,
          endDate,
          data
        };
        if (beginDate.toISOString) {
          config.beginDate = beginDate.toISOString();
        }
        if (endDate.toISOString) {
          config.endDate = endDate.toISOString();
        }
        // 更新
        if (id) {
          const info = diff(config, this.originalValue);
          if (info.modifiedCount) {
            await this.updateConfigByID({
              id,
              data: info.data
            });
          }
          this.$message.info("修改配置成功");
        } else {
          await this.addConfig(config);
          this.$message.info("添加配置成功");
        }
        this.$router.back();
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    const { id } = this.$props;
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
