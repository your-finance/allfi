<script setup>
/**
 * PinInput - PIN 格子输入组件
 * 显示 N 个格子，用户输入时自动跳转
 */
import { ref, watch, onMounted, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  length: {
    type: Number,
    default: 6
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const inputs = ref([])
const pins = ref([])

// 初始化格子
onMounted(() => {
  pins.value = props.modelValue.split('').concat(
    Array(props.length - props.modelValue.length).fill('')
  )
})

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  pins.value = newVal.split('').concat(
    Array(props.length - newVal.length).fill('')
  )
})

// 处理输入
const handleInput = (index, event) => {
  const value = event.target.value

  // 只允许数字
  if (!/^\d*$/.test(value)) {
    event.target.value = pins.value[index]
    return
  }

  pins.value[index] = value.slice(-1)

  // 更新 modelValue
  const newValue = pins.value.join('')
  emit('update:modelValue', newValue)

  // 自动跳转到下一个格子
  if (value && index < props.length - 1) {
    nextTick(() => {
      inputs.value[index + 1]?.focus()
    })
  }

  // 检查是否完成
  if (newValue.length === props.length) {
    emit('complete', newValue)
  }
}

// 处理退格
const handleKeydown = (index, event) => {
  if (event.key === 'Backspace' && !pins.value[index] && index > 0) {
    nextTick(() => {
      inputs.value[index - 1]?.focus()
    })
  }
}

// 处理粘贴
const handlePaste = (event) => {
  event.preventDefault()
  const paste = event.clipboardData.getData('text')

  if (!/^\d+$/.test(paste)) return

  const chars = paste.slice(0, props.length).split('')
  pins.value = chars.concat(Array(props.length - chars.length).fill(''))

  const newValue = pins.value.join('')
  emit('update:modelValue', newValue)

  if (newValue.length === props.length) {
    emit('complete', newValue)
  }
}
</script>

<template>
  <div class="pin-input">
    <input
      v-for="(_, index) in length"
      :key="index"
      ref="inputs"
      type="password"
      inputmode="numeric"
      maxlength="1"
      :value="pins[index]"
      :disabled="disabled"
      class="pin-cell"
      @input="handleInput(index, $event)"
      @keydown="handleKeydown(index, $event)"
      @paste="index === 0 && handlePaste($event)"
    />
  </div>
</template>

<style scoped>
.pin-input {
  display: flex;
  gap: var(--gap-sm);
  justify-content: center;
}

.pin-cell {
  width: 48px;
  height: 56px;
  text-align: center;
  font-size: 24px;
  font-family: var(--font-mono);
  background: var(--color-bg-tertiary);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  transition: border-color var(--transition-fast);
}

.pin-cell:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.pin-cell:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .pin-cell {
    width: 40px;
    height: 48px;
    font-size: 20px;
  }
}
</style>
