import axios from "axios";

import { RANDOM_KEYS } from "../urls";

export async function getRandomKeys(params) {
  const { data } = await axios.get(RANDOM_KEYS, {
    params
  });
  return data;
}
