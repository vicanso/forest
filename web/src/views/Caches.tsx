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
  NIcon,
} from "naive-ui";
import { Trash } from "@vicons/fa";
import { showError, showWarning, toast } from "../helpers/util";
import { useAdminStore } from "../stores/admin";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";

const cacheTableClass = css`
  top: 10px;
`;

export default defineComponent({
  name: "CachesView",
  setup() {
    const message = useMessage();
    const keyword = ref("");
    const adminStore = useAdminStore();
    const { cacheKeys, fetchingCacheKeys } = storeToRefs(adminStore);

    const fetch = async () => {
      if (!keyword.value) {
        showWarning(message, "请输入要查询的key");
        return;
      }
      try {
        await adminStore.listCacheKeys({
          keyword: keyword.value,
        });
      } catch (err) {
        showError(message, err);
      }
    };

    const del = async (key: string) => {
      if (!key) {
        showWarning(message, "请输入要删除的key");
        return;
      }
      try {
        await adminStore.removeCache(key);
        toast(message, "已成功清除数据");
      } catch (err) {
        showError(message, err);
      }
    };

    return {
      fetchingCacheKeys,
      cacheKeys,
      keyword,
      fetch,
      del,
    };
  },
  render() {
    const size = "large";
    const { fetch, fetchingCacheKeys, cacheKeys, del } = this;
    const columns: TableColumn[] = [
      {
        title: "KEY",
        key: "key",
      },
      {
        title: "操作",
        key: "actions",
        render: (row) => {
          return (
            <NButton
              bordered={false}
              onClick={() => {
                del(row.key as string);
              }}
            >
              <NIcon>
                <Trash />
              </NIcon>
              删除
            </NButton>
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
            <NGridItem span={18}>
              <NInput
                placeholder="请输入缓存的key"
                size={size}
                clearable
                onUpdateValue={(value) => {
                  this.keyword = value;
                }}
              />
            </NGridItem>
            <NGridItem span={6}>
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
