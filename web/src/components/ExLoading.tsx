import { NSpin } from "naive-ui";
import { defineComponent } from "vue";

const spinStyle = {
  float: "left",
  marginRight: "10px",
};
const loadingStyle = {
  margin: "auto",
  width: "200px",
};

export default defineComponent({
  name: "ExLoading",
  props: {
    style: {
      type: Object,
      default: () => {
        return {
          marginTop: "60px",
        };
      },
    },
  },
  render() {
    const style = Object.assign({}, loadingStyle, this.$props.style);
    return (
      <div style={style}>
        <NSpin size="small" style={spinStyle} />
        正在加载中，请稍候...
      </div>
    );
  },
});
