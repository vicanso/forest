import { createRouter, createWebHashHistory } from "vue-router";

import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Logins from "../views/Logins.vue";
import Users from "../views/Users.vue";
import Trackers from "../views/Trackers.vue";
import MockTime from "../views/configs/MockTime.vue";
import BlockIP from "../views/configs/BlockIP.vue";
import SignedKey from "../views/configs/SignedKey.vue";
import RouterMock from "../views/configs/Router.vue";
import RouterConcurrency from "../views/configs/RouterConcurrency.vue";
import SessionInterceptor from "../views/configs/SessionInterceptor.vue";
import ValidConfiguration from "../views/configs/ValidConfiguration.vue";

const home = "home";
const login = "login";
const register = "register";
const logins = "logins";
const users = "users";
const trackers = "trackers";
const mockTime = "mockTime";
const blockIP = "blockIP";
const signedKey = "signedKey";
const routerMock = "routerMock";
const routerConcurrency = "routerConcurrency";
const sessionInterceptor = "sessionInterceptor";
const validConfiguration = "validConfiguration";

interface Location {
  name: string;
  path: string;
}

let currentLocation: Location = {
  name: "",
  path: "",
};
let prevLocation: Location = {
  name: "",
  path: "",
};

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: home,
      component: Home,
    },
    {
      path: "/login",
      name: login,
      component: Login,
    },
    {
      path: "/register",
      name: register,
      component: Register,
    },
    {
      path: "/users",
      name: users,
      component: Users,
    },
    {
      path: "/logins",
      name: logins,
      component: Logins,
    },
    {
      path: "/trackers",
      name: trackers,
      component: Trackers,
    },
    {
      path: "/mock-time",
      name: mockTime,
      component: MockTime,
    },
    {
      path: "/block-ip",
      name: blockIP,
      component: BlockIP,
    },
    {
      path: "/signed-key",
      name: signedKey,
      component: SignedKey,
    },
    {
      path: "/router-mock",
      name: routerMock,
      component: RouterMock,
    },
    {
      path: "/router-concurrency",
      name: routerConcurrency,
      component: RouterConcurrency,
    },
    {
      path: "/session-interceptor",
      name: sessionInterceptor,
      component: SessionInterceptor,
    },
    {
      path: "/valid-configuration",
      name: validConfiguration,
      component: ValidConfiguration,
    },
  ],
});

export function getHomeRouteName(): string {
  return home;
}
export function getLoginRouteName(): string {
  return login;
}
export function getRegisterRouteName(): string {
  return register;
}
export function getLoginsRouteName(): string {
  return logins;
}
export function getUsersRouteName(): string {
  return users;
}
export function getTrackersRouteName(): string {
  return trackers;
}
export function getMockTimeRouteName(): string {
  return mockTime;
}
export function getBlockIPRouteName(): string {
  return blockIP;
}
export function getSignedKeyRouteName(): string {
  return signedKey;
}
export function getRouterMockRouteName(): string {
  return routerMock;
}
export function getRouterConcurrencyRouteName(): string {
  return routerConcurrency;
}
export function getSessionInterceptorRouteName(): string {
  return sessionInterceptor;
}
export function getValidConfigurationRouteName(): string {
  return validConfiguration;
}

export function getCurrentLocation(): Location {
  return currentLocation;
}

router.beforeEach((to, from) => {
  if (from.name) {
    prevLocation.name = from.name.toString();
    prevLocation.path = from.fullPath;
  }
  if (to.name) {
    currentLocation.name = to.name.toString();
    currentLocation.path = to.fullPath;
  }
});

export default router;