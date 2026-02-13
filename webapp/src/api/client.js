/**
 * API 客户端 - HTTP 请求封装
 * 提供统一的请求方法、错误处理和响应拦截
 */

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const API_TIMEOUT = 30000 // 30秒超时

/**
 * 自定义 API 错误类
 */
export class ApiError extends Error {
  constructor(code, message, data = null) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.data = data
  }
}

/**
 * 创建请求头
 * @returns {Headers}
 */
const createHeaders = () => {
  const headers = new Headers({
    'Content-Type': 'application/json'
  })

  // 从 localStorage 获取 token（如果有认证系统）
  const authData = localStorage.getItem('allfi-auth')
  if (authData) {
    try {
      const { token } = JSON.parse(authData)
      if (token) {
        headers.set('Authorization', `Bearer ${token}`)
      }
    } catch (e) {
      // 忽略解析错误
    }
  }

  return headers
}

/**
 * 发送 HTTP 请求
 * @param {string} endpoint - API 端点
 * @param {object} options - 请求选项
 * @returns {Promise<any>}
 */
async function request(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`

  const config = {
    headers: createHeaders(),
    ...options
  }

  // 设置超时
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), API_TIMEOUT)
  config.signal = controller.signal

  try {
    const response = await fetch(url, config)
    clearTimeout(timeoutId)

    // 解析响应
    const data = await response.json()

    // 处理 401 未授权（Token 过期或无效）
    if (response.status === 401) {
      localStorage.removeItem('allfi-auth')
      // 非认证接口收到 401 时跳转登录页
      if (!endpoint.startsWith('/auth/')) {
        window.location.href = '/login'
      }
      throw new ApiError(1005, data.message || '认证已过期，请重新登录')
    }

    // 检查业务错误码
    if (data.code !== 0 && data.code !== 200) {
      throw new ApiError(data.code, data.message || '请求失败', data.error)
    }

    return data.data
  } catch (error) {
    clearTimeout(timeoutId)

    // 处理网络错误
    if (error.name === 'AbortError') {
      throw new ApiError(0, '请求超时，请检查网络连接')
    }

    if (error instanceof ApiError) {
      throw error
    }

    // 处理其他错误
    throw new ApiError(0, error.message || '网络请求失败')
  }
}

/**
 * GET 请求
 * @param {string} endpoint - API 端点
 * @param {object} params - 查询参数
 * @returns {Promise<any>}
 */
export async function get(endpoint, params = {}) {
  const queryString = new URLSearchParams(params).toString()
  const url = queryString ? `${endpoint}?${queryString}` : endpoint
  return request(url, { method: 'GET' })
}

/**
 * POST 请求
 * @param {string} endpoint - API 端点
 * @param {object} data - 请求体
 * @returns {Promise<any>}
 */
export async function post(endpoint, data = {}) {
  return request(endpoint, {
    method: 'POST',
    body: JSON.stringify(data)
  })
}

/**
 * PUT 请求
 * @param {string} endpoint - API 端点
 * @param {object} data - 请求体
 * @returns {Promise<any>}
 */
export async function put(endpoint, data = {}) {
  return request(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

/**
 * DELETE 请求
 * @param {string} endpoint - API 端点
 * @returns {Promise<any>}
 */
export async function del(endpoint) {
  return request(endpoint, { method: 'DELETE' })
}

/**
 * 健康检查
 * @returns {Promise<boolean>}
 */
export async function healthCheck() {
  try {
    await get('/health')
    return true
  } catch {
    return false
  }
}

export default {
  get,
  post,
  put,
  del,
  healthCheck,
  ApiError
}
