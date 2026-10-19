<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{ open: boolean; recommendedLevel: number; level: number }>()
const emit = defineEmits<{ confirm: []; cancel: [] }>()

const dialog = ref<HTMLDialogElement | null>(null)

watch(
  () => props.open,
  (open) => {
    const el = dialog.value
    if (!el) return
    if (open && !el.open) el.showModal()
    if (!open && el.open) el.close()
  },
)

function onBackdrop(e: MouseEvent) {
  if (e.target === dialog.value) emit('cancel')
}
</script>

<template>
  <dialog ref="dialog" class="modal" @cancel.prevent="emit('cancel')" @click="onBackdrop">
    <div class="body">
      <h2>Ты не готов</h2>
      <p>
        Эта карта рассчитана на уровень {{ recommendedLevel }}, а ты пока на {{ level }}. Будет
        тяжело, ты уверен?
      </p>
      <div class="actions">
        <button class="btn btn-ghost" @click="emit('cancel')">Вернуться</button>
        <button class="btn btn-primary" @click="emit('confirm')">Всё равно пойду</button>
      </div>
    </div>
  </dialog>
</template>

<style scoped>
.modal {
  border: 1px solid var(--border);
  border-radius: 18px;
  background: var(--surface);
  color: var(--text);
  padding: 0;
  max-width: 400px;
  width: calc(100% - 32px);
}
.modal::backdrop {
  background: rgb(0 0 0 / 0.6);
}
.body {
  padding: 24px;
}
h2 {
  margin: 0 0 8px;
  color: var(--accent);
}
p {
  margin: 0 0 20px;
  color: var(--muted);
}
.actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  flex-wrap: wrap;
}
</style>
