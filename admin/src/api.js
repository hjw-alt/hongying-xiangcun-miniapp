const TOKEN_KEY = 'shadowart-admin-token'
const BASE = import.meta.env.VITE_API_BASE || ''

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token) {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
  } else {
    localStorage.removeItem(TOKEN_KEY)
  }
}

async function parseBody(response) {
  const text = await response.text()
  if (!text) return null
  try {
    return JSON.parse(text)
  } catch (error) {
    return null
  }
}

export async function request(path, options = {}) {
  const { method = 'GET', body, auth = true } = options
  const headers = {}
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  if (auth && getToken()) {
    headers['X-Admin-Token'] = getToken()
  }

  const response = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body)
  })
  const data = await parseBody(response)

  if (!response.ok) {
    if (response.status === 401) {
      setToken('')
      window.dispatchEvent(new CustomEvent('admin-unauthorized'))
    }
    throw new ApiError((data && data.error) || `请求失败（${response.status}）`, response.status)
  }
  return data
}

async function upload(path, file) {
  const form = new FormData()
  form.append('file', file)
  const headers = {}
  if (getToken()) {
    headers['X-Admin-Token'] = getToken()
  }
  const response = await fetch(`${BASE}${path}`, { method: 'POST', headers, body: form })
  const data = await parseBody(response)
  if (!response.ok || !data || !data.url) {
    throw new ApiError((data && data.error) || '文件上传失败', response.status)
  }
  return data.url
}

export const uploadImage = file => upload('/api/upload', file)
export const uploadVideo = file => upload('/api/upload/video', file)

export function assetUrl(url) {
  if (!url) return ''
  if (/^https?:\/\//.test(url)) return url
  return `${BASE}${url}`
}

export function formatDateTime(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = number => String(number).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function toLocalInput(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = number => String(number).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function toRFC3339(local) {
  if (!local) return ''
  return `${local}:00+08:00`
}
