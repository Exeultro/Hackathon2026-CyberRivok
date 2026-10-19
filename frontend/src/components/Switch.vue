<script setup lang="ts">
defineProps<{ modelValue: boolean; label: string; hint?: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>

<template>
  <label class="row">
    <span class="text">
      <b>{{ label }}</b>
      <small v-if="hint">{{ hint }}</small>
    </span>
    <input
      type="checkbox"
      role="switch"
      :checked="modelValue"
      @change="emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
    />
    <span class="track" aria-hidden="true" />
  </label>
</template>

<style scoped>
.row {
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  position: relative;
}
.text {
  flex: 1;
  display: grid;
}
.text small {
  color: var(--muted);
}
input {
  position: absolute;
  opacity: 0;
  inset: 0;
  margin: 0;
  cursor: pointer;
}
.track {
  flex: none;
  width: 46px;
  height: 28px;
  border-radius: 14px;
  background: var(--surface-2);
  border: 2px solid var(--border);
  position: relative;
  transition: background 0.15s;
}
.track::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--text);
  transition: transform 0.15s;
}
input:checked ~ .track {
  background: var(--accent);
  border-color: var(--accent);
}
input:checked ~ .track::after {
  transform: translateX(18px);
  background: var(--accent-ink);
}
input:focus-visible ~ .track {
  outline: 3px solid var(--accent);
  outline-offset: 2px;
}
</style>
