const api = require('../../utils/request.js')

Page({
  data: {
    summary: { farmers: 0, products: 0, orders: 0, supportCount: 0 },
    farmers: [],
    products: [],
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 5, hidden: false })
    }
    this.loadData()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [summary, farmers, productResult] = await Promise.all([
        api.request('/api/farm-summary'),
        api.request('/api/farmers'),
        api.request('/api/products?category=' + encodeURIComponent('湘村风物助农'))
      ])
      this.setData({
        summary,
        farmers: farmers.map(api.decorateFarmer),
        products: (productResult.items || []).map(api.decorateProduct),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onProduct(event) {
    wx.navigateTo({ url: `/pages/product/product?id=${event.currentTarget.dataset.id}` })
  },

  goMall() {
    wx.switchTab({ url: '/pages/mall/mall' })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 助农帮扶专区', path: '/pages/farm/farm' }
  }
})
