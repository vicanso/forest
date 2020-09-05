<template>
  <el-upload
    class="upload"
    list-type="picture"
    drag
    :action="action"
    :on-success="handleSuccess"
    :on-error="handleError"
    :on-remove="handleRemove"
    :before-upload="handleBeforeUpload"
    :on-exceed="handleExceed"
    :limit="$props.limit"
    :file-list="fileList"
  >
    <i class="el-icon-upload"></i>
    <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
  </el-upload>
</template>
<script>
import { contains } from "@/helpers/util";
import { FILES_IMAGES } from "@/constants/url";
export default {
  props: {
    bucket: {
      type: String,
      default: "origin-pics"
    },
    limit: {
      type: Number,
      default: 1
    },
    files: Array,
    contentType: String,
    width: Number,
    height: Number
  },
  data() {
    const { files, bucket, width, height } = this.$props;
    let fileList = [];
    if (files) {
      fileList = files.slice(0);
    }
    const params = [`bucket=${bucket}`];
    if (width) {
      params.push(`width=${width}`);
    }
    if (height) {
      params.push(`height=${height}`);
    }
    return {
      fileList,
      action: FILES_IMAGES + `?${params.join("&")}`
    };
  },
  methods: {
    handleRemove(file) {
      const fileList = this.fileList.filter(item => {
        if (file.url && item.url === file.url) {
          return false;
        }
        return item.url !== file.response.url;
      });
      this.fileList = fileList;
      this.$emit("change", this.fileList.slice(0));
    },
    handleSuccess(file) {
      this.fileList.push(file);
      this.$emit("change", this.fileList.slice(0));
    },
    handleError(err) {
      this.$message.error(err.message);
    },
    handleBeforeUpload(file) {
      const { contentType } = this.$props;
      const validTypes = [];
      // 暂时仅支持图片类上传
      if (!contentType) {
        validTypes.push("image/jpeg", "image/png");
      } else {
        validTypes.push(contentType);
      }
      if (!contains(validTypes, file.type)) {
        this.$message.warning(`仅支持上传${validTypes.join("，")}格式`);
        return false;
      }
      const tooLarge = file.size / 1024 / 1024 > 1;

      if (tooLarge) {
        this.$message.error("上传图片不能超过1MB");
        return false;
      }
      return true;
    },
    handleExceed() {
      const { limit } = this.$props;
      this.$message.warning(`图片限制上传${limit}张，请先删除图片再上传`);
    }
  }
};
</script>
