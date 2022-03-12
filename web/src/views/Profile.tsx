import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NInput,
  useMessage,
} from "naive-ui";
import { defineComponent, ref } from "vue";
import { storeToRefs } from "pinia";

import ExLoading from "../components/ExLoading";
import { showError, showWarning } from "../helpers/util";
import { useUserStore } from "../stores/user";

export default defineComponent({
  name: "ProfileView",
  setup() {
    const message = useMessage();
    const processing = ref(false);
    const userStore = useUserStore();
    const { account, email, name } = storeToRefs(userStore);
    const fetch = async () => {
      processing.value = true;
      try {
        await userStore.detail();
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };
    const form = ref({
      name: "",
      email: "",
    });
    const update = async () => {
      if (processing.value) {
        return;
      }
      const { name, email } = form.value;
      if (!name && !email) {
        showWarning(message, "信息无修改无需要更新");
        return;
      }
      processing.value = true;
      try {
        await userStore.update({
          name,
          email,
        });
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };
    fetch();
    return {
      processing,
      name,
      account,
      email,
      form,
      update,
    };
  },
  render() {
    const { processing, name, email, account, form, update } = this;
    if (processing && !account) {
      return <ExLoading />;
    }
    let text = "更新";
    if (processing) {
      text = "更新中..";
    }
    const size = "large";
    return (
      <NCard title={"个人信息"}>
        <NForm>
          <NGrid xGap={24}>
            <NGridItem span={12}>
              <NFormItem label="用户：">
                <NInput
                  placeholder="请输入用户名称"
                  defaultValue={name}
                  clearable
                  size={size}
                  onUpdateValue={(value) => {
                    form.name = value;
                  }}
                />
              </NFormItem>
            </NGridItem>
            <NGridItem span={12}>
              <NFormItem label="邮箱地址：">
                <NInput
                  placeholder="请输入邮箱地址"
                  defaultValue={email}
                  clearable
                  size={size}
                  onUpdateValue={(value) => {
                    form.email = value;
                  }}
                />
              </NFormItem>
            </NGridItem>
            <NGridItem span={24}>
              <NButton class="widthFull" size={size} onClick={update}>
                {text}
              </NButton>
            </NGridItem>
          </NGrid>
        </NForm>
      </NCard>
    );
  },
});
