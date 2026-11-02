const api = require('../../utils/request.js')

Page({
  data: {
    id: '',
    activity: null,
    loading: true,
    errorMessage: ''
  },

  onLoad(options) {
    this.setData({ id: options.id || '' })
  },

  onShow() {
    if (this.data.id) {
      this.loadDetail()
    } else {
      this.setData({ loading: false, errorMessage: '缺少活动信息' })
    }
  },

  async loadDetail() {
    this.setData({ loading: true, errorMessage: '' })
    try {
      const item = await api.request(`/api/activities/${this.data.id}`)
      this.setData({
        activity: api.decorateActivity(item),
        loading: false
      })
    } catch (error) {
      this.setData({
        loading: false,
        errorMessage: error.message || '活动详情读取失败'
      })
    }
  },

  onBook() {
    const activity = this.data.activity
    if (!activity) return
    if (activity.bookingDisabled) {
      wx.showToast({ title: '名额已满，可以先看看别的活动', icon: 'none' })
      return
    }
    getApp().globalData.bookingActivityId = activity.id
    getApp().globalData.experienceTab = 'events'
    wx.switchTab({ url: '/pages/experience/experience' })
  },

  onShareAppMessage() {
    const activity = this.data.activity
    return {
      title: activity ? `影绘湘潭 · ${activity.title}` : '影绘湘潭研学活动',
      path: `/pages/activity/activity?id=${this.data.id}`
    }
  }
})
