import request from '@/utils/request'

export function javatosql(data) {
  return request({
    url: '/api/javatosql',
    method: 'post',
    headers: { 'Content-Type': 'application/json;charset=UTF-8' },
    data: data
  })
}

