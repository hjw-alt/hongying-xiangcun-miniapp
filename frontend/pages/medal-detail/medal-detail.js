const api = require('../../utils/request.js')

Page({
  data: {
    id: '',
    medal: null,
    loading: true
  },

  onLoad(options) {
    this.setData({ id: options.id || '' })
  },

  onShow() {
    if (this.data.id) this.loadDetail()
  },

  async loadDetail() {
    this.setData({ loading: true })
    try {
      const medal = await api.request(`/api/medals/${this.data.id}`)
      this.setData({ medal: api.decorateMedal(medal), loading: false })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  goMedals() {
    wx.navigateBack({ delta: 1 })
  },

  onShareAppMessage() {
    return { title: this.data.medal ? `红映湘村 · ${this.data.medal.name}` : '红映湘村徽章', path: `/pages/medal-detail/medal-detail?id=${this.data.id}` }
  }
})
