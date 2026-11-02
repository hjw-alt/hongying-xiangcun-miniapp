const api = require('../../utils/request.js')

const emptyProfile = api.decorateProfile({
  points: 0,
  level: 2,
  levelName: '初见',
  experience: 0,
  currentLevelExp: 0,
  nextLevelExp: 30,
  progress: 0
})

function decorateRecord(item = {}) {
  return {
    ...item,
    createdAtText: api.formatDateTime(item.createdAt)
  }
}

Page({
  data: {
    tasks: [],
    villages: [],
    villageId: '',
    done: {},
    points: 0,
    records: [],
    profile: emptyProfile,
    medals: [],
    activeTask: null,
    note: '',
    photoUrl: '',
    submitting: false,
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 4, hidden: false })
    }
    this.loadData()
  },

  setTabBarVisible(visible) {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ hidden: !visible })
    }
  },

  onHide() {
    this.setData({ activeTask: null })
    this.setTabBarVisible(true)
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [taskResult, villageResult, medalResult] = await Promise.all([
        api.request(`/api/tasks?villageId=${encodeURIComponent(this.data.villageId)}`),
        api.request('/api/villages'),
        api.request('/api/medals')
      ])
      const profile = api.decorateProfile(taskResult.profile || emptyProfile)
      this.setData({
        tasks: (taskResult.items || []).map(api.decorateTask),
        records: (taskResult.records || []).map(decorateRecord),
        done: taskResult.done || {},
        points: Number(profile.points || taskResult.points || 0),
        profile,
        medals: (medalResult.items || medalResult || []).map(api.decorateMedal),
        villages: (villageResult.items || []).map(api.decorateVillage),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onVillage(event) {
    const id = event.currentTarget.dataset.id
    this.setData({ villageId: this.data.villageId === id ? '' : id })
    this.loadData()
  },

  openTask(event) {
    const id = event.currentTarget.dataset.id
    if (this.data.done[id]) {
      wx.showToast({ title: '这个任务已经打卡了', icon: 'none' })
      return
    }
    const task = this.data.tasks.find(item => item.id === id)
    this.setData({ activeTask: task, note: '', photoUrl: '' })
    this.setTabBarVisible(false)
  },

  closeTask() {
    this.setData({ activeTask: null })
    this.setTabBarVisible(true)
  },

  noop() {},

  onNote(event) {
    this.setData({ note: event.detail.value })
  },

  choosePhoto() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: async result => {
        const file = result.tempFiles && result.tempFiles[0]
        if (!file) return
        wx.showLoading({ title: '上传中', mask: true })
        try {
          const url = await api.uploadImage(file.tempFilePath)
          this.setData({ photoUrl: api.absoluteUrl(url) })
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        } finally {
          wx.hideLoading()
        }
      }
    })
  },

  async submitTask() {
    if (this.data.submitting || !this.data.activeTask) return
    const note = this.data.note.trim()
    if (!note) {
      wx.showToast({ title: '写一句打卡心得吧', icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      const result = await api.request(`/api/tasks/${this.data.activeTask.id}/checkin`, 'POST', {
        note,
        photoUrl: this.data.photoUrl
      })
      const record = result.record ? decorateRecord(result.record) : null
      const nextRecords = record ? [record, ...this.data.records] : this.data.records
      this.setData({
        activeTask: null,
        profile: api.decorateProfile(result.profile || this.data.profile),
        points: Number((result.profile && result.profile.points) || this.data.points),
        medals: (result.medals || this.data.medals).map(api.decorateMedal),
        records: nextRecords,
        done: { ...this.data.done, [this.data.activeTask.id]: true }
      })
      this.setTabBarVisible(true)
      wx.showToast({ title: '打卡成功，积分和经验已到账', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' })
  },

  goMedals() {
    wx.navigateTo({ url: '/pages/medals/medals' })
  },

  openMedal(event) {
    const id = event.currentTarget.dataset.id
    if (id) wx.navigateTo({ url: `/pages/medal-detail/medal-detail?id=${id}` })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 研学打卡', path: '/pages/study/study' }
  }
})
