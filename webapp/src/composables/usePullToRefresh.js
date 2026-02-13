/**
 * 下拉刷新组合函数
 * 移动端手势触发刷新，阈值 60px，带触觉反馈（如浏览器支持）
 * 用法: const { pullDistance, isRefreshing, isPulling } = usePullToRefresh(onRefresh, containerRef)
 */
import { ref, onMounted, onUnmounted } from 'vue'

// 下拉刷新阈值（px）
const THRESHOLD = 60
// 最大下拉距离（px）
const MAX_PULL = 120
// 阻尼系数（越大越难拉）
const DAMPING = 2.5

export function usePullToRefresh(onRefresh, containerRef) {
  const pullDistance = ref(0)
  const isPulling = ref(false)
  const isRefreshing = ref(false)

  let startY = 0
  let currentY = 0

  // 检查容器是否在顶部（scrollTop === 0）
  const isAtTop = () => {
    if (containerRef?.value) {
      return containerRef.value.scrollTop <= 0
    }
    return window.scrollY <= 0
  }

  const onTouchStart = (e) => {
    if (isRefreshing.value) return
    if (!isAtTop()) return
    startY = e.touches[0].clientY
    isPulling.value = true
  }

  const onTouchMove = (e) => {
    if (!isPulling.value || isRefreshing.value) return

    currentY = e.touches[0].clientY
    const diff = currentY - startY

    // 仅处理下拉（diff > 0）
    if (diff <= 0) {
      pullDistance.value = 0
      return
    }

    // 如果容器可滚动且未在顶部，放弃下拉
    if (!isAtTop()) {
      isPulling.value = false
      pullDistance.value = 0
      return
    }

    // 阻尼效果
    pullDistance.value = Math.min(diff / DAMPING, MAX_PULL)

    // 达到阈值时触觉反馈
    if (pullDistance.value >= THRESHOLD && navigator.vibrate) {
      navigator.vibrate(10)
    }

    // 阻止页面滚动
    if (pullDistance.value > 0) {
      e.preventDefault()
    }
  }

  const onTouchEnd = async () => {
    if (!isPulling.value) return
    isPulling.value = false

    if (pullDistance.value >= THRESHOLD && !isRefreshing.value) {
      // 触发刷新
      isRefreshing.value = true
      pullDistance.value = THRESHOLD / DAMPING // 保持指示器在阈值位置

      try {
        await onRefresh()
      } finally {
        isRefreshing.value = false
        pullDistance.value = 0
      }
    } else {
      // 未达阈值，回弹
      pullDistance.value = 0
    }
  }

  onMounted(() => {
    const el = containerRef?.value || document
    el.addEventListener('touchstart', onTouchStart, { passive: true })
    el.addEventListener('touchmove', onTouchMove, { passive: false })
    el.addEventListener('touchend', onTouchEnd, { passive: true })
  })

  onUnmounted(() => {
    const el = containerRef?.value || document
    el.removeEventListener('touchstart', onTouchStart)
    el.removeEventListener('touchmove', onTouchMove)
    el.removeEventListener('touchend', onTouchEnd)
  })

  return {
    pullDistance,
    isPulling,
    isRefreshing
  }
}
