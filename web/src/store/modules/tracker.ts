import { createStore, Store } from "vuex";
import { Commit } from "vuex/types";
import request from "../../helpers/request";
import { USERS_TRACKERS } from "../../constants/url";

const prefix = "tracker";
const mutationTrackerListProcessing = `${prefix}.processing`;
const mutationTrackerList = `${prefix}.list`;

interface TrackerList {
  processing: boolean;
  items: any[];
}

interface TrackerState {
  trackers: TrackerList;
}

const trackers: TrackerList = {
  processing: false,
  items: [],
};

const state: TrackerState = {
  trackers,
};

export const trackerStore = createStore<TrackerState>({
  state,
  mutations: {
    // 设置正在查询用户行为
    [mutationTrackerListProcessing](state: TrackerState, processing: boolean) {
      state.trackers.processing = processing;
    },
    // 设置用户行为记录
    [mutationTrackerList](state: TrackerState, data: { trackers: any[] }) {
      state.trackers.items = (data.trackers || [])
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
    },
  },
  actions: {
    async list(context: { commit: Commit }, params) {
      context.commit(mutationTrackerListProcessing, true);
      try {
        const { data } = await request.get(USERS_TRACKERS, {
          params,
        });
        context.commit(mutationTrackerList, data);
      } finally {
        context.commit(mutationTrackerListProcessing, false);
      }
    },
  },
});

// getTrackerStore 获取用户行为store
export function getTrackerStore(): Store<TrackerState> {
  return trackerStore;
}
