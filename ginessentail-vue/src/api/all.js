import request from '@/utils/request'

export function javatosql(data) {
  return request({
    url: '/api/javatosql',
    method: 'post',
    headers: { 'Content-Type': 'application/json;charset=UTF-8' },
    data: data
  })
}

export function compareFile(data) {
  return request({
    url: '/api/compareFile',
    method: 'post',
    data: data
  })
}

