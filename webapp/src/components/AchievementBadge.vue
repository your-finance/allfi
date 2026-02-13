<script setup>
/**
 * 成就徽章组件
 * 已解锁：彩色图标 + 名称
 * 未解锁：灰色图标 + 描述
 * hover 显示解锁条件
 */
import {
  PhStar,
  PhCrown,
  PhRocketLaunch,
  PhShieldCheck,
  PhBird,
  PhCalendarBlank,
  PhTimer,
  PhTrendUp,
  PhDatabase,
  PhCoins,
  PhMedal,
  PhDiamond,
  PhTarget
} from '@phosphor-icons/vue'

const props = defineProps({
  achievement: { type: Object, required: true },
})

// 图标映射
const iconMap = {
  seed: PhStar,
  collection: PhDatabase,
  btc: PhCoins,
  money: PhDiamond,
  rocket: PhRocketLaunch,
  bird: PhBird,
  data: PhDatabase,
  calendar: PhCalendarBlank,
  timer: PhTimer,
  trend: PhTrendUp,
  crown: PhCrown,
  shield: PhShieldCheck,
  scatter: PhTarget,
}

// 类别颜色
const categoryColors = {
  milestone: '#F59E0B',
  persistence: '#3B82F6',
  investment: '#10B981',
}
</script>

<template>
  <div
    class="badge"
    :class="{ unlocked: achievement.unlocked }"
    :title="achievement.unlocked ? achievement.name : achievement.condition"
  >
    <div
      class="badge-icon"
      :style="achievement.unlocked ? { background: categoryColors[achievement.category] + '18', color: categoryColors[achievement.category] } : {}"
    >
      <component
        :is="iconMap[achievement.icon] || PhMedal"
        :size="20"
        :weight="achievement.unlocked ? 'fill' : 'regular'"
      />
    </div>
    <div class="badge-info">
      <span class="badge-name">{{ achievement.unlocked ? achievement.name : '???' }}</span>
      <span class="badge-desc">{{ achievement.description }}</span>
    </div>
    <span v-if="achievement.unlocked && achievement.unlockedAt" class="badge-date">
      {{ achievement.unlockedAt }}
    </span>
  </div>
</template>

<style scoped>
.badge {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  cursor: default;
}

.badge:hover {
  background: var(--color-bg-tertiary);
}

.badge:not(.unlocked) {
  opacity: 0.45;
}

.badge-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-bg-tertiary);
  color: var(--color-text-muted);
}

.badge-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.badge-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.badge:not(.unlocked) .badge-name {
  color: var(--color-text-muted);
}

.badge-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge-date {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  flex-shrink: 0;
}
</style>
