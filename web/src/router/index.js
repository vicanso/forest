import Vue from "vue";
import VueRouter from "vue-router";
import store from "@/store";
import {
  HOME,
  LOGIN,
  REGISTER,
  CONFIG_MOCK_TIME,
  CONFIG_BLOCK_IP,
  CONFIG_SIGNED_KEY,
  CONFIG_ROUTER
} from "@/constants/route";
import Home from "@/views/Home.vue";
import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";
import MockTime from "@/views/configs/MockTime.vue";
import BlockIP from "@/views/configs/BlockIP.vue";
import SignedKey from "@/views/configs/SignedKey.vue";
import Router from "@/views/configs/Router.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: HOME,
    component: Home
  },
  {
    path: "/login",
    name: LOGIN,
    component: Login
  },
  {
    path: "/register",
    name: REGISTER,
    component: Register
  },
  {
    path: "/configs/mockTime",
    name: CONFIG_MOCK_TIME,
    component: MockTime
  },
  {
    path: "/configs/blockIP",
    name: CONFIG_BLOCK_IP,
    component: BlockIP
  },
  {
    path: "/configs/signedKey",
    name: CONFIG_SIGNED_KEY,
    component: SignedKey
  },
  {
    path: "/configs/router",
    name: CONFIG_ROUTER,
    component: Router
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue")
  }
];

const router = new VueRouter({
  routes
});

let fetchedUserInfo = false;
function waitForFetchingUserInfo() {
  const check = resolve => {
    if (fetchedUserInfo) {
      resolve();
      return;
    }
    if (!store.state.user.processing) {
      fetchedUserInfo = true;
    }
    setTimeout(() => {
      check(resolve);
    }, 30);
  };

  return new Promise(resolve => {
    check(resolve);
  });
}

router.beforeEach(async (to, from, next) => {
  if (!fetchedUserInfo) {
    await waitForFetchingUserInfo();
  }
  if (!to.meta.requiresAuth) {
    return next();
  }
  if (!store.state.user.account) {
    return next({
      name: LOGIN
    });
  }
  return next();
});

export default router;
