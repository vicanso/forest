<template>
  <div class="orderCommissionData">
    <el-col :span="12">
      <el-form-item label="佣金分组：">
        <el-select
          class="orderCommissionGroupSelect"
          placeholder="请输入佣金分组，通用的则选择 *"
          v-model="group"
        >
          <el-option
            v-for="item in marketingGroups"
            :key="item.value"
            :label="item.name"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
    </el-col>
    <el-col :span="12">
      <el-form-item label="佣金比例：">
        <el-input
          type="number"
          placeholder="请输入佣金比例（不大于0.1）"
          v-model="ratio"
        />
      </el-form-item>
    </el-col>
  </div>
</template>
<script>
import { mapActions } from "vuex";

const marketingGroups = [];

export default {
  name: "OrderCommissionData",
  props: {
    data: String
  },
  data() {
    const data = {
      marketingGroups,
      processing: false,
      group: "",
      ratio: 0
    };
    if (this.$props.data) {
      Object.assign(data, JSON.parse(this.$props.data));
    }
    return data;
  },
  watch: {
    group() {
      this.handleChange();
    },
    ratio() {
      this.handleChange();
    }
  },
  methods: {
    ...mapActions(["listUserMarketingGroup"]),
    handleChange() {
      const { group, ratio } = this;
      const ratioValue = Number(ratio);
      let value = "";
      if (ratioValue > 0.1) {
        this.$message.error("佣金比例不能大于0.1");
        return;
      }
      if (group) {
        value = JSON.stringify({
          group,
          ratio: ratioValue
        });
      }
      this.$emit("change", value);
    }
  },
  async beforeMount() {
    this.processing = true;
    try {
      const { marketingGroups } = await this.listUserMarketingGroup();
      this.marketingGroups.length = 0;
      this.marketingGroups.push({
        name: "*",
        value: "*"
      });
      this.marketingGroups.push(...marketingGroups);
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
<style lang="sass" scoped>
.orderCommissionGroupSelect
  width: 100%
</style>
