import { defineComponent, ref } from "vue";
import { css } from "@linaria/core";
import { storeToRefs } from "pinia";
import {
  NCard,
  NGrid,
  NGridItem,
  NInput,
  NButton,
  useMessage,
  NSpin,
  NDataTable,
  NPopconfirm,
  NPopover,
  NInputNumber,
} from "naive-ui";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { showError, showWarning, toast, formatJSON } from "../helpers/util";
import { useAdminStore } from "../stores/admin";

const cacheTableClass = css`
  top: 10px;
`;

export default defineComponent({
  name: "CachesView",
  setup() {
    const message = useMessage();
    const keyword = ref("");
    const limit = ref(0);
    const offset = ref(0);
    const adminStore = useAdminStore();
    const { cacheKeys, fetchingCacheKeys } = storeToRefs(adminStore);
    const cacheData = ref("");

    const fetch = async () => {
      if (!keyword.value) {
        showWarning(message, "请输入要查询的key");
        return;
      }
      try {
        await adminStore.listCacheKeys({
          keyword: keyword.value,
          limit: limit.value,
          offset: offset.value,
        });
      } catch (err) {
        showError(message, err);
      }
    };

    const del = async (key: string) => {
      try {
        await adminStore.removeCache(key);
        toast(message, "已成功清除数据");
      } catch (err) {
        showError(message, err);
      }
    };
    const getCache = async (key: string) => {
      try {
        cacheData.value = "正在获取中，请稍候...";
        const data = await adminStore.getCache(key);
        cacheData.value = formatJSON(data);
      } catch (err) {
        showError(message, err);
      }
    };

    return {
      fetchingCacheKeys,
      cacheKeys,
      keyword,
      limit,
      offset,
      fetch,
      del,
      getCache,
      cacheData,
    };
  },
  render() {
    const size = "large";
    const { fetch, fetchingCacheKeys, cacheKeys, del, getCache, cacheData } =
      this;
    const columns: TableColumn[] = [
      {
        title: "KEY",
        key: "key",
      },
      {
        title: "操作",
        key: "actions",
        width: 150,
        align: "center",
        render: (row) => {
          const delSlots = {
            trigger: () => <NButton bordered={false}>删除</NButton>,
          };
          const viewSlots = {
            trigger: () => (
              <NButton
                bordered={false}
                onClick={() => {
                  getCache(row.key as string);
                }}
              >
                查看
              </NButton>
            ),
          };
          return (
            <div>
              <NPopover v-slots={viewSlots} trigger="click" placement="left">
                <pre>{cacheData}</pre>
              </NPopover>
              <NPopconfirm
                v-slots={delSlots}
                onPositiveClick={() => {
                  del(row.key as string);
                }}
              >
                <span>确认要删除此数据吗？</span>
              </NPopconfirm>
            </div>
          );
        },
      },
    ];
    const data = cacheKeys.map((item) => {
      return {
        key: item,
      };
    });
    return (
      <NSpin show={fetchingCacheKeys}>
        <NCard title="缓存查询与清除">
          <NGrid xGap={24}>
            <NGridItem span={12}>
              <NInput
                placeholder="请输入缓存的key"
                size={size}
                clearable
                onUpdateValue={(value) => {
                  this.keyword = value;
                }}
              />
            </NGridItem>
            <NGridItem span={4}>
              <NInputNumber
                placeholder="请输入Scan的数量"
                size={size}
                clearable
                defaultValue={1000}
                onUpdateValue={(value) => {
                  this.limit = Number(value);
                }}
              />
            </NGridItem>
            <NGridItem span={4}>
              <NInputNumber
                placeholder="请输入Cursor的数量"
                size={size}
                clearable
                onUpdateValue={(value) => {
                  this.offset = Number(value);
                }}
              />
            </NGridItem>
            <NGridItem span={4}>
              <NButton class="widthFull" size={size} onClick={() => fetch()}>
                查询
              </NButton>
            </NGridItem>
          </NGrid>
          <NDataTable class={cacheTableClass} columns={columns} data={data} />
        </NCard>
      </NSpin>
    );
  },
});
