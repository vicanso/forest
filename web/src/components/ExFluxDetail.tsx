import { css } from "@linaria/core";
import { InfoCircle } from "@vicons/fa";
import { NIcon, NPopover } from "naive-ui";
import { defineComponent, PropType } from "vue";

const infoListCalss = css`
  margin: 0;
  padding: 0 20px;
  max-width: 400px;
  word-wrap: break-word;
  word-break: break-all;
  white-space: normal;
  list-style-position: insied;
`;

export default defineComponent({
  name: "ExFluxDetail",
  props: {
    data: {
      type: Object as PropType<Record<string, unknown>>,
      required: true,
    },
  },
  render() {
    const { data } = this.$props;
    const slots = {
      trigger: () => (
        <NIcon>
          <InfoCircle />
        </NIcon>
      ),
    };
    const tags: Record<string, string> = {};
    Object.keys(data).forEach((key) => {
      const v = data[key];
      if (!v) {
        return;
      }
      tags[key] = v as string;
    });

    const ignoreKeys = [
      "_measurement",
      "_start",
      "_stop",
      "_time",
      "result",
      "table",
    ];
    const values: Record<string, unknown>[] = [];
    Object.keys(data).forEach((key) => {
      if (ignoreKeys.includes(key)) {
        return;
      }
      values.push({
        name: key,
        value: data[key],
      });
    });
    const arr = values.map((item) => {
      return (
        <li>
          <span class="mright5">{item.name}:</span> {String(item.value)}
        </li>
      );
    });
    return (
      <NPopover
        trigger="hover"
        placement="top-end"
        delay={500}
        duration={1000}
        v-slots={slots}
      >
        <ul class={infoListCalss}>{arr}</ul>
      </NPopover>
    );
  },
});
