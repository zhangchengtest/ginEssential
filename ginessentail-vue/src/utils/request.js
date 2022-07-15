import axios from 'axios'
import router from '../router'
import { Message, MessageBox, Loading } from 'element-ui'
import store from '../store'
import { getToken } from '@/utils/auth'
import ENVIRONMENT from '@/utils/environment'
import qs from 'qs'
let loadingInstance
let loadingCount = 0
// 创建axios实例
const service = axios.create({
  baseURL: ENVIRONMENT.BASE_API, // api 的 base_url
  timeout: 1000 * 10 // 请求超时时间
  // withCredentials: true
})
// axios.defaults.withCredentials = true//带上cookie
// request拦截器
service.interceptors.request.use(
  (config) => {
    if (loadingCount === 0) {
      loadingInstance = Loading.service({
        lock: true,
        customClass: 'z-index999',
        text: '数据加载中，请稍后...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
    }
    loadingCount++
    // 设置浏览器不缓存
    config.headers['Cache-Control'] = 'no-cache'
    config.headers['Pragma'] = 'no-cache'
    // 设置请求数据格式
    config.headers['Content-Type'] =
      config.headers['Content-Type'] || 'application/json;charset=UTF-8'
    if (store.getters.token) {
      const Token = getToken()
      config.headers['Authorization'] = `Bearer ${Token}`
    }
    if (config.method === 'post') {
      // 如果为post请求，要求Content-Type为json时，数据需要转换成JSON字符串
      if (config.headers['Content-Type'] === 'application/json;charset=UTF-8') {
        config.data = JSON.stringify(config.data)
      } else {
        // 如果是FormData参数据不需要qs序列化处理，否则默认Content-Type数据格式都需要qs序列化处理
        if (!(config.data instanceof FormData)) {
          config.data = qs.stringify(config.data)
        }
      }
    }
    return config
  },
  (error) => {
    // Do something with request error
    console.log(error) // for debug
    Promise.reject(error)
  }
)

// response 拦截器
service.interceptors.response.use(
  (response) => {
    loadingCount--
    if (loadingInstance && loadingCount === 0) {
      loadingInstance.close()
    }
    const res = response.data
    console.log(res.code)
    if (res.code !== 200) {
      if (res.code === '401') {
        if (
          router.currentRoute.path !== '/' &&
          router.currentRoute.path !== '/login'
        ) {
          MessageBox.confirm(
            '你已被登出，可以取消继续留在该页面，或者重新登录',
            '确定登出',
            {
              confirmButtonText: '重新登录',
              cancelButtonText: '取消',
              type: 'warning'
            }
          ).then(() => {
            store.dispatch('FedLogOut').then(() => {
              location.reload() // 为了重新实例化vue-router对象 避免bug
            })
          })
        }
      } else {
        if (
          response.config.params &&
          response.config.params.export &&
          res &&
          !res.failed
        ) {
          // 如果是导出，response.data返回为乱码则不判断
          return response.data
        }
        // 如果需要自定义对错误做处理
        if (
          response.config.data &&
          response.config.data.indexOf('custom') > -1
        ) {
          return response.data
        }
        // 统一处理错误信息
        Message({
          message: res.message,
          type: 'error',
          duration: 3 * 1000
        })
      }
      return Promise.reject(response)
    } else {
      return response
    }
  },
  (error) => {
    console.log('err' + error) // for debug
    loadingCount--
    if (loadingInstance && loadingCount === 0) {
      loadingInstance.close()
    }
    if (error.message.includes('timeout')) {
      // 判断请求异常信息中是否含有超时timeout字符串
      Message({
        message: '请求超时',
        type: 'error'
      })
    } else if (error.response) {
      switch (error.response.status) {
        case 200:
          break
        case 401:
          break
        case 403:
        case 404:
        case 500:
          break
        default:
          break
      }
    }
    return Promise.reject(error.response)
  }
)

export default service
