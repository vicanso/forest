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

//- session查询清理
mixin SessionQueryDelete
  el-card(
    v-loading="sessionProcessing"
  )
    template(
      #header
    )
      span Session查询与清除
    el-row(
      :gutter="15"
    )
      el-col(
        :span="12"
      ): el-input(
        v-model="sessionID"
        placeholder="请输入session id"
        clearable
      )
      el-col(
        :span=6
        category="smallText"
      ): ex-button(
        buttonCategory="default"
        :onClick="cleanSession"
      ) 清除
      el-col(
        :span="6"
      ): ex-button(
        :onClick="findSession"
      ) 查询

    pre.sessionData(
      v-if="sessionData"
    ) {{sessionData}}

.others
  //- 随机字符串
  +RandomKeys
  //- session查询与删除
  +SessionQueryDelete

</template>

<script lang="ts">
import { defineComponent } from "vue";

import ExButton from "../components/ExButton.vue";
import useCommonState, { commonListRandomKey } from "../states/common";
import { adminFindSessionByID, adminCleanSessionByID } from "../states/admin";

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
      sessionID: "",
      sessionProcessing: false,
      sessionData: "",
    };
  },
  methods: {
    // 查询session
    async findSession() {
      const { sessionID, sessionProcessing } = this;
      if (!sessionID) {
        this.$message.warning("session id不能为空");
        return;
      }
      if (sessionProcessing) {
        return;
      }
      try {
        this.sessionProcessing = true;
        const resp = await adminFindSessionByID(sessionID);
        if (!resp.data) {
          this.$message.warning("该session信息不存的，请确认ID是否有误");
          return;
        }
        const sessionData = JSON.parse(resp.data);
        if (sessionData["user-session-info"]) {
          sessionData["user-session-info"] = JSON.parse(
            sessionData["user-session-info"]
          );
        }
        this.sessionData = JSON.stringify(sessionData, null, 2);
      } catch (err) {
        this.$error(err);
      } finally {
        this.sessionProcessing = false;
      }
    },
    // 清除session
    async cleanSession() {
      const { sessionID, sessionProcessing } = this;
      if (!sessionID) {
        this.$message.warning("session id不能为空");
        return;
      }
      if (sessionProcessing) {
        return;
      }
      try {
        this.sessionProcessing = true;
        await adminCleanSessionByID(sessionID);
        this.sessionData = "";
      } catch (err) {
        this.$error(err);
      } finally {
        this.sessionProcessing = false;
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

.sessionData
  margin $mainMargin
</style>
