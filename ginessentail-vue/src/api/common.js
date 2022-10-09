import request from '@/utils/request'
/**
 *
 *  获取阿里云临时目录
 */
 export function authTemp() {
    return request({
      url: `/api/oss/authTemp`,
      method: 'get'
    })
  }
  /**
   *
   *  转换路径
   */
  export function tempFormal(data) {
    return request({
      url: `/api/oss/temp2formal`,
      timeout: 0,
      method: 'post',
      data
    })
  }
