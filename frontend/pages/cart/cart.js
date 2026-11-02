const api = require('../../utils/request.js')

const emptyProfile = api.decorateProfile({
  points: 0,
  level: 2,
  levelName: '初见',
  experience: 0,
  nextLevelExp: 30,
  progress: 0
})

Page({
  data: {
    items: [],
    profile: emptyProfile,
    selectedIds: [],
    cashTotalText: '¥0.00',
    pointsTotalText: '0 积分',
    totalText: '¥0.00',
    rewardPointsTotal: 0,
    rewardExpTotal: 0,
    paymentType: 'cash',
    paymentTypeText: '现金购买',
    paymentHint: '现金购买后可获得积分与成长经验',
    selectedCount: 0,
    allSelected: false,
    canPayPoints: false,
    canCheckout: false,
    loading: true
  },

  onShow() {
    this.loadData()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [result, profile] = await Promise.all([
        api.request('/api/cart'),
        api.request('/api/profile')
      ])
      const items = (result.items || []).map(item => ({
        ...item,
        product: item.product ? api.decorateProduct(item.product) : null,
        selected: this.data.selectedIds.includes(item.id)
      }))
      this.setData({
        items,
        profile: api.decorateProfile(profile),
        loading: false
      }, () => this.refreshSummary())
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  refreshSummary() {
    const validItems = this.data.items.filter(item => item.product)
    const selected = validItems.filter(item => item.selected)
    const cashTotal = selected.reduce((sum, item) => sum + Number(item.product.priceCents || 0) * item.quantity, 0)
    const pointsTotal = selected.reduce((sum, item) => sum + Number(item.product.pointsPrice || 0) * item.quantity, 0)
    const rewardPointsTotal = selected.reduce((sum, item) => sum + Number(item.product.rewardPoints || 0) * item.quantity, 0)
    const rewardExpTotal = selected.reduce((sum, item) => sum + Number(item.product.rewardExp || 0) * item.quantity, 0)
    const pointsUnavailable = selected.some(item => Number(item.product.pointsPrice || 0) <= 0)
    const pointsEnough = Number(this.data.profile.points || 0) >= pointsTotal
    const canPayPoints = selected.length > 0 && !pointsUnavailable && pointsEnough
    const paymentType = this.data.paymentType === 'points' && canPayPoints ? 'points' : 'cash'
    let paymentHint = '现金购买后可获得积分与成长经验'
    if (selected.length === 0) {
      paymentHint = '请选择商品后再选择支付方式'
    } else if (pointsUnavailable) {
      paymentHint = '所选商品中有暂不支持积分兑换的商品'
    } else if (!pointsEnough) {
      paymentHint = `积分兑换还差 ${Math.max(0, pointsTotal - Number(this.data.profile.points || 0))} 积分`
    } else if (paymentType === 'points') {
      paymentHint = `本次将消耗 ${pointsTotal} 积分，并获得 ${rewardExpTotal} 经验`
    } else {
      paymentHint = `现金购买预计获得 ${rewardPointsTotal} 积分 · ${rewardExpTotal} 经验`
    }
    const selectedIds = selected.map(item => item.id)
    this.setData({
      selectedIds,
      selectedCount: selected.length,
      cashTotalText: `¥${(cashTotal / 100).toFixed(2)}`,
      pointsTotalText: `${pointsTotal} 积分`,
      totalText: paymentType === 'points' ? `${pointsTotal} 积分` : `¥${(cashTotal / 100).toFixed(2)}`,
      rewardPointsTotal,
      rewardExpTotal,
      paymentType,
      paymentTypeText: paymentType === 'points' ? '积分兑换' : '现金购买',
      paymentHint,
      canPayPoints,
      canCheckout: selected.length > 0 && (paymentType === 'cash' || canPayPoints),
      allSelected: validItems.length > 0 && selected.length === validItems.length
    })
  },

  toggleItem(event) {
    const id = event.currentTarget.dataset.id
    const items = this.data.items.map(item => (item.id === id ? { ...item, selected: !item.selected } : item))
    this.setData({ items }, () => this.refreshSummary())
  },

  toggleAll() {
    const next = !this.data.allSelected
    const items = this.data.items.map(item => ({ ...item, selected: next }))
    this.setData({ items }, () => this.refreshSummary())
  },

  choosePayment(event) {
    const type = event.currentTarget.dataset.type
    if (type === 'points' && !this.data.canPayPoints) {
      wx.showToast({ title: this.data.paymentHint || '当前商品不支持积分兑换', icon: 'none' })
      return
    }
    this.setData({ paymentType: type }, () => this.refreshSummary())
  },

  async changeQuantity(event) {
    const id = event.currentTarget.dataset.id
    const delta = event.currentTarget.dataset.type === 'add' ? 1 : -1
    const item = this.data.items.find(row => row.id === id)
    if (!item) return
    const quantity = Math.min(99, Math.max(1, item.quantity + delta))
    if (quantity === item.quantity) return
    try {
      await api.request(`/api/cart/${id}`, 'PUT', { quantity, customText: item.customText })
      const items = this.data.items.map(row => (row.id === id ? { ...row, quantity } : row))
      this.setData({ items }, () => this.refreshSummary())
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  removeItem(event) {
    const id = event.currentTarget.dataset.id
    wx.showModal({
      title: '移出购物车',
      content: '确定要移除这件商品吗？',
      confirmColor: '#B32A2A',
      success: async result => {
        if (!result.confirm) return
        try {
          const res = await api.request(`/api/cart/${id}`, 'DELETE')
          getApp().globalData.cartCount = res.count
          await this.loadData()
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  },

  checkout() {
    const selected = this.data.items.filter(item => item.selected && item.product)
    if (selected.length === 0) {
      wx.showToast({ title: '请先选择要结算的商品', icon: 'none' })
      return
    }
    if (!this.data.canCheckout) {
      wx.showToast({ title: this.data.paymentHint || '当前无法结算', icon: 'none' })
      return
    }
    getApp().globalData.checkoutDraft = {
      source: 'cart',
      paymentType: this.data.paymentType,
      items: selected.map(item => ({
        productId: item.productId,
        quantity: item.quantity,
        customText: item.customText,
        name: item.product.name,
        priceCents: item.product.priceCents,
        pointsPrice: item.product.pointsPrice,
        rewardPoints: item.product.rewardPoints,
        rewardExp: item.product.rewardExp,
        coverSrc: item.product.coverSrc,
        farmerName: item.product.farmerName
      })),
      cartItemIds: selected.map(item => item.id)
    }
    wx.navigateTo({ url: '/pages/checkout/checkout' })
  },

  goMall() {
    wx.switchTab({ url: '/pages/mall/mall' })
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 红文创商城', path: '/pages/mall/mall' }
  }
})
