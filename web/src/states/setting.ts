import { reactive, readonly, DeepReadonly } from "vue";
import { userSettings } from "../store";

interface Setting {
  processing: boolean;
  // 主侧边栏是否隐藏
  mainNavShrinking: boolean;
}

const setting: Setting = reactive({
  processing: false,
  mainNavShrinking: false,
});

// settingLoad 加载配置
export async function settingLoad(): Promise<void> {
  if (setting.processing) {
    return;
  }
  try {
    setting.processing = true;
    const data = await userSettings.load();
    Object.assign(setting, data);
  } finally {
    setting.processing = false;
  }
}

// settingSave 保存配置
export async function settingSave(params: {
  mainNavShrinking?: boolean;
}): Promise<void> {
  const mainNavShrinking = params.mainNavShrinking || false;
  await userSettings.set("mainNavShrinking", mainNavShrinking);
  setting.mainNavShrinking = mainNavShrinking;
}

export default function useSettingState(): DeepReadonly<Setting> {
  return readonly(setting);
}
