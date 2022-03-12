import { ChartBar, Cogs, Deezer, User } from "@vicons/fa";
import { css } from "@linaria/core";
import { NButton, NIcon, NMenu } from "naive-ui";
import { Component, defineComponent, h } from "vue";
import { containsAny } from "./helpers/util";
import { goTo, goToLogin } from "./routes";
import { names } from "./routes/routes";
import { useUserStore } from "./stores/user";
import { storeToRefs } from "pinia";
import { useCommonStore } from "./stores/common";

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const loginButtonClass = css`
  margin: 50px auto;
  text-align: center;
`;

const navigationOptions = [
  {
    label: "用户",
    key: "user",
    icon: renderIcon(User),
    children: [
      {
        label: "用户列表",
        key: names.users,
      },
      {
        label: "登录记录",
        key: names.logins,
      },
    ],
  },
  {
    label: "统计",
    key: "stats",
    icon: renderIcon(ChartBar),
    children: [
      {
        label: "用户行为",
        key: names.userTrackers,
      },
      {
        label: "响应出错记录",
        key: names.httpErrors,
      },
      {
        label: "后端HTTP调用",
        key: names.requests,
      },
    ],
  },
  {
    label: "配置",
    key: "settings",
    disabled: true,
    icon: renderIcon(Cogs),
    children: [
      {
        label: "所有配置",
        key: names.configs,
      },
      {
        label: "MockTime配置",
        key: names.mockTime,
      },
      {
        label: "黑名单IP",
        key: names.blockIPs,
      },
      {
        label: "SignedKey配置",
        key: names.signedKeys,
      },
      {
        label: "路由Mock配置",
        key: names.routerMocks,
      },
      {
        label: "路由并发配置",
        key: names.routerConcurrencies,
      },
      {
        label: "HTTP实例并发配置",
        key: names.requestConcurrencies,
      },
      {
        label: "HTTP服务拦截配置",
        key: names.httpServerInterceptors,
      },
      {
        label: "接收邮箱配置",
        key: names.emails,
      },
    ],
  },
  {
    label: "其它",
    key: "others",
    disabled: true,
    icon: renderIcon(Deezer),
    children: [
      {
        label: "缓存",
        key: names.caches,
      },
    ],
  },
];

export default defineComponent({
  name: "AppNavigation",
  setup() {
    const userStore = useUserStore();
    const { processing, account, roles } = storeToRefs(userStore);
    const { setting } = storeToRefs(useCommonStore());
    return {
      setting,
      processing,
      account,
      roles,
      handleNavigation(key: string): void {
        goTo(key, {
          replace: false,
        });
      },
    };
  },
  render() {
    const { account, processing, roles, $router, setting } = this;
    if (processing) {
      return <p class="tac">...</p>;
    }
    if (!account) {
      if (setting.collapsed) {
        return <div />;
      }
      return (
        <div class={loginButtonClass}>
          <NButton type="info" onClick={() => goToLogin()}>
            立即登录
          </NButton>
        </div>
      );
    }
    const options = navigationOptions.slice(0);
    if (containsAny(roles, ["su", "admin"])) {
      options.forEach((item) => {
        if (item.disabled) {
          item.disabled = false;
        }
      });
    }
    const currentRoute = $router.currentRoute.value.name?.toString();
    return (
      <NMenu
        value={currentRoute}
        defaultExpandAll={true}
        onUpdate:value={this.handleNavigation}
        options={options}
        collapsedWidth={64}
        collapsed={setting.collapsed}
      />
    );
  },
});
