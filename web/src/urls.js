const prefix = "";

const userRoute = `${prefix}/users`;
export const USERS = `${userRoute}/v1`;
export const USERS_UPDATE = `${userRoute}/v1/update/:id`;
export const USERS_ME = `${userRoute}/v1/me`;
export const USERS_LOGIN = `${userRoute}/v1/me/login`;
export const USERS_LOGOUT = `${userRoute}/v1/me/logout`;
export const USERS_LOGIN_RECORDS = `${userRoute}/v1/login-records`;

const configurationRoute = `${prefix}/configurations`;
export const CONFIGURATIONS_ADD = `${configurationRoute}/v1`;
export const CONFIGURATIONS_UPDATE = `${configurationRoute}/v1/:id`;
export const CONFIGURATIONS_DELETE = `${configurationRoute}/v1/:id`;
export const CONFIGURATIONS_LIST = `${configurationRoute}/v1`;
export const CONFIGURATIONS_LIST_AVAILABLE = `${configurationRoute}/v1/available`;
export const CONFIGURATIONS_LIST_UNAVAILABLE = `${configurationRoute}/v1/unavailable`;

const commonRoute = `${prefix}/commons`;
export const ROUTERS = `${commonRoute}/routers`;
export const RANDOM_KEYS = `${commonRoute}/random-keys`;
export const CAPTCHA = `${commonRoute}/captcha`;
