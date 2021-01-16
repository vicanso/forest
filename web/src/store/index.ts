import { Store } from "vuex";

import { getUserStore, userStore } from "./modules/user";
import { getCommonStore, commonStore } from "./modules/common";
import { getFluxStore, fluxStore } from "./modules/flux";
import { getConfigStore, configStore } from "./modules/config";

const stores: Store<any>[] = [userStore, commonStore, fluxStore, configStore];

export const useUserStore = getUserStore;
export const useCommonStore = getCommonStore;
export const useFluxStore = getFluxStore;
export const useConfigStore = getConfigStore;

export default stores;
