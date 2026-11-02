<script setup>
import { onMounted, ref } from 'vue'
import { request } from '../api'

const summary = ref(null)
const loading = ref(true)

const labels = {
  villages: '红色村落',
  spots: '红色景点',
  routes: '主题线路',
  stories: '红迹故事',
  products: '文创商品',
  farmers: '合作农户',
  tasks: '研学任务',
  orders: '订单',
  taskRecords: '打卡记录',
  banners: '首页轮播',
  profiles: '成长账户',
  medals: '启用徽章',
  pointLedger: '积分流水'
}

onMounted(async () => {
  try {
    summary.value = await request('/api/admin/red-summary')
  } catch (error) {
    window.dispatchEvent(new CustomEvent('admin-notify', { detail: error.message }))
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <h1 class="page-title">数据概览</h1>
        <p class="page-desc">小程序端内容与成长体系由后台统一维护：首页轮播、红色村落、线路、故事、文创商品、农户、研学任务、积分规则与徽章。</p>
      </div>
    </div>

    <div v-if="loading" class="empty">正在读取数据…</div>
    <div v-else class="stat-grid">
      <div v-for="(label, key) in labels" :key="key" class="stat-card">
        <div class="stat-label">{{ label }}</div>
        <div class="stat-value">{{ summary ? summary[key] : 0 }}</div>
      </div>
    </div>
  </div>
</template>
