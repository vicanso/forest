import { create } from "naive-ui";
import { createApp } from "vue";
import { createPinia } from "pinia";
import Root from "./Root";
import router from "./routes/router";
import { settingStorage } from "./storages/local";

const naive = create();
const app = createApp(Root);
app.use(router).use(naive).use(createPinia());

const loadSetting = async () => {
  try {
    await settingStorage.load();
  } catch (err) {
    console.error(err);
  }
};

router
  .isReady()
  .then(loadSetting)
  .then(() => app.mount("#app"))
  .catch(console.error);
