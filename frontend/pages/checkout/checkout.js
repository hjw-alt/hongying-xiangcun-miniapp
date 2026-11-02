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
    address: null,
    remark: '',
    paymentType: 'cash',
    paymentTypeText: '现金购买',
    cashTotalText: '¥0.00',
    pointsTotalText: '0 积分',
    totalText: '¥0.00',
    rewardPointsTotal: 0,
    rewardExpTotal: 0,
    paymentHint: '现金购买后可获得积分与成长经验',
    canPayPoints: false,
    canSubmit: false,
    submitting: false,
    loading: true,
    empty: false
  },

  onLoad() {
    this.loadDraft()
    this.loadProfile()
  },

  onShow() {
    const picked = getApp().globalData.selectedAddress
    if (picked) {
      getApp().globalData.selectedAddress = null
      this.setData({ address: picked })
      return
    }
    if (!this.data.address) this.loadDefaultAddress()
  },

  async loadProfile() {
    try {
      const profile = await api.request('/api/profile')
      this.setData({ profile: api.decorateProfile(profile) }, () => this.syncSummary())
    } catch (error) {
      // 允许先浏览结算页，提交时仍由后端校验账户余额。
    }
  },

  async loadDefaultAddress() {
    try {
      const list = await api.request('/api/addresses')
      const preferred = list.find(item => item.isDefault) || list[0] || null
      this.setData({ address: preferred })
    } catch (error) {
      this.setData({ address: null })
    }
  },

  async loadDraft() {
    const draft = getApp().globalData.checkoutDraft
    if (!draft || !draft.items || draft.items.length === 0) {
      this.setData({ loading: false, empty: true })
      return
    }
    const items = []
    for (const row of draft.items) {
      const hasRewards = row.pointsPrice !== undefined && row.rewardExp !== undefined
      if (row.name && row.priceCents !== undefined && hasRewards) {
        items.push({
          ...row,
          priceText: row.priceText || `¥${(Number(row.priceCents || 0) / 100).toFixed(2)}`,
          coverSrc: row.coverSrc || row.coverUrl || ''
        })
        continue
      }
      try {
        const product = api.decorateProduct(await api.request(`/api/products/${row.productId}`))
        items.push({
          ...row,
          name: product.name,
          priceCents: product.priceCents,
          priceText: product.priceText,
          pointsPrice: product.pointsPrice,
          rewardPoints: product.rewardPoints,
          rewardExp: product.rewardExp,
          coverSrc: product.coverSrc,
          farmerName: product.farmerName
        })
      } catch (error) {
        // 商品可能已下架，跳过
      }
    }
    this.setData({
      items,
      paymentType: draft.paymentType || 'cash',
      loading: false,
      empty: items.length === 0
    }, () => this.syncSummary())
  },

  syncSummary() {
    const selected = this.data.items
    const cashTotal = selected.reduce((sum, item) => sum + Number(item.priceCents || 0) * Number(item.quantity || 1), 0)
    const pointsTotal = selected.reduce((sum, item) => sum + Number(item.pointsPrice || 0) * Number(item.quantity || 1), 0)
    const rewardPointsTotal = selected.reduce((sum, item) => sum + Number(item.rewardPoints || 0) * Number(item.quantity || 1), 0)
    const rewardExpTotal = selected.reduce((sum, item) => sum + Number(item.rewardExp || 0) * Number(item.quantity || 1), 0)
    const pointsUnavailable = selected.some(item => Number(item.pointsPrice || 0) <= 0)
    const pointsEnough = Number(this.data.profile.points || 0) >= pointsTotal
    const canPayPoints = selected.length > 0 && !pointsUnavailable && pointsEnough
    const paymentType = this.data.paymentType === 'points' && canPayPoints ? 'points' : 'cash'
    let paymentHint = '现金购买后可获得积分与成长经验'
    if (selected.length === 0) {
      paymentHint = '没有可结算的商品'
    } else if (pointsUnavailable) {
      paymentHint = '所选商品中有暂不支持积分兑换的商品'
    } else if (!pointsEnough) {
      paymentHint = `积分兑换还差 ${Math.max(0, pointsTotal - Number(this.data.profile.points || 0))} 积分`
    } else if (paymentType === 'points') {
      paymentHint = `本次将消耗 ${pointsTotal} 积分，并获得 ${rewardExpTotal} 经验`
    } else {
      paymentHint = `现金购买预计获得 ${rewardPointsTotal} 积分 · 成长 +${rewardExpTotal} 经验`
    }
    this.setData({
      paymentType,
      paymentTypeText: paymentType === 'points' ? '积分兑换' : '现金购买',
      cashTotalText: `¥${(cashTotal / 100).toFixed(2)}`,
      pointsTotalText: `${pointsTotal} 积分`,
      totalText: paymentType === 'points' ? `${pointsTotal} 积分` : `¥${(cashTotal / 100).toFixed(2)}`,
      rewardPointsTotal,
      rewardExpTotal,
      paymentHint,
      canPayPoints,
      canSubmit: selected.length > 0 && (paymentType === 'cash' || canPayPoints)
    })
  },

  choosePayment(event) {
    const type = event.currentTarget.dataset.type
    if (type === 'points' && !this.data.canPayPoints) {
      wx.showToast({ title: this.data.paymentHint || '当前无法使用积分兑换', icon: 'none' })
      return
    }
    this.setData({ paymentType: type }, () => this.syncSummary())
  },

  chooseAddress() {
    const current = this.data.address
    wx.navigateTo({
      url: `/pages/address/address?select=1${current ? `&selectedId=${current.id}` : ''}`
    })
  },

  onRemark(event) {
    this.setData({ remark: event.detail.value })
  },

  async submit() {
    if (this.data.submitting) return
    if (!this.data.address) {
      wx.showToast({ title: '请先选择收货地址', icon: 'none' })
      return
    }
    if (this.data.items.length === 0) {
      wx.showToast({ title: '没有可结算的商品', icon: 'none' })
      return
    }
    if (!this.data.canSubmit) {
      wx.showToast({ title: this.data.paymentHint || '当前无法提交订单', icon: 'none' })
      return
    }
    const draft = getApp().globalData.checkoutDraft || {}
    this.setData({ submitting: true })
    try {
      const result = await api.request('/api/orders', 'POST', {
        addressId: this.data.address.id,
        remark: this.data.remark,
        paymentType: this.data.paymentType,
        items: this.data.items.map(item => ({
          productId: item.productId,
          quantity: item.quantity,
          customText: item.customText || ''
        })),
        cartItemIds: draft.cartItemIds || []
      })
      if (draft.cartItemIds && draft.cartItemIds.length) {
        getApp().globalData.cartCount = 0
      }
      getApp().globalData.checkoutDraft = null
      const order = result.order || {}
      const rewardText = this.data.paymentType === 'points'
        ? `积分兑换成功，成长 +${Number(order.rewardExp || this.data.rewardExpTotal)} 经验`
        : `下单成功，获得 ${Number(order.rewardPoints || this.data.rewardPointsTotal)} 积分 · +${Number(order.rewardExp || this.data.rewardExpTotal)} 经验`
      wx.showToast({ title: rewardText, icon: 'none', duration: 2200 })
      setTimeout(() => {
        wx.redirectTo({ url: '/pages/profile/profile' })
      }, 1300)
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  }
})
