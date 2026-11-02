<script setup>
import { inject, onMounted, ref } from 'vue'
import { formatDateTime, request } from '../api'

const notify = inject('notify')
const items = ref([])
const loading = ref(true)
const statuses = ['待发货', '已发货', '已完成', '已取消']

async function load() {
  loading.value = true
  try {
    const result = await request('/api/admin/orders')
    items.value = (result.items || []).map(item => ({
      ...item,
      totalText: `¥${(item.totalCents / 100).toFixed(2)}`,
      createdAtText: formatDateTime(item.createdAt)
    }))
  } catch (error) {
    notify(error.message)
  } finally {
    loading.value = false
  }
}

async function changeStatus(item, event) {
  const status = event.target.value
  try {
    await request(`/api/admin/orders/${item.id}/status`, { method: 'POST', body: { status } })
    notify('订单状态已更新')
    await load()
  } catch (error) {
    notify(error.message)
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <h1 class="page-title">订单管理</h1>
        <p class="page-desc">文创与助农订单，共 {{ items.length }} 单；修改状态即同步到小程序「我的订单」</p>
      </div>
      <button class="btn" type="button" @click="load">刷新</button>
    </div>

    <div v-if="loading" class="empty">正在加载…</div>
    <div v-else-if="!items.length" class="empty">暂无订单</div>
    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th style="width: 180px">订单号</th>
            <th>商品</th>
            <th style="width: 160px">收货信息</th>
            <th style="width: 100px">金额</th>
            <th style="width: 150px">下单时间</th>
            <th style="width: 140px">状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id">
            <td>{{ item.orderNo }}</td>
            <td>
              <div v-for="row in item.items" :key="row.productId" class="cell-sub">
                {{ row.name }} × {{ row.quantity }}
                <span v-if="row.customText">（定制：{{ row.customText }}）</span>
              </div>
            </td>
            <td>
              <div>{{ item.receiver }} {{ item.phone }}</div>
              <div class="cell-sub">{{ item.address }}</div>
            </td>
            <td>{{ item.totalText }}</td>
            <td>{{ item.createdAtText }}</td>
            <td>
              <select :value="item.status" @change="changeStatus(item, $event)">
                <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
              </select>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
