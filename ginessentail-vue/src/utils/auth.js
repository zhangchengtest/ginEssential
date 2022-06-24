import Cookies from 'js-cookie'
import config from '../../package.json'
// 项目名 + 当前环境 + 项目版本 + 缓存key
const name = `${config.name}-${process.env.NODE_ENV ? 'loc' : process.env.VUE_APP_TITLE}-${config.version}-`
const TokenKey = `${name}accessToken`
const UserKey = `${name}userInfo`
const CurrentOrgKey = `${name}currentOrg`

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token, { expires: 5 })
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getUserInfo() {
  return localStorage.getItem(UserKey)
    ? JSON.parse(localStorage.getItem(UserKey))
    : {}
}

export function setUserInfo(userInfo) {
  return localStorage.setItem(UserKey, JSON.stringify(userInfo))
}

export function removeUserInfo() {
  return localStorage.removeItem(UserKey)
}

export function getCurrentOrg() {
  return localStorage.getItem(CurrentOrgKey)
    ? JSON.parse(localStorage.getItem(CurrentOrgKey))
    : {}
}

export function setCurrentOrg(currentOrg) {
  return localStorage.setItem(CurrentOrgKey, JSON.stringify(currentOrg))
}

export function removeCurrentOrg() {
  return localStorage.removeItem(CurrentOrgKey)
}

// Cookies存储
export function setCookies(key, value) {
  Cookies.set(name + key, value, { expires: 5 })
}
export function getCookies(key) {
  return Cookies.get(name + key)
}
export function clearOneCookies(key) {
  Cookies.remove(name + key)
}

// localStorage存储
export function setLocal(key, value) {
  localStorage.setItem(name + key, value)
}
export function getLocal(key) {
  return localStorage.getItem(name + key) ? JSON.parse(localStorage.getItem(name + key)) : {}
}
export function clearOneLocal(key) {
  localStorage.removeItem(name + key)
}
export function clearAllLocal() {
  localStorage.clear()
}

// setSession存储
export function setSession(key, value) {
  sessionStorage.setItem(name + key, value)
}
export function getSession(key) {
  return sessionStorage.getItem(name + key) ? JSON.parse(sessionStorage.getItem(name + key)) : {}
}
export function clearOneSession(key) {
  sessionStorage.removeItem(name + key)
}
export function clearAllSession() {
  sessionStorage.clear()
}
