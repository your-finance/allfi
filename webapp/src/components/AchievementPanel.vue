<script setup>
/**
 * 成就面板
 * 按类别分组展示：里程碑 / 坚持 / 投资
 * 进度概览：已解锁 X/Y 个
 */
import { onMounted } from 'vue'
import { PhTrophy, PhX } from '@phosphor-icons/vue'
import AchievementBadge from './AchievementBadge.vue'
import AchievementUnlock from './AchievementUnlock.vue'
import { useAchievementStore } from '../stores/achievementStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
})

const emit = defineEmits(['close'])

const achStore = useAchievementStore()
const { t } = useI18n()

// 类别配置
const categories = [
  { key: 'milestone', labelKey: 'achievement.categoryMilestone', icon: '🏆' },
  { key: 'persistence', labelKey: 'achievement.categoryPersistence', icon: '🔥' },
  { key: 'investment', labelKey: 'achievement.categoryInvestment', icon: '📈' },
]

onMounted(() => {
  if (achStore.achievements.length === 0) {
    achStore.fetchAchievements()
  }
})
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="panel-overlay" @click.self="emit('close')">
      <div class="achievement-panel">
        <!-- 头部 -->
        <div class="panel-header">
          <div class="header-left">
            <PhTrophy :size="20" weight="fill" class="trophy-icon" />
            <h2>{{ t('achievement.title') }}</h2>
          </div>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="18" />
          </button>
        </div>

        <!-- 进度概览 -->
        <div class="progress-overview">
          <div class="progress-text">
            <span class="progress-count font-mono">{{ achStore.unlockedCount }}</span>
            <span class="progress-total"> / {{ achStore.totalCount }}</span>
            <span class="progress-label">{{ t('achievement.unlocked') }}</span>
          </div>
          <div class="progress-bar-track">
            <div
              class="progress-bar"
              :style="{ width: achStore.totalCount > 0 ? (achStore.unlockedCount / achStore.totalCount * 100) + '%' : '0%' }"
            />
          </div>
        </div>

        <!-- 加载中 -->
        <div v-if="achStore.isLoading" class="loading-state">
          {{ t('common.loading') }}
        </div>

        <!-- 成就列表（按类别分组） -->
        <div v-else class="achievement-list">
          <div v-for="cat in categories" :key="cat.key" class="category-group">
            <div class="category-header">
              <span class="category-icon">{{ cat.icon }}</span>
              <span class="category-name">{{ t(cat.labelKey) }}</span>
              <span class="category-count font-mono">
                {{ achStore.byCategory[cat.key].filter(a => a.unlocked).length }}/{{ achStore.byCategory[cat.key].length }}
              </span>
            </div>
            <div class="badge-grid">
              <AchievementBadge
                v-for="ach in achStore.byCategory[cat.key]"
                :key="ach.id"
                :achievement="ach"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 解锁动画 -->
      <AchievementUnlock
        v-if="achStore.pendingUnlock"
        :achievement="achStore.pendingUnlock"
        @close="achStore.dismissUnlock()"
      />
    </div>
  </Transition>
</template>

<style scoped>
.panel-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.achievement-panel {
  width: 100%;
  max-width: 560px;
  max-height: 85vh;
  overflow-y: auto;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--gap-xl);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-lg);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.trophy-icon {
  color: var(--color-warning);
}

.header-left h2 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--gap-xs);
}

.close-btn:hover {
  color: var(--color-text-primary);
}

/* 进度概览 */
.progress-overview {
  margin-bottom: var(--gap-lg);
}

.progress-text {
  display: flex;
  align-items: baseline;
  gap: 2px;
  margin-bottom: var(--gap-xs);
}

.progress-count {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-accent-primary);
}

.progress-total {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.progress-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-left: var(--gap-xs);
}

.progress-bar-track {
  height: 6px;
  background: var(--color-bg-tertiary);
  border-radius: 3px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--color-accent-primary);
  border-radius: 3px;
  transition: width 0.6s ease;
}

.loading-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
}

/* 成就列表 */
.achievement-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

.category-header {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-sm);
}

.category-icon {
  font-size: 1rem;
}

.category-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  flex: 1;
}

.category-count {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.badge-grid {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

/* 过渡 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
