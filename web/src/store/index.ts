import localforage from "localforage";

const store = localforage.createInstance({
  name: "forest",
});

// LimitStore limit store
class LimitStore {
  private key: string;
  private max: number;
  private data?: Record<string, unknown>[];
  constructor(key: string, max: number) {
    this.key = key;
    this.max = max;
  }
  // load 从本地存储中加载
  async load(): Promise<void> {
    const data = await store.getItem(this.key);
    if (!data) {
      return;
    }
    const arr = JSON.parse(data);
    this.data = arr;
  }
  // add 添加记录
  async add(item: Record<string, unknown>): Promise<void> {
    if (!this.data) {
      this.data = [];
    }
    this.data.push(item);
    if (this.data.length > this.max) {
      this.data.shift();
    }
    await store.setItem(this.key, JSON.stringify(this.data));
  }
  // clear 清除记录
  async clear() {
    const data = this.data;
    this.data = [];
    await store.removeItem(this.key);
    return data;
  }
  // 获取长度
  get size() {
    return this.data?.length || 0;
  }
}

class LocalStore {
  private key: string;
  constructor(key: string) {
    this.key = key;
  }
  // load 加载数据
  async load() {
    const data: string = await store.getItem(this.key);
    return JSON.parse(data || "{}");
  }
  // set 设置数据
  async set(key: string, value: string | number | boolean) {
    const data = await this.load();
    data[key] = value;
    await store.setItem(this.key, JSON.stringify(data));
    return data;
  }
}

// userActions 保存用户行为
export const userActions = new LimitStore("userActions", 50);
// userSettings 用户设置
export const userSettings = new LocalStore("userSettings");
// httpRequests http请求记录，测试环境使用
export const httpRequests = new LimitStore("httpRequests", 20);
