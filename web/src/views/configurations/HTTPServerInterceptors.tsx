import { useMessage } from "naive-ui";
import { defineComponent, onMounted } from "vue";

import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../stores/configs";
import {
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";
import { useCommonStore } from "../../stores/common";
import { storeToRefs } from "pinia";

export default defineComponent({
  name: "HTTPServerInterceptorConfigs",
  setup() {
    const message = useMessage();
    const commonStore = useCommonStore();
    const { routers, fetchingRouters } = storeToRefs(commonStore);

    onMounted(async () => {
      try {
        await commonStore.listRouter();
      } catch (err) {
        showError(message, err);
      }
    });

    return {
      routers,
      processing: fetchingRouters,
    };
  },
  render() {
    const { routers, processing } = this;
    if (processing) {
      return <ExLoading />;
    }
    const extraFormItems: FormItem[] = [
      {
        name: "路由：",
        key: "data.router",
        type: FormItemTypes.Select,
        placeholder: "请选择路由",
        options: routers.map((item) => {
          const value = `${item.method} ${item.route}`;
          return {
            label: value,
            value,
          };
        }),
      },
      {
        name: "IP：",
        span: 24,
        key: "data.ip",
        placeholder: "请填写需要拦截的IP，如果未指定则表示所有IP均处理",
      },
      {
        name: "前置脚本：",
        key: "data.before",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理前的相关处理脚本",
      },
      {
        name: "后置脚本：",
        key: "data.after",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理后的相关处理脚本",
      },
    ];
    const rules = getDefaultFormRules({
      "data.router": newRequireRule("路由不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="HTTP服务拦截配置"
        editorTitle="添加/更新HTTP服务拦截配置"
        editorDescription="设置HTTP服务各路由的拦截配置"
        category={ConfigCategory.HTTPServerInterceptor}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
