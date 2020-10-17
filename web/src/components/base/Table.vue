<script>
import { CONFIG_EDIT_MODE } from "@/constants/route";
export default {
  name: "BaseTable",
  computed: {
    currentPage() {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    },
    editMode() {
      return this.$route.query.mode === CONFIG_EDIT_MODE;
    }
  },
  methods: {
    handleCurrentChange(page) {
      this.query.offset = (page - 1) * this.query.limit;
      this.fetch();
    },
    handleSizeChange(pageSize) {
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
          mode: CONFIG_EDIT_MODE
        }
      });
    },
    modify(item) {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
          id: item.id
        }
      });
    },
    filter(params) {
      Object.assign(this.query, params);
      this.query.offset = 0;
      this.fetch();
    }
  },
  watch: {
    $route() {
      if (!this.editMode) {
        this.fetch();
      }
    }
  },
  beforeMount() {
    if (!this.editMode) {
      this.fetch();
    }
  }
};
</script>
