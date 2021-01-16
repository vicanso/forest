import { createStore, Store } from "vuex";
import { Commit } from "vuex/types";
import request from "../../helpers/request";
import { FLUXES_TRACKERS, FLUXES_HTTP_ERRORS } from "../../constants/url";

const prefix = "flux";

const prefixTracker = `${prefix}.tracker`;
const mutationTrackerListProcessing = `${prefixTracker}.processing`;
const mutationTrackerList = `${prefixTracker}.list`;

const prefixHTTPError = `${prefix}.httpError`;
const mutationHTTPErrorListProcessing = `${prefixHTTPError}.processing`;
const mutationHTTPErrorList = `${prefixHTTPError}.list`;

interface TrackerList {
  processing: boolean;
  items: any[];
}
interface HTTPErrorList {
  processing: boolean;
  items: any[];
}

interface FluxState {
  trackers: TrackerList;
  httpErrors: HTTPErrorList;
}

const trackers: TrackerList = {
  processing: false,
  items: [],
};

const httpErrors: HTTPErrorList = {
  processing: false,
  items: [],
};

const state: FluxState = {
  trackers,
  httpErrors,
};

function fluxItemsSort(items: any[]) {
  return (items || [])
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

export const fluxStore = createStore<FluxState>({
  state,
  mutations: {
    // 设置正在查询用户行为
    [mutationTrackerListProcessing](state: FluxState, processing: boolean) {
      state.trackers.processing = processing;
    },
    // 设置用户行为记录
    [mutationTrackerList](state: FluxState, data: { trackers: any[] }) {
      state.trackers.items = fluxItemsSort(data.trackers);
    },
    // 设置正在查询出错列表
    [mutationHTTPErrorListProcessing](state: FluxState, processing: boolean) {
      state.httpErrors.processing = processing;
    },
    // 设置出错列表
    [mutationHTTPErrorList](state: FluxState, data: { httpErrors: any[] }) {
      state.httpErrors.items = fluxItemsSort(data.httpErrors);
    },
  },
  actions: {
    async listTracker(context: { commit: Commit }, params) {
      context.commit(mutationTrackerListProcessing, true);
      try {
        const { data } = await request.get(FLUXES_TRACKERS, {
          params,
        });
        context.commit(mutationTrackerList, data);
      } finally {
        context.commit(mutationTrackerListProcessing, false);
      }
    },
    async listHTTPError(context: { commit: Commit }, params) {
      context.commit(mutationHTTPErrorListProcessing, true);
      try {
        const { data } = await request.get(FLUXES_HTTP_ERRORS, {
          params,
        });
        context.commit(mutationHTTPErrorList, data);
      } finally {
        context.commit(mutationHTTPErrorListProcessing, false);
      }
    },
  },
});

// getFluxStore 获取flux store
export function getFluxStore(): Store<FluxState> {
  return fluxStore;
}
