import { sha256 } from "@/helpers/crypto";

const hash = "JT";
const oneDayMS = 24 * 3600 * 1000;

export function generatePassword(pass) {
  return sha256(pass + hash);
}

// formatDate 格式化日期
export function formatDate(str) {
  if (!str) {
    return "--";
  }
  const d = new Date(str);
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString()}`;
}

// formatDuration 格式化duration
export function formatDuration(d) {
  if (!d) {
    return "--";
  }
  if (d > 1000) {
    const v = d / 100;
    let fix = 1;
    if (Number.parseInt(v) % 10 === 0) {
      fix = 0;
    }
    return `${(d / 1000).toFixed(fix)}秒`;
  }
  return `${d}毫秒`;
}

// delay 延时promise
export function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function isEqual(value, originalValue) {
  // 使用json stringify对比是否相同
  return JSON.stringify(value) == JSON.stringify(originalValue);
}

// diff  对比两个object的差异
export function diff(current, original) {
  const data = {};
  let modifiedCount = 0;
  Object.keys(current).forEach(key => {
    const value = current[key];
    if (!isEqual(value, original[key])) {
      data[key] = value;
      modifiedCount++;
    }
  });
  return {
    modifiedCount,
    data
  };
}

// isAllowedUser 判断是否允许该用户
export function isAllowedUser(allowedRoles, userRoles) {
  let allowed = false;
  allowedRoles.forEach(item => {
    userRoles.forEach(userRole => {
      if (userRole === item) {
        allowed = true;
      }
    });
  });
  return allowed;
}

// queryOmitEmpty 删除query中的空值
export function queryOmitEmpty(query) {
  const params = {};
  Object.keys(query).forEach(key => {
    if (query[key]) {
      params[key] = query[key];
    }
  });
  return params;
}

// omitNil omit nil(undefined null)
export function omitNil(data) {
  const params = {};
  Object.keys(data).forEach(key => {
    const value = data[key];
    if (value !== undefined && value !== null) {
      params[key] = value;
    }
  });
  return params;
}

// today 获取当天0点时间
export function today() {
  return new Date(new Date(new Date().toLocaleDateString()).getTime());
}

// tomorrow 获取明天0点时间
export function tomorrow() {
  return new Date(today().getTime() + oneDayMS);
}

// today 获取当天0点时间
export function yesterday() {
  return new Date(today().getTime() - oneDayMS);
}
export function formatBegin(begin) {
  return begin.toISOString();
}
export function formatEnd(end) {
  return new Date(end.getTime() + 24 * 3600 * 1000 - 1).toISOString();
}

// contains 判断是否包括该元素
export function contains(arr, value) {
  let exists = false;
  arr.forEach(item => {
    if (item === value) {
      exists = true;
    }
  });
  return exists;
}

// addNoCacheQueryParam 添加不缓存query参数
export function addNoCacheQueryParam(params = {}) {
  params["cacheControl"] = "no-cache";
  return params;
}

// validateForm validate form
export function validateForm(form) {
  return new Promise((resolve, reject) => {
    form.validate((valid, rules) => {
      if (valid) {
        return resolve();
      }
      const messagesArr = [];
      Object.keys(rules).forEach(key => {
        const arr = rules[key];
        arr.forEach(item => {
          messagesArr.push(item.message);
        });
      });
      return reject(new Error(messagesArr.join("，")));
    });
  });
}

// toUploadFiles to upload files
export function toUploadFiles(files) {
  if (!files) {
    return [];
  }
  if (Array.isArray(files)) {
    return files.map(file => {
      return {
        url: file
      };
    });
  }
  return [
    {
      url: files
    }
  ];
}

// findByID find from array by id
export function findByID(arr, id) {
  let found = null;
  if (!arr || !arr.length) {
    return found;
  }
  arr.forEach(item => {
    if (item.id === id) {
      found = item;
    }
  });
  return found;
}

// getFieldRules get field rules
export function getFieldRules(fields) {
  const rules = {};
  fields.forEach(field => {
    if (field.rules) {
      rules[field.key] = field.rules;
    }
  });
  return rules;
}
