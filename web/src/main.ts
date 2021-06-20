import { create } from "naive-ui";
import { createApp } from "vue";
import Root from "./Root";
import router from "./routes/router";

const naive = create();
const app = createApp(Root);
app.use(router).use(naive);

router
  .isReady()
  .then(() => app.mount("#app"))
  .catch(console.error);
