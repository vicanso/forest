import { CONFIG_EDIT_MODE } from "../constants/common";
import { diff } from "../helpers/util";

export default {
  computed: {
    currentPage(): number {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    },
    editMode(): boolean {
      return this.$route.query.mode === CONFIG_EDIT_MODE;
    },
  },
  methods: {
    handleCurrentChange(page: number) {
      this.query.offset = (page - 1) * this.query.limit;
      this.fetch();
    },
    handleSizeChange(pageSize: number) {
      this.query.limit = pageSize;
      this.query.offset = 0;
      this.fetch();
    },
    handleSortChange({ prop, order }) {
      let key = prop.replace("Desc", "");
      if (order === "descending") {
        key = `-${key}`;
      }
      this.query.order = key;
      this.query.offset = 0;
      this.fetch();
    },
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
        },
      });
    },
    modify(item: any) {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
          id: item.id,
        },
      });
    },
    filter(params: any) {
      Object.assign(this.query, params);
      this.query.offset = 0;
      this.fetch();
    },
  },
  watch: {
    "$route.query"(query: any, prevQuery: any) {
      // 如果路由已更换，则直接返回
      if (this.$route.name !== this._currentRoute) {
        return;
      }
      if (!diff(query, prevQuery).modifiedCount) {
        return;
      }

      if (!this.editMode) {
        this.fetch();
      }
    },
  },
  beforeMount() {
    this._currentRoute = this.$route.name;
    if (!this.editMode && !this.disableBeforeMountFetch) {
      this.fetch();
    }
  },
};
