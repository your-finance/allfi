/**
 * 系统管理 API 服务
 * 版本查询、更新检测、一键更新、回滚
 */
import { get, post } from './client.js'

// Mock 模式判断
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'

// 模拟延迟
const delay = (ms = 500) => new Promise(r => setTimeout(r, ms))

export const systemService = {
  /** 获取当前版本信息 */
  async getVersion() {
    if (USE_MOCK) {
      await delay(300)
      return {
        version: __APP_VERSION__,
        build_time: new Date().toISOString(),
        git_commit: 'mock123',
        run_mode: 'host',
        go_version: 'go1.24',
      }
    }
    return get('/system/version')
  },

  /** 检查新版本 */
  async checkUpdate() {
    if (USE_MOCK) {
      await delay(800)
      return {
        has_update: true,
        current_version: __APP_VERSION__,
        latest_version: '999.0.0',
        release_notes: '这是 Mock 模式的模拟更新说明',
        release_url: 'https://github.com/your-finance/allfi/releases',
        published_at: new Date().toISOString(),
      }
    }
    return get('/system/update/check')
  },

  /** 执行更新 */
  async applyUpdate(targetVersion) {
    if (USE_MOCK) {
      await delay(500)
      return { status: 'started', message: 'Mock: 更新已启动' }
    }
    return post('/system/update/apply', { target_version: targetVersion })
  },

  /** 版本回滚 */
  async rollback(targetVersion) {
    if (USE_MOCK) {
      await delay(500)
      return { status: 'started', message: 'Mock: 回滚已启动' }
    }
    return post('/system/update/rollback', { target_version: targetVersion })
  },

  /** 获取更新状态 */
  async getUpdateStatus() {
    if (USE_MOCK) {
      await delay(200)
      return { state: 'idle', step: 0, total: 3, step_name: '', message: '' }
    }
    return get('/system/update/status')
  },

  /** 获取更新历史 */
  async getUpdateHistory() {
    if (USE_MOCK) {
      await delay(300)
      return { records: [] }
    }
    return get('/system/update/history')
  },
}
