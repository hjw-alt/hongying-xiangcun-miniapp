const api = require('../../utils/request.js')

Page({
  data: {
    heroImage: api.absoluteUrl('/static/seed/art-puppet.png'),
    artworks: [],
    keyword: '',
    theme: '全部',
    themes: ['全部', '纸影人物', '湘潭印象', '戏台故事'],
    loading: true
  },

  onShow() {
    this.loadList()
  },

  async loadList() {
    this.setData({ loading: true })
    try {
      const query = `keyword=${encodeURIComponent(this.data.keyword)}&theme=${encodeURIComponent(this.data.theme)}`
      const response = await api.request(`/api/artworks?${query}`)
      this.setData({
        artworks: response.items.map(api.decorateArtwork),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onKeywordInput(event) {
    this.setData({ keyword: event.detail.value })
  },

  onSearch() {
    this.loadList()
  },

  onThemeSelect(event) {
    this.setData({ theme: event.currentTarget.dataset.theme })
    this.loadList()
  },

  onEdit(event) {
    const { id } = event.currentTarget.dataset
    getApp().globalData.editingId = id
    wx.switchTab({ url: '/pages/create/create' })
  },

  onDetail(event) {
    const { id } = event.currentTarget.dataset
    wx.navigateTo({ url: `/pages/detail/detail?id=${id}` })
  },

  onDelete(event) {
    const { id, title } = event.currentTarget.dataset
    wx.showModal({
      title: '删除确认',
      content: `确定要删除《${title}》吗？相关评论也会删除。`,
      confirmColor: '#C8502E',
      success: async result => {
        if (!result.confirm) return
        try {
          await api.request(`/api/artworks/${id}`, 'DELETE')
          wx.showToast({ title: '已删除', icon: 'success' })
          this.loadList()
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  }
})
