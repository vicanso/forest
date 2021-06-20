import localforage from "localforage";

const store = localforage.createInstance({
  name: "forest",
});

class LocalStorage {
  private key: string;
  constructor(key: string) {
    this.key = key;
  }
  // load 加载数据
  async load(): Promise<Record<string, unknown>> {
    const data = await store.getItem(this.key);
    if (!data) {
      return {};
    }
    const str = data as string;
    return JSON.parse(str || "{}");
  }
  // set 设置数据
  async set(key: string, value: string | number | boolean) {
    const data = await this.load();
    data[key] = value;
    await store.setItem(this.key, JSON.stringify(data));
    return data;
  }
}

export const settingStorage = new LocalStorage("settings");
