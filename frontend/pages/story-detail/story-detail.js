const api = require('../../utils/request.js')

Page({
  data: {
    id: '',
    story: null,
    village: null,
    favorited: false,
    arVisible: false,
    loading: true
  },

  onLoad(options) {
    this.setData({ id: options.id || '' })
  },

  onShow() {
    if (this.data.id) this.loadDetail()
  },

  async loadDetail() {
    try {
      const story = api.decorateStory(await api.request(`/api/stories/${this.data.id}`))
      const favorite = await api.request(`/api/favorites/item/story/${this.data.id}`)
      let village = null
      if (story.villageId) {
        try {
          const detail = await api.request(`/api/villages/${story.villageId}`)
          village = api.decorateVillage(detail.village)
        } catch (error) {
          village = null
        }
      }
      this.setData({ story, village, favorited: favorite.favorited, loading: false })
      wx.setNavigationBarTitle({ title: story.title })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async onFavorite() {
    try {
      const result = await api.request(`/api/favorites/item/story/${this.data.id}/toggle`, 'POST')
      this.setData({ favorited: result.favorited })
      wx.showToast({ title: result.favorited ? '已收藏' : '已取消收藏', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  openAR() {
    if (!this.data.story.arSrc) {
      wx.showToast({ title: '该故事暂无 AR 内容', icon: 'none' })
      return
    }
    this.setData({ arVisible: true })
  },

  closeAR() {
    this.setData({ arVisible: false })
  },

  goVillage() {
    if (!this.data.village) return
    wx.navigateTo({ url: `/pages/village/village?id=${this.data.village.id}` })
  },

  onShareAppMessage() {
    const story = this.data.story
    return {
      title: story ? `红映湘村 · ${story.title}` : '红映湘村',
      path: `/pages/story-detail/story-detail?id=${this.data.id}`
    }
  }
})
