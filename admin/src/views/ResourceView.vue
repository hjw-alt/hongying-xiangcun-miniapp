<script setup>
import { computed, inject, reactive, ref, watch } from 'vue'
import { assetUrl, request } from '../api'
import { resources } from '../config'
import UploadField from '../components/UploadField.vue'

const props = defineProps({
  resourceKey: { type: String, required: true }
})

const notify = inject('notify')
const resource = computed(() => resources[props.resourceKey])
const items = ref([])
const refs = reactive({})
const loading = ref(true)
const keyword = ref('')
const drawerOpen = ref(false)
const saving = ref(false)
const editingId = ref('')
const form = reactive({})

function emptyForm() {
  const data = {}
  for (const field of resource.value.fields) {
    data[field.key] = field.type === 'bool' ? 0 : field.type === 'number' ? 0 : ''
  }
  return data
}

function resetForm() {
  const data = emptyForm()
  Object.keys(form).forEach(key => delete form[key])
  Object.assign(form, data)
}

async function loadRefs() {
  const needed = [...new Set(resource.value.fields.filter(f => f.type === 'ref').map(f => f.ref))]
  for (const name of needed) {
    if (refs[name]) continue
    try {
      const result = await request(`/api/admin/data/${name}`)
      refs[name] = result.items || []
    } catch (error) {
      refs[name] = []
    }
  }
}

async function load() {
  loading.value = true
  try {
    const result = await request(`/api/admin/data/${props.resourceKey}`)
    items.value = result.items || []
  } catch (error) {
    notify(error.message)
  } finally {
    loading.value = false
  }
}

function labelFor(field, value) {
  if (!value) return '-'
  const list = refs[field.ref] || []
  const found = list.find(item => item.id === value)
  return found ? (found.name || found.id) : value
}

function cellText(column, row) {
  const value = row[column.key]
  if (column.type === 'ref') return labelFor({ ref: column.ref }, value)
  if (column.type === 'list') return Array.isArray(value) ? value.join('、') : value
  return value === null || value === undefined || value === '' ? '-' : String(value)
}

function formatValue(field, raw) {
  if (field.type === 'list') {
    if (Array.isArray(raw)) return raw.join('\n')
    if (typeof raw === 'string' && raw.startsWith('[')) {
      try {
        return JSON.parse(raw).join('\n')
      } catch (error) {
        return raw
      }
    }
    return raw || ''
  }
  if (field.type === 'bool') return raw ? 1 : 0
  return raw === null || raw === undefined ? '' : raw
}

function openCreate() {
  editingId.value = ''
  resetForm()
  drawerOpen.value = true
}

function openEdit(row) {
  editingId.value = row.id
  const data = {}
  for (const field of resource.value.fields) {
    data[field.key] = formatValue(field, row[field.key])
  }
  Object.keys(form).forEach(key => delete form[key])
  Object.assign(form, data)
  drawerOpen.value = true
}

function payload() {
  const data = {}
  for (const field of resource.value.fields) {
    const raw = form[field.key]
    // 新建记录时允许后端自动生成 ID；编辑时仍保留现有 ID 但后端会忽略主键更新。
    if (field.key === 'id' && !editingId.value && !String(raw || '').trim()) continue
    if (field.type === 'list') {
      data[field.key] = String(raw || '').split('\n').map(line => line.trim()).filter(Boolean)
    } else if (field.type === 'number') {
      data[field.key] = Number(raw || 0)
    } else if (field.type === 'bool') {
      data[field.key] = raw ? 1 : 0
    } else {
      data[field.key] = raw === undefined || raw === null ? '' : raw
    }
  }
  return data
}

async function save() {
  if (saving.value) return
  saving.value = true
  try {
    if (editingId.value) {
      await request(`/api/admin/data/${props.resourceKey}/${editingId.value}`, { method: 'PUT', body: payload() })
    } else {
      await request(`/api/admin/data/${props.resourceKey}`, { method: 'POST', body: payload() })
    }
    drawerOpen.value = false
    notify(editingId.value ? '已保存修改' : '已新增记录')
    await load()
  } catch (error) {
    notify(error.message)
  } finally {
    saving.value = false
  }
}

