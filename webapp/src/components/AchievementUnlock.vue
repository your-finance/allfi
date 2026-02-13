<script setup>
/**
 * 成就解锁动画
 * 简单的缩放 + 渐入动画
 * 成就名称 + 描述 + 解锁时间
 */
import { ref, onMounted } from 'vue'
import { PhTrophy } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  achievement: { type: Object, required: true },
})

const emit = defineEmits(['close'])

const { t } = useI18n()
const isVisible = ref(false)

// 类别颜色
const categoryColors = {
  milestone: '#F59E0B',
  persistence: '#3B82F6',
  investment: '#10B981',
}

onMounted(() => {
  // 延迟触发动画
  requestAnimationFrame(() => {
    isVisible.value = true
  })
})

const handleClose = () => {
  isVisible.value = false
  setTimeout(() => emit('close'), 200)
}
</script>

<template>
  <Transition name="unlock">
    <div v-if="isVisible" class="unlock-overlay" @click.self="handleClose">
      <div class="unlock-card">
        <div class="unlock-icon" :style="{ background: categoryColors[achievement.category] + '18', color: categoryColors[achievement.category] }">
          <PhTrophy :size="32" weight="fill" />
        </div>
        <div class="unlock-label">{{ t('achievement.newUnlock') }}</div>
        <div class="unlock-name">{{ achievement.name }}</div>
        <div class="unlock-desc">{{ achievement.description }}</div>
        <div class="unlock-date" v-if="achievement.unlockedAt">
          {{ achievement.unlockedAt }}
        </div>
        <button class="unlock-btn" @click="handleClose">
          {{ t('common.confirm') }}
        </button>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.unlock-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
}

.unlock-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-xl) var(--gap-2xl);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  text-align: center;
  min-width: 280px;
  animation: scaleIn 0.3s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.85);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.unlock-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--gap-xs);
}

.unlock-label {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-accent-primary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.unlock-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.unlock-desc {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  max-width: 240px;
}

.unlock-date {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.unlock-btn {
  margin-top: var(--gap-sm);
  padding: 6px 24px;
  font-size: 0.8125rem;
  font-weight: 500;
  background: var(--color-accent-primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: opacity var(--transition-fast);
}

.unlock-btn:hover {
  opacity: 0.9;
}

/* 过渡 */
.unlock-enter-active {
  transition: opacity 0.2s ease;
}

.unlock-leave-active {
  transition: opacity 0.2s ease;
}

.unlock-enter-from,
.unlock-leave-to {
  opacity: 0;
}
</style>
