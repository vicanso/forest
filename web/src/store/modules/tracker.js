import request from "@/helpers/request";
import { USERS_TRACKERS } from "@/constants/url";
import { formatDate } from "@/helpers/util";

const prefix = "tracker";

const mutationTrackerList = `${prefix}.list`;
const mutationTrackerListProcessing = `${mutationTrackerList}.processing`;

const state = {
  listProcessing: false,
  list: {
    data: null
  }
};

export default {
  state,
  mutations: {
    [mutationTrackerListProcessing](state, processing) {
      state.listProcessing = processing;
    },
    [mutationTrackerList](state, { trackers = [] }) {
      trackers.forEach(item => {
        item.timeDesc = formatDate(item._time);
      });
      state.list.data = trackers
        .sort((item1, item2) => {
          if (item1._time < item2._time) {
            return -1;
          }
          if (item1._time > item2._time) {
            return 1;
          }
          return 0;
        })
        .reverse();
    }
  },
  actions: {
    async listTracker({ commit }, params) {
      commit(mutationTrackerListProcessing, true);
      try {
        const { data } = await request.get(USERS_TRACKERS, {
          params
        });
        commit(mutationTrackerList, data);
      } finally {
        commit(mutationTrackerListProcessing, false);
      }
    }
  }
};
