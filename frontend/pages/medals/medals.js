const api = require('../../utils/request.js')

Page({
  data: {
    profile: null,
    medals: [],
    unlockedCount: 0,
    loading: true
  },

  onShow() {
    this.loadData()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [profile, result] = await Promise.all([
        api.request('/api/profile'),
        api.request('/api/medals')
      ])
      const medals = (result.items || []).map(api.decorateMedal)
      this.setData({
        profile: api.decorateProfile(profile),
        medals,
        unlockedCount: medals.filter(item => item.unlocked).length,
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  openDetail(event) {
    const id = event.currentTarget.dataset.id
    if (id) wx.navigateTo({ url: `/pages/medal-detail/medal-detail?id=${id}` })
  },

  goProfile() {
    wx.navigateBack({ delta: 1 })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 我的徽章', path: '/pages/medals/medals' }
  }
})
