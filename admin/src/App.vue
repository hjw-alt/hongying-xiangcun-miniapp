<script setup>
import { onMounted, provide, reactive, ref } from 'vue'
import { getToken, request, setToken } from './api'
import { resources } from './config'
import DashboardView from './views/DashboardView.vue'
import OrdersView from './views/OrdersView.vue'
import ResourceView from './views/ResourceView.vue'

const menu = [
  { key: 'dashboard', label: '数据概览', type: 'dashboard' },
  { key: 'orders', label: '订单管理', type: 'orders' },
  { key: 'banners', label: '首页轮播', type: 'resource' },
  { key: 'villages', label: '红色村落', type: 'resource' },
  { key: 'spots', label: '村落景点', type: 'resource' },
  { key: 'routes', label: '红色线路', type: 'resource' },
  { key: 'stories', label: '红迹故事', type: 'resource' },
  { key: 'products', label: '红文创商品', type: 'resource' },
  { key: 'farmers', label: '合作农户', type: 'resource' },
  { key: 'tasks', label: '研学任务', type: 'resource' },
  { key: 'medals', label: '徽章管理', type: 'resource' },
  { key: 'point_rules', label: '积分/成长规则', type: 'resource' }
]

const current = ref('dashboard')
const authed = ref(false)
const ready = ref(false)
const username = ref('')
const loginForm = reactive({ username: '', password: '' })
const loginError = ref('')
const logging = ref(false)
const toastMessage = ref('')
let toastTimer = 0

function notify(message) {
  toastMessage.value = message
  window.clearTimeout(toastTimer)
  toastTimer = window.setTimeout(() => {
    toastMessage.value = ''
  }, 2400)
}

provide('notify', notify)

async function verifySession() {
  if (!getToken()) {
    ready.value = true
    return
  }
  try {
    const session = await request('/api/admin/session')
    username.value = session.username
    authed.value = true
  } catch (error) {
    setToken('')
  } finally {
    ready.value = true
  }
}

async function submitLogin() {
  if (logging.value) return
  loginError.value = ''
  logging.value = true
  try {
    const result = await request('/api/admin/login', {
      method: 'POST',
      auth: false,
      body: { username: loginForm.username, password: loginForm.password }
    })
    setToken(result.token)
    username.value = result.username
    authed.value = true
    current.value = 'dashboard'
    notify('登录成功')
  } catch (error) {
    loginError.value = error.message
  } finally {
    logging.value = false
  }
}

function logout() {
  setToken('')
  authed.value = false
  username.value = ''
  loginForm.password = ''
  notify('已退出登录')
}

onMounted(() => {
  window.addEventListener('admin-unauthorized', () => {
    authed.value = false
    username.value = ''
    notify('登录已过期，请重新登录')
  })
  window.addEventListener('admin-notify', event => notify(event.detail))
  verifySession()
})
</script>

<template>
  <div v-if="!ready" class="login-page">
    <div class="login-card">
      <div class="login-sub">正在检查登录状态…</div>
    </div>
  </div>

  <div v-else-if="!authed" class="login-page">
    <form class="login-card" @submit.prevent="submitLogin">
      <span class="login-seal">红映湘村</span>
      <h1 class="login-title">内容管理后台</h1>
      <p class="login-sub">村落 · 线路 · 故事 · 文创 · 农户 · 研学 · 订单</p>
      <div v-if="loginError" class="login-error">{{ loginError }}</div>
      <div class="login-field">
        <label for="admin-username">管理员账号</label>
        <input id="admin-username" v-model="loginForm.username" autocomplete="username" placeholder="请输入账号" />
      </div>
      <div class="login-field">
        <label for="admin-password">密码</label>
        <input id="admin-password" v-model="loginForm.password" type="password" autocomplete="current-password" placeholder="请输入密码" />
      </div>
      <button class="btn btn-primary login-submit" type="submit" :disabled="logging">
        {{ logging ? '登录中…' : '登录' }}
      </button>
    </form>
  </div>

  <div v-else class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-name">红映湘村</div>
        <div class="brand-sub">内容管理后台</div>
      </div>
      <nav class="nav">
        <button
          v-for="item in menu"
          :key="item.key"
          class="nav-item"
          :class="{ active: current === item.key }"
          type="button"
          @click="current = item.key"
        >{{ item.label }}</button>
      </nav>
      <div class="sidebar-foot">
        <div class="sidebar-user">{{ username }}</div>
        <button class="btn" type="button" style="margin-top: 10px; width: 100%" @click="logout">退出登录</button>
      </div>
    </aside>

    <main class="main">
      <DashboardView v-if="current === 'dashboard'" />
      <OrdersView v-else-if="current === 'orders'" />
      <ResourceView v-else :key="current" :resource-key="current" />
    </main>
  </div>

  <div v-if="toastMessage" class="toast">{{ toastMessage }}</div>
</template>
