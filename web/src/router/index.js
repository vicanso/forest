import Vue from "vue";
import VueRouter from "vue-router";
import Home from "@/views/Home.vue";

import store from "@/store";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home
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
  // if (!store.state.user.account) {
  //   return next({
  //     name: ROUTE_LOGIN
  //   })
  // }
  return next();
});

export default router;
