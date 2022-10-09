import Cookies from 'js-cookie'

const TokenKey = 'tokenPorta'
const RefreshTokenKey = 'refreshTokenPorta'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token, { expires: 5 })
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getRefreshToken() {
  return Cookies.get(RefreshTokenKey)
}

export function setRefreshToken(refreshToken) {
  return Cookies.set(RefreshTokenKey, refreshToken, { expires: 5 })
}

export function removeRefreshToken() {
  return Cookies.remove(RefreshTokenKey)
}

export function getlocalStorage(key) {
  return localStorage.getItem(key) ? JSON.parse(localStorage.getItem(key)) : {}
}

export function setlocalStorage(key, advList) {
  return localStorage.setItem(key, JSON.stringify(advList))
}

export function removelocalStorage(key) {
  return localStorage.removeItem(key)
}
