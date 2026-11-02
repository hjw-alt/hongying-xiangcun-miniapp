<script setup>
import { ref } from 'vue'
import { assetUrl, uploadImage, uploadVideo } from '../api'

const props = defineProps({
  modelValue: { type: String, default: '' },
  kind: { type: String, default: 'image' },
  label: { type: String, default: '上传文件' },
  hint: { type: String, default: '' }
})

const emit = defineEmits(['update:modelValue'])
const busy = ref(false)

async function onPick(event) {
  const file = event.target.files && event.target.files[0]
  if (!file) return
  busy.value = true
  try {
    const url = props.kind === 'video' ? await uploadVideo(file) : await uploadImage(file)
    emit('update:modelValue', url)
  } catch (error) {
    window.dispatchEvent(new CustomEvent('admin-notify', { detail: error.message }))
  } finally {
    busy.value = false
    event.target.value = ''
  }
}

function clear() {
  emit('update:modelValue', '')
}
</script>

<template>
  <div class="upload-field">
    <div class="upload-preview">
      <video v-if="modelValue && kind === 'video'" :src="assetUrl(modelValue)" muted />
      <img v-else-if="modelValue" :src="assetUrl(modelValue)" alt="" />
      <span v-else>暂无文件</span>
    </div>
    <div class="upload-actions">
      <label class="btn" :class="{ disabled: busy }">
        {{ busy ? '上传中…' : label }}
        <input type="file" :accept="kind === 'video' ? 'video/*' : 'image/*'" hidden @change="onPick" />
      </label>
      <button v-if="modelValue" class="btn btn-danger" type="button" @click="clear">清除</button>
      <span v-if="hint" class="hint">{{ hint }}</span>
    </div>
  </div>
</template>