async function remove(row) {
  if (!window.confirm(`确定删除「${row.name || row.title || row.id}」吗？`)) return
  try {
    await request(`/api/admin/data/${props.resourceKey}/${row.id}`, { method: 'DELETE' })
    notify('已删除')
    await load()
  } catch (error) {
    notify(error.message)
  }
}

const filtered = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  if (!key) return items.value
  return items.value.filter(row => JSON.stringify(row).toLowerCase().includes(key))
})

watch(() => props.resourceKey, async () => {
  keyword.value = ''
  resetForm()
  await loadRefs()
  await load()
}, { immediate: true })
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <h1 class="page-title">{{ resource.label }}</h1>
        <p class="page-desc">{{ resource.desc }} · 共 {{ items.length }} 条</p>
      </div>
      <button class="btn btn-primary" type="button" @click="openCreate">新增{{ resource.label }}</button>
    </div>

    <div class="toolbar">
      <input v-model="keyword" placeholder="搜索当前列表" />
      <button class="btn" type="button" @click="load">刷新</button>
    </div>

    <div v-if="loading" class="empty">正在加载…</div>
    <div v-else-if="!filtered.length" class="empty">暂无数据</div>
    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th v-for="column in resource.listColumns" :key="column.key" :style="column.width ? `width:${column.width}` : ''">
              {{ column.label }}
            </th>
            <th style="width: 150px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in filtered" :key="row.id">
            <td v-for="column in resource.listColumns" :key="column.key">
              <img v-if="column.type === 'image'" class="thumb" :src="assetUrl(row[column.key])" alt="" />
              <span v-else>{{ cellText(column, row) }}</span>
            </td>
            <td>
              <div class="row-actions">
                <button class="btn" type="button" @click="openEdit(row)">编辑</button>
                <button class="btn btn-danger" type="button" @click="remove(row)">删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="drawerOpen" class="drawer-mask" @click.self="drawerOpen = false">
      <div class="drawer">
        <div class="drawer-head">
          <div class="drawer-title">{{ editingId ? '编辑' : '新增' }}{{ resource.label }}</div>
          <button class="btn" type="button" @click="drawerOpen = false">关闭</button>
        </div>
        <div class="drawer-body">
          <div v-for="field in resource.fields" :key="field.key" class="field">
            <label>{{ field.label }}</label>
            <input v-if="field.type === 'text'" v-model="form[field.key]" />
            <input v-else-if="field.type === 'number'" v-model="form[field.key]" type="number" />
            <textarea v-else-if="field.type === 'textarea'" v-model="form[field.key]" :style="field.big ? 'min-height:150px' : ''" />
            <textarea v-else-if="field.type === 'list'" v-model="form[field.key]" style="min-height:90px" placeholder="每行一条" />
            <select v-else-if="field.type === 'select'" v-model="form[field.key]">
              <option value="">请选择</option>
              <option v-for="option in field.options" :key="option" :value="option">{{ option }}</option>
            </select>
            <select v-else-if="field.type === 'ref'" v-model="form[field.key]">
              <option value="">请选择</option>
              <option v-for="option in refs[field.ref] || []" :key="option.id" :value="option.id">
                {{ option.name || option.title || option.id }}
              </option>
            </select>
            <UploadField v-else-if="field.type === 'image'" v-model="form[field.key]" label="上传图片" />
            <label v-else-if="field.type === 'bool'" class="bool-row">
              <input v-model="form[field.key]" type="checkbox" :true-value="1" :false-value="0" />
              <span>勾选表示启用</span>
            </label>
            <div v-if="field.hint" class="hint">{{ field.hint }}</div>
          </div>
        </div>
        <div class="drawer-foot">
          <button class="btn" type="button" @click="drawerOpen = false">取消</button>
          <button class="btn btn-primary" type="button" :disabled="saving" @click="save">
            {{ saving ? '保存中…' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
