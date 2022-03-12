import { User } from "@vicons/fa";
import {
  NButton,
  NH2,
  NIcon,
  NLayoutHeader,
  NSpace,
  NText,
  useMessage,
} from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, onBeforeMount } from "vue";
import { storeToRefs } from "pinia";
import { mainHeaderHeight, padding } from "./constants/style";
import { showError } from "./helpers/util";
import { goToHome, goToLogin, goToRegister, goToProfile } from "./routes";
import { useUserStore } from "./stores/user";
import { useCommonStore } from "./stores/common";

const userInfoClass = css`
  margin-right: 5px;
`;
const headerClass = css`
  height: ${mainHeaderHeight}px;
  line-height: ${mainHeaderHeight}px;
  padding: 0 ${3 * padding}px;
`;

const logoClass = css`
  float: left;
  cursor: pointer;
`;

export default defineComponent({
  name: "AppHeader",
  setup() {
    const message = useMessage();
    const userStore = useUserStore();
    const commonStore = useCommonStore();
    const { anonymous, processing, account } = storeToRefs(userStore);
    const { setting } = storeToRefs(commonStore);
    onBeforeMount(async () => {
      try {
        await commonStore.getSetting();
        await userStore.fetch();
      } catch (err) {
        showError(message, err);
      } finally {
        if (anonymous.value) {
          goToLogin();
        }
      }
    });

    // 退出登录
    const logout = async () => {
      try {
        await userStore.logout();
      } catch (err) {
        showError(message, err);
      }
    };

    // 主题选择
    const renderToggleTheme = () => {
      const isDark = setting.value.theme === "dark";
      const toggleTheme = isDark ? "light" : "dark";
      const text = isDark ? "浅色" : "深色";
      return (
        <NButton
          bordered={false}
          onClick={async () => {
            try {
              await commonStore.setTheme(toggleTheme);
            } catch (err) {
              showError(message, err);
            }
          }}
        >
          {text}
        </NButton>
      );
    };

    // 用户信息
    const renderUserInfo = () => (
      <>
        <NButton bordered={false} onClick={() => goToProfile()}>
          <NIcon class={userInfoClass}>
            <User />
          </NIcon>
          {account.value}
        </NButton>
        <NButton bordered={false} onClick={logout}>
          退出登录
        </NButton>
      </>
    );
    // 登录注册等功能按钮
    const renderCtrls = () => (
      <>
        <NButton bordered={false} onClick={() => goToLogin()}>
          登录
        </NButton>
        <NButton bordered={false} onClick={() => goToRegister()}>
          注册
        </NButton>
      </>
    );
    return {
      renderToggleTheme,
      renderUserInfo,
      renderCtrls,
      anonymous,
      processing,
    };
  },
  render() {
    const { processing, anonymous } = this;
    return (
      <NLayoutHeader bordered class={headerClass}>
        <NText tag="div" class={logoClass}>
          <NH2>
            <a
              href="#"
              style="color: var(--primary-color);text-decoration: none;"
              onClick={(e) => {
                e.preventDefault();
                goToHome();
              }}
            >
              Forest
            </a>
          </NH2>
        </NText>
        <NSpace justify="end">
          {this.renderToggleTheme()}
          {processing && <span>正在加载中，请稍候...</span>}
          {!processing && !anonymous && this.renderUserInfo()}
          {!processing && anonymous && this.renderCtrls()}
        </NSpace>
      </NLayoutHeader>
    );
  },
});
