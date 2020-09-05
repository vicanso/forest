<template>
  <el-card class="baseEditor" v-loading="processing">
    <div slot="header">
      <i v-if="$props.icon" :class="$props.icon"></i>
      <span>{{ $props.title }}</span>
    </div>
    <el-form
      v-if="inited"
      :label-width="$props.labelWidth"
      ref="baseEditorForm"
      :rules="rules"
      :model="current"
    >
      <el-row :gutter="15">
        <el-col
          v-for="field in $props.fields"
          :span="field.span || 8"
          :key="field.key"
        >
          <el-form-item
            :label="field.label"
            :label-width="field.labelWidth"
            :class="field.itemClass"
            :prop="field.key"
          >
            <el-select
              class="select"
              v-if="field.type === 'select'"
              :placeholder="field.placeholder"
              v-model="current[field.key]"
              :multiple="field.multiple || false"
            >
              <el-option
                v-for="item in field.options"
                :key="item.key || item.value"
                :label="item.label || item.name"
                :value="item.value"
              />
            </el-select>
            <el-input
              v-else-if="field.type === 'specsUnit'"
              :placeholder="field.placeholder"
              v-model="current[field.key]"
              :clearable="field.clearable"
            >
              <el-select
                class="inputSelect"
                v-model="current[field.selectKey]"
                slot="append"
                :placeholder="field.selectPlaceholder"
              >
                <el-option
                  v-for="item in field.options"
                  :key="item.name"
                  :label="item.name"
                  :value="item.value"
                />
              </el-select>
            </el-input>

            <el-input
              type="textarea"
              v-else-if="field.type === 'textarea'"
              v-model="current[field.key]"
              :placeholder="field.placeholder"
              :autosize="field.autosize"
            />
            <Upload
              v-else-if="field.type === 'upload'"
              :files="current[field.key]"
              :bucket="field.bucket"
              :limit="field.limit"
              :contentType="field.contentType"
              :width="field.width"
              :height="field.height"
              @change="handleUpload"
            />
            <el-date-picker
              v-else-if="field.type === 'datePicker'"
              v-model="current[field.key]"
              :type="field.pickerType || 'date'"
              :placeholder="field.placeholder"
            />
            <el-input
              v-else
              v-model="current[field.key]"
              :clearable="field.clearable"
              :disabled="field.disabled || false"
              :placeholder="field.placeholder"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" type="primary" @click="submit">{{
              submitText
            }}</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" @click="goBack">返回</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </el-card>
</template>
<script>
import { diff, validateForm, omitNil, getFieldRules } from "@/helpers/util";
import Upload from "@/components/Upload.vue";

export default {
  name: "BaseEditor",
  components: {
    Upload
  },
  props: {
    icon: String,
    title: {
      type: String,
      required: true
    },
    labelWidth: {
      type: String,
      default: "80px"
    },
    fields: {
      type: Array,
      required: true
    },
    id: Number,
    findByID: Function,
    updateByID: Function,
    add: Function
  },
  data() {
    const { id, fields } = this.$props;
    const submitText = id ? "更新" : "添加";
    const current = {};
    fields.forEach(item => {
      current[item.key] = null;
    });

    return {
      inited: false,
      originData: null,
      processing: false,
      submitText,
      current,
      rules: getFieldRules(fields)
    };
  },
  methods: {
    handleUpload(files) {
      this.current.files = files;
    },
    async handleAdd(data) {
      const { add } = this.$props;
      const { rules } = this;
      this.processing = true;
      try {
        if (rules) {
          await validateForm(this.$refs.baseEditorForm);
        }
        await add(data);
        this.$message.info("已成功添加");
        this.goBack();
      } catch (err) {
        this.$message.error(err.message);
      } finally {
        this.processing = false;
      }
    },
    async handleUpdate(data) {
      const { id, updateByID } = this.$props;
      const { originData, rules } = this;
      const updateInfo = diff(omitNil(data), originData);
      if (updateInfo.modifiedCount === 0) {
        this.$message.warning("请先修改要更新的信息");
        return;
      }

      this.processing = true;
      try {
        if (rules) {
          await validateForm(this.$refs.baseEditorForm);
        }
        await updateByID({
          id,
          data: updateInfo.data
        });
        this.$message.info("已成功更新");
        this.goBack();
      } catch (err) {
        this.$message.error(err.message);
      } finally {
        this.processing = false;
      }
    },
    submit() {
      const { current } = this;
      const { id, fields } = this.$props;
      const data = Object.assign({}, current);
      fields.forEach(item => {
        if (item.dataType === "number") {
          data[item.key] = Number(data[item.key]);
        }
      });
      if (!id) {
        this.handleAdd(data);
        return;
      }
      this.handleUpdate(data);
    },
    goBack() {
      this.$router.back();
    }
  },
  async beforeMount() {
    const { id, findByID } = this.$props;
    if (!id) {
      this.inited = true;
      return;
    }
    try {
      const data = await findByID(id);
      this.originData = data;
      Object.assign(this.current, data);
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.inited = true;
    }
  }
};
</script>
<style lang="sass" scoped>
.baseEditor
  i
    margin-right: 5px
  .select, .btn
    width: 100%
  .inputSelect
    min-width: 60px
</style>
