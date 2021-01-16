import dayjs from "dayjs";

import { sha256 } from "./crypto";

const hash = "JT";
const oneDayMS = 24 * 3600 * 1000;

export function generatePassword(pass: string): string {
  return sha256(hash + sha256(pass + hash));
}

// formatDate 格式化日期
export function formatDate(str: string): string {
  if (!str) {
    return "--";
  }

  return dayjs(str).format("YYYY-MM-DD HH:mm:ss");
}

// isAllowedUser 判断是否允许该用户
export function isAllowedUser(
  allowList: string[],
  currentList: string[]
): boolean {
  if (!allowList || allowList.length === 0) {
    return true;
  }
  let allowed = false;
  allowList.forEach((item) => {
    currentList.forEach((current) => {
      if (current === item) {
        allowed = true;
      }
    });
  });
  return allowed;
}

// today 获取当天0点时间
export function today(): Date {
  return new Date(new Date(new Date().toLocaleDateString()).getTime());
}

// tomorrow 获取明天0点时间
export function tomorrow(): Date {
  return new Date(today().getTime() + oneDayMS);
}

// today 获取当天0点时间
export function yesterday(): Date {
  return new Date(today().getTime() - oneDayMS);
}
export function formatBegin(begin: Date): string {
  return begin.toISOString();
}
export function formatEnd(end: Date): string {
  return new Date(end.getTime() + 24 * 3600 * 1000 - 1).toISOString();
}

function isEqual(value: any, originalValue: any) {
  // 使用json stringify对比是否相同
  return JSON.stringify(value) == JSON.stringify(originalValue);
}

// diff  对比两个object的差异
export function diff(current: any, original: any) {
  const data: any = {};
  let modifiedCount = 0;
  Object.keys(current).forEach((key) => {
    const value = current[key];
    if (!isEqual(value, original[key])) {
      data[key] = value;
      modifiedCount++;
    }
  });
  return {
    modifiedCount,
    data,
  };
}

// validateForm validate form
export function validateForm(form: any) {
  return new Promise((resolve, reject) => {
    form.validate((valid: any, rules: any) => {
      if (valid) {
        return resolve();
      }
      const messagesArr: string[] = [];
      Object.keys(rules).forEach((key) => {
        const arr = rules[key];
        arr.forEach((item) => {
          messagesArr.push(item.message);
        });
      });
      return reject(new Error(messagesArr.join("，")));
    });
  });
}

// omitNil omit nil(undefined null)
export function omitNil(data: any) {
  const params: any = {};
  Object.keys(data).forEach((key) => {
    const value = data[key];
    if (value !== undefined && value !== null) {
      params[key] = value;
    }
  });
  return params;
}

// getFieldRules get field rules
export function getFieldRules(fields: any) {
  const rules: any = {};
  fields.forEach((field) => {
    if (field.rules) {
      rules[field.key] = field.rules;
    }
  });
  return rules;
}
