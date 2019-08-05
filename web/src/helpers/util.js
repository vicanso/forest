import { sha256 } from "./crypto";

const hash = "JT";

export function setBeginOfDay(date) {
  return date
    .clone()
    .hours(0)
    .minutes(0)
    .seconds(0)
    .milliseconds(0);
}

export function setEndOfDay(date) {
  return date
    .clone()
    .hours(23)
    .minutes(59)
    .seconds(59)
    .milliseconds(999);
}

export function generatePassword(pass) {
  return sha256(pass + hash);
}
