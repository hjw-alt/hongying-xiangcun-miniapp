const api = require('../../utils/request.js')

const defaultProfile = api.decorateProfile({
  nickname: '湘村体验者',
  points: 0,
  experience: 0,
  level: 2,
  levelName: '初见',
  currentLevelExp: 0,
  nextLevelExp: 30,
  progress: 0,
  coupons: 0,
  balanceCents: 0,
  medalCount: 0
})

Page({
  data: {
    userKey: '',
    profile: defaultProfile,
    medals: [],
    unlockedMedals: [],
    cartCount: 0,
    addressCount: 0,
    tabs: [
      { key: 'orders', title: '我的订单' },
      { key: 'checkins', title: '研学打卡' },
      { key: 'ledger', title: '积分流水' },
      { key: 'favorites', title: '我的收藏' }
    ],
    activeTab: 'orders',
    orders: [],
    records: [],
    ledger: [],
    points: 0,
    loading: true
  },

  onShow() {
    this.setData({ userKey: api.getUserKey() })
    this.loadData()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [profile, medalsResult, orders, taskResult, cartCount, addresses, ledgerResult] = await Promise.all([
        api.request('/api/profile'),
        api.request('/api/medals'),
        api.request('/api/orders'),
        api.request('/api/tasks'),
        api.request('/api/cart/count'),
        api.request('/api/addresses'),
        api.request('/api/points/ledger')
      ])
      const decoratedMedals = (medalsResult.items || []).map(api.decorateMedal)
      const decoratedProfile = api.decorateProfile({
        ...profile,
        nickname: profile.nickname || '湘村体验者'
      })
      this.setData({
        profile: decoratedProfile,
        medals: decoratedMedals,
        unlockedMedals: decoratedMedals.filter(item => item.unlocked),
        orders: (orders || []).map(api.decorateOrder),
        records: taskResult.records || [],
        ledger: (ledgerResult.items || []).map(api.decorateLedger),
        points: Number(decoratedProfile.points || taskResult.points || 0),
        cartCount: cartCount.count || 0,
        addressCount: (addresses || []).length,
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onTab(event) {
    this.setData({ activeTab: event.currentTarget.dataset.key })
  },

  goCart() {
    wx.navigateTo({ url: '/pages/cart/cart' })
  },

  goAddress() {
    wx.navigateTo({ url: '/pages/address/address' })
  },

  goStudy() {
    wx.switchTab({ url: '/pages/study/study' })
  },

  goMall() {
    wx.switchTab({ url: '/pages/mall/mall' })
  },

  goMedals() {
    wx.navigateTo({ url: '/pages/medals/medals' })
  },

  openMedal(event) {
    wx.navigateTo({ url: `/pages/medal-detail/medal-detail?id=${event.currentTarget.dataset.id}` })
  },

  goLedger() {
    this.setData({ activeTab: 'ledger' })
  },

  goStory(event) {
    wx.navigateTo({ url: `/pages/story-detail/story-detail?id=${event.currentTarget.dataset.id}` })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 我的成长', path: '/pages/profile/profile' }
  }
})
