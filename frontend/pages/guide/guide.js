const api = require('../../utils/request.js')

Page({
  data: {
    villages: [],
    routes: [],
    markers: [],
    selectedRoute: null,
    keyword: '',
    center: { latitude: 27.829, longitude: 112.6 },
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 2, hidden: false })
    }
    this.loadData()
  },

  setTabBarVisible(visible) {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ hidden: !visible })
    }
  },

  onHide() {
    this.setData({ selectedRoute: null })
    this.setTabBarVisible(true)
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [villageResult, routes] = await Promise.all([
        api.request(`/api/villages?keyword=${encodeURIComponent(this.data.keyword)}`),
        api.request('/api/routes')
      ])
      const villages = (villageResult.items || []).map(api.decorateVillage)
      this.setData({
        villages,
        routes: routes.map(api.decorateRoute),
        markers: villages.filter(item => item.lat && item.lng).map(item => ({
          id: item.id,
          latitude: item.lat,
          longitude: item.lng,
          width: 28,
          height: 28,
          callout: {
            content: item.name,
            color: '#7E1B1B',
            fontSize: 12,
            borderRadius: 6,
            padding: 6,
            display: 'ALWAYS'
          }
        })),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onKeyword(event) {
    this.setData({ keyword: event.detail.value })
  },

  onSearch() {
    this.loadData()
  },

  onMarkerTap(event) {
    const id = event.detail.markerId
    wx.navigateTo({ url: `/pages/village/village?id=${id}` })
  },

  onVillage(event) {
    wx.navigateTo({ url: `/pages/village/village?id=${event.currentTarget.dataset.id}` })
  },

  onRoute(event) {
    const route = this.data.routes.find(item => item.id === event.currentTarget.dataset.id)
    if (!route) return
    this.setData({ selectedRoute: route })
    this.setTabBarVisible(false)
  },

  closeRoute() {
    this.setData({ selectedRoute: null })
    this.setTabBarVisible(true)
  },

  noop() {},

  onRouteSpot(event) {
    const villageId = event.currentTarget.dataset.village || (this.data.villages[0] && this.data.villages[0].id)
    if (!villageId) return
    this.setData({ selectedRoute: null })
    this.setTabBarVisible(true)
    wx.navigateTo({ url: `/pages/village/village?id=${villageId}` })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 红色村落导览', path: '/pages/guide/guide' }
  }
})
