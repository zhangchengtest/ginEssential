let baseUrl = ''
let jumpUrl = ''
if (process.env.NODE_ENV === 'development') {
  // baseUrl = 'http://192.168.78.198:7772/' // 高一金:7772
  baseUrl = 'http://127.0.0.1:8080/' // 测试环境
  jumpUrl = 'https://aq-test.cunwedu.com.cn/#/passwordUpdate' // 开发环境
} else if (process.env.VUE_APP_TITLE === 'test') {
  baseUrl = 'https://test.cunwedu.com.cn/aq/' // 测试环境
  jumpUrl = 'https://aq-test.cunwedu.com.cn/#/passwordUpdate' // 测试环境
} else if (process.env.VUE_APP_TITLE === 'dev') {
  baseUrl = 'http://127.0.0.1:8080/'// 开发环境
  // baseUrl = 'https://test.cunwedu.com.cn/aq/'
} else if (process.env.VUE_APP_TITLE === 'pre') {
  baseUrl = 'https://test.cunwedu.com.cn/aq/' // 预发布环境
} else {
  baseUrl = 'http://172.31.21.3:8080/' // 生产环境
}
const ENVIRONMENT = {
  BASE_API: baseUrl, // api访问地址
  JUMP_API: jumpUrl // api访问地址
}

export default ENVIRONMENT