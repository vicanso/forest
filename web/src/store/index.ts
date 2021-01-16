import { Store } from "vuex";

import { getUserStore, userStore } from "./modules/user";
import { getCommonStore, commonStore } from "./modules/common";
import { getTrackerStore, trackerStore } from "./modules/tracker";
import { getConfigStore, configStore } from "./modules/config";

const stores: Store<any>[] = [
  userStore,
  commonStore,
  trackerStore,
  configStore,
];

export const useUserStore = getUserStore;
export const useCommonStore = getCommonStore;
export const useTrackerStore = getTrackerStore;
export const useConfigStore = getConfigStore;

export default stores;
