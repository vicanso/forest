<template lang="pug">
//- 随机字符串生成
mixin RandomKeys
  el-card.card
    template(
      #header
    )
      span 随机字符串生成
    .randomKeys(
      v-loading="randomKeys.processing"      
    )
      ul(
        v-if="randomKeys.items.length"
      ): li(
        v-for="item in randomKeys.items"
      ) {{item}}
      p.tips(
        v-else
      ) 点击按钮生成随机字符串
    ex-button(
      :onClick="listRandomKey"
    ) 生成随机字符串

//- 缓存查询与删除
mixin CacheQueryDelete
  el-card(
    v-loading="cacheProcessing"
  )
    template(
      #header
    )
      span 缓存查询与删除
    p.tips session的缓存格式 ss:sessionID
    el-row(
      :gutter="15"
    )
      el-col(
        :span="12"
      ): el-input(
        v-model="cacheKey"
        placeholder="请输入缓存的key"
        clearable
      )
     
      el-col(
        :span="6"
      ): ex-button(
        :onClick="findCache"
      ) 查询
      el-col(
        :span=6
        category="smallText"
      ): ex-button(
        buttonCategory="default"
        :onClick="cleanCache"
      ) 清除

    pre.cacheData(
      v-if="cacheData"
    ) {{cacheData}}

.others
  //- 随机字符串
  +RandomKeys
  //- 缓存查询与删除
  +CacheQueryDelete

</template>

<script lang="ts">
import { defineComponent } from "vue";

import ExButton from "../components/ExButton.vue";
import useCommonState, { commonListRandomKey } from "../states/common";
import { adminFindCacheByKey, adminCleanCacheByKey } from "../states/admin";

export default defineComponent({
  name: "Others",
  components: {
    ExButton,
  },
  setup() {
    const commonState = useCommonState();
    return {
      randomKeys: commonState.randomKeys,
      listRandomKey: commonListRandomKey,
    };
  },
  data() {
    return {
      cacheKey: "",
      cacheProcessing: false,
      cacheData: "",
    };
  },
  methods: {
    // 查询缓存
    async findCache() {
      const { cacheKey, cacheProcessing } = this;
      if (!cacheKey) {
        this.$message.warning("缓存key不能为空");
        return;
      }
      if (cacheProcessing) {
        return;
      }
      try {
        this.cacheProcessing = true;
        const resp = await adminFindCacheByKey(cacheKey);
        if (!resp.data) {
          this.$message.warning("请缓存不存在，请确认是否正确");
          return;
        }
        this.cacheData = resp.data;
        try {
          const data = JSON.parse(resp.data);
          this.cacheData = JSON.stringify(data, null, 2);
        } catch (err) {
          // 如果json处理出错，忽略
        }
      } catch (err) {
        this.$error(err);
      } finally {
        this.cacheProcessing = false;
      }
    },
    // 清除缓存
    async cleanCache() {
      const { cacheKey, cacheProcessing } = this;
      if (!cacheKey) {
        this.$message.warning("key不能为空");
        return;
      }
      if (cacheProcessing) {
        return;
      }
      try {
        this.cacheProcessing = true;
        await adminCleanCacheByKey(cacheKey);
        this.cacheData = "";
      } catch (err) {
        this.$error(err);
      } finally {
        this.cacheProcessing = false;
      }
    },
  },
});
</script>
<style lang="stylus" scoped>
@import "../common";

.others
  margin $mainMargin

.card
  margin-bottom $mainMargin
.randomKeys
  line-height 2em
  margin 0
  padding-bottom 10px
  height 40px
  li
    display inline-block
    margin-right 15px
.btn
  width 100%

.tips
  color $darkGray
  font-size 14px
  line-height 2em

.cacheData
  margin $mainMargin
</style>
