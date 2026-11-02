const api = require('../../utils/request.js')

const emptyProfile = api.decorateProfile({
  points: 0,
  level: 2,
  levelName: '初见',
  experience: 0,
  currentLevelExp: 0,
  nextLevelExp: 30,
  progress: 0
})

Page({
  data: {
    id: '',
    product: null,
    profile: emptyProfile,
    favorited: false,
    cartCount: 0,
    customText: '',
    quantity: 1,
    paymentType: 'cash',
    paymentTypeText: '现金购买',
    subtotalText: '¥0.00',
    pointsTotalText: '0 积分',
    rewardPointsTotal: 0,
    rewardExpTotal: 0,
    rewardText: '',
    canPayPoints: false,
    paymentHint: '现金购买后可获得积分与成长经验',
    adding: false,
    loading: true
  },

  onLoad(options) {
    this.setData({ id: options.id || '' })
  },

  onShow() {
    if (this.data.id) this.loadDetail()
    this.loadCartCount()
  },

  async loadDetail() {
    this.setData({ loading: true })
    try {
      const [product, favorite, profile] = await Promise.all([
        api.request(`/api/products/${this.data.id}`),
        api.request(`/api/favorites/item/product/${this.data.id}`),
        api.request('/api/profile')
      ])
      this.setData({
        product: api.decorateProduct(product),
        favorited: favorite.favorited,
        profile: api.decorateProfile(profile),
        loading: false
      }, () => this.syncSummary())
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  syncSummary() {
    const product = this.data.product
    if (!product) return
    const quantity = Math.max(1, Number(this.data.quantity || 1))
    const subtotal = Number(product.priceCents || 0) * quantity
    const pointsTotal = Number(product.pointsPrice || 0) * quantity
    const rewardPointsTotal = Number(product.rewardPoints || 0) * quantity
    const rewardExpTotal = Number(product.rewardExp || 0) * quantity
    const supportsPoints = pointsTotal > 0
    const canPayPoints = supportsPoints && Number(this.data.profile.points || 0) >= pointsTotal
    const paymentType = this.data.paymentType === 'points' && canPayPoints ? 'points' : 'cash'
    let paymentHint = '现金购买后可获得积分与成长经验'
    if (!supportsPoints) {
      paymentHint = '该商品暂不支持积分兑换'
    } else if (!canPayPoints) {
      paymentHint = `积分兑换还差 ${Math.max(0, pointsTotal - Number(this.data.profile.points || 0))} 积分`
    } else if (paymentType === 'points') {
      paymentHint = `本次将消耗 ${pointsTotal} 积分，并获得 ${rewardExpTotal} 经验`
    } else {
      paymentHint = `现金购买预计获得 ${rewardPointsTotal} 积分 · ${rewardExpTotal} 经验`
    }
    this.setData({
      quantity,
      paymentType,
      paymentTypeText: paymentType === 'points' ? '积分兑换' : '现金购买',
      subtotalText: `¥${(subtotal / 100).toFixed(2)}`,
      pointsTotalText: `${pointsTotal} 积分`,
      rewardPointsTotal,
      rewardExpTotal,
      rewardText: rewardPointsTotal > 0 || rewardExpTotal > 0
        ? `现金购买预计 +${rewardPointsTotal} 积分 · 成长 +${rewardExpTotal} 经验`
        : '',
      canPayPoints,
      paymentHint
    })
  },

  async loadCartCount() {
    try {
      const result = await api.request('/api/cart/count')
      this.setData({ cartCount: result.count || 0 })
      getApp().globalData.cartCount = result.count || 0
    } catch (error) {
      this.setData({ cartCount: 0 })
    }
  },

  async onFavorite() {
    try {
      const result = await api.request(`/api/favorites/item/product/${this.data.id}/toggle`, 'POST')
      this.setData({ favorited: result.favorited })
      wx.showToast({ title: result.favorited ? '已收藏' : '已取消收藏', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onCustom(event) {
    this.setData({ customText: event.detail.value })
  },

  changeQuantity(event) {
    const delta = event.currentTarget.dataset.type === 'add' ? 1 : -1
    const next = Math.min(99, Math.max(1, this.data.quantity + delta))
    this.setData({ quantity: next }, () => this.syncSummary())
  },

  choosePayment(event) {
    const type = event.currentTarget.dataset.type
    if (type === 'points' && !this.data.canPayPoints) {
      wx.showToast({ title: this.data.paymentHint || '当前积分不足', icon: 'none' })
      return
    }
    this.setData({ paymentType: type }, () => this.syncSummary())
  },

  async addToCart() {
    if (this.data.adding || !this.data.product) return
    if (this.data.product.stock <= 0) {
      wx.showToast({ title: '该商品暂时缺货', icon: 'none' })
      return
    }
    this.setData({ adding: true })
    try {
      const result = await api.request('/api/cart', 'POST', {
        productId: this.data.product.id,
        quantity: this.data.quantity,
        customText: this.data.customText
      })
      this.setData({ cartCount: result.count || 0 })
      getApp().globalData.cartCount = result.count || 0
      wx.showToast({ title: '已加入购物车', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ adding: false })
    }
  },

  buyNow() {
    const product = this.data.product
    if (!product) return
    if (product.stock <= 0) {
      wx.showToast({ title: '该商品暂时缺货', icon: 'none' })
      return
    }
    if (this.data.paymentType === 'points' && !this.data.canPayPoints) {
      wx.showToast({ title: this.data.paymentHint || '当前积分不足', icon: 'none' })
      return
    }
    getApp().globalData.checkoutDraft = {
      source: 'buy',
      paymentType: this.data.paymentType,
      items: [{
        productId: product.id,
        quantity: this.data.quantity,
        customText: this.data.customText,
        name: product.name,
        priceCents: product.priceCents,
        pointsPrice: product.pointsPrice,
        rewardPoints: product.rewardPoints,
        rewardExp: product.rewardExp,
        coverSrc: product.coverSrc,
        farmerName: product.farmerName
      }],
      cartItemIds: []
    }
    wx.navigateTo({ url: '/pages/checkout/checkout' })
  },

  goCart() {
    wx.navigateTo({ url: '/pages/cart/cart' })
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' })
  },

  onShareAppMessage() {
    const product = this.data.product
    return {
      title: product ? `红映湘村 · ${product.name}` : '红映湘村',
      path: `/pages/product/product?id=${this.data.id}`
    }
  }
})
