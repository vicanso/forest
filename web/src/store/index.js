import Vue from "vue";
import Vuex from "vuex";

import userStore from "@/store/modules/user";
import configStore from "@/store/modules/config";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    user: userStore,
    config: configStore
  }
});
