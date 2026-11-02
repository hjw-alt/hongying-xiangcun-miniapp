const api = require('../../utils/request.js')

Page({
  data: {
    id: '',
    village: null,
    spots: [],
    stories: [],
    products: [],
    favorited: false,
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
      const [detail, favorite] = await Promise.all([
        api.request(`/api/villages/${this.data.id}`),
        api.request(`/api/favorites/item/village/${this.data.id}`)
      ])
      this.setData({
        village: api.decorateVillage(detail.village),
        spots: (detail.spots || []).map(api.decorateSpot),
        stories: (detail.stories || []).map(api.decorateStory),
        products: (detail.products || []).map(api.decorateProduct),
        favorited: favorite.favorited,
        loading: false
      })
      wx.setNavigationBarTitle({ title: detail.village.name })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async onFavorite() {
    try {
      const result = await api.request(`/api/favorites/item/village/${this.data.id}/toggle`, 'POST')
      this.setData({ favorited: result.favorited })
      wx.showToast({ title: result.favorited ? '已收藏' : '已取消收藏', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  openLocation() {
    const village = this.data.village
    if (!village.lat || !village.lng) return
    wx.openLocation({
      latitude: village.lat,
      longitude: village.lng,
      name: village.name,
      address: village.region,
      scale: 14
    })
  },

  onStory(event) {
    wx.navigateTo({ url: `/pages/story-detail/story-detail?id=${event.currentTarget.dataset.id}` })
  },

  onProduct(event) {
    wx.navigateTo({ url: `/pages/product/product?id=${event.currentTarget.dataset.id}` })
  },

  onShareAppMessage() {
    const village = this.data.village
    return {
      title: village ? `红映湘村 · ${village.name}` : '红映湘村',
      path: `/pages/village/village?id=${this.data.id}`
    }
  }
})
