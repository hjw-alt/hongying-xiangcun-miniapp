const api = require('../../utils/request.js')
const fallbackProducts = [
  { id: 'product-bag-poem', name: '毛泽东诗词书法帆布袋', category: '红色经典IP', priceCents: 5800, originalCents: 7800, stock: 200, sales: 126, coverUrl: '/assets/red/product-bag.jpg', tags: ['帆布袋', '诗词', '日常可用'] },
  { id: 'product-badge-peng', name: '彭德怀“横刀立马”金属徽章', category: '红色经典IP', priceCents: 2900, originalCents: 3900, stock: 300, sales: 208, coverUrl: '/assets/red/product-badge.jpg', tags: ['徽章', '元帅故里', '收藏'] },
  { id: 'product-sticker-guju', name: '韶山故居线稿贴纸套装', category: '红色经典IP', priceCents: 1800, originalCents: 2600, stock: 500, sales: 342, coverUrl: '/assets/red/product-sticker.jpg', tags: ['贴纸', '手绘', '手账'] },
  { id: 'product-blindbox', name: '红军长征主题盲盒手办', category: '红色经典IP', priceCents: 6900, originalCents: 8900, stock: 150, sales: 87, coverUrl: '/assets/red/product-blindbox.jpg', tags: ['盲盒', '手办', '研学奖品'] },
  { id: 'product-scarf-shrimp', name: '齐白石虾图真丝丝巾', category: '齐白石艺术', priceCents: 29800, originalCents: 36800, stock: 60, sales: 34, coverUrl: '/assets/red/product-scarf.jpg', tags: ['真丝', '齐白石', '礼品'] },
  { id: 'product-seal', name: '齐白石篆刻印章文创礼盒', category: '齐白石艺术', priceCents: 39800, originalCents: 49800, stock: 40, sales: 21, coverUrl: '/assets/red/product-seal.jpg', tags: ['篆刻', '定制', '礼盒'] },
  { id: 'product-lotus-gift', name: '湘莲红色主题包装礼盒', category: '湘村风物助农', priceCents: 8800, originalCents: 10800, stock: 300, sales: 178, coverUrl: '/assets/red/product-lotus.jpg', tags: ['助农', '湘莲', '产地直发'] },
  { id: 'product-camellia-oil', name: '湘潭茶油红色文化伴手礼', category: '湘村风物助农', priceCents: 19800, originalCents: 23800, stock: 150, sales: 96, coverUrl: '/assets/red/product-camellia.jpg', tags: ['助农', '茶油', '伴手礼'] },
  { id: 'product-bamboo-craft', name: '乌石峰“三十六坊”竹编手作', category: '湘村风物助农', priceCents: 12800, originalCents: 15800, stock: 120, sales: 65, coverUrl: '/assets/red/product-bamboo.jpg', tags: ['助农', '非遗手作', '竹编'] },
  { id: 'product-chili-sauce', name: '泉湖剁辣椒与腊味伴手礼', category: '湘村风物助农', priceCents: 13800, originalCents: 16800, stock: 100, sales: 58, coverUrl: '/assets/red/product-chili.jpg', tags: ['助农', '湘味', '组合装'] },
  { id: 'product-digital-card', name: '红色村落数字纪念卡', category: '数字文创', priceCents: 3900, originalCents: 4900, stock: 999, sales: 412, coverUrl: '/assets/red/product-digital.jpg', tags: ['数字文创', '纪念卡', '收藏'] }
]

function getFallbackProducts(category, keyword, sort) {
  const normalizedKeyword = (keyword || '').trim().toLowerCase()
  let items = fallbackProducts.filter(item => {
    const categoryMatched = !category || category === '全部' || item.category === category
    const keywordMatched = !normalizedKeyword || `${item.name}${item.category}${item.tags.join('')}`.toLowerCase().includes(normalizedKeyword)
    return categoryMatched && keywordMatched
  })
  if (sort === 'sales') items.sort((a, b) => b.sales - a.sales)
  if (sort === 'price') items.sort((a, b) => a.priceCents - b.priceCents)
  return items.map(item => api.decorateProduct({
    ...item,
    pointsPrice: item.pointsPrice || Math.max(100, Math.round(item.priceCents / 10)),
    rewardPoints: item.rewardPoints || Math.max(1, Math.round(item.priceCents / 1000)),
    rewardExp: item.rewardExp || Math.max(5, Math.round(item.priceCents / 800))
  }))
}
Page({
  data: {
    categories: ['全部', '红色经典IP', '齐白石艺术', '湘村风物助农', '数字文创'],
    category: '全部',
    sorts: [
      { key: 'default', label: '综合' },
      { key: 'sales', label: '销量' },
      { key: 'price', label: '价格' }
    ],
    sort: 'default',
    keyword: '',
    products: getFallbackProducts('全部', '', 'default'),
    cartCount: 0,
    profile: api.decorateProfile({ points: 0, level: 2, levelName: '初见', experience: 0, nextLevelExp: 30, progress: 0 }),
    addingId: '',
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 1, hidden: false })
    }
    this.loadProducts()
    this.loadCartCount()
    this.loadProfile()
  },

  async loadProducts() {
    this.setData({ loading: true })
    try {
      const query = `category=${encodeURIComponent(this.data.category)}&sort=${this.data.sort}&keyword=${encodeURIComponent(this.data.keyword)}`
      const result = await api.request(`/api/products?${query}`)
      this.setData({
        products: (result.items || []).map(api.decorateProduct),
        loading: false
      })
    } catch (error) {
      this.setData({
        products: getFallbackProducts(this.data.category, this.data.keyword, this.data.sort),
        loading: false
      })
      console.warn('商城接口暂不可用，已使用本地示例商品', error)
    }
  },


  async loadProfile() {
    try {
      const profile = await api.request('/api/profile')
      this.setData({ profile: api.decorateProfile(profile) })
    } catch (error) {
      // 保留默认成长账户，商城仍可浏览。
    }
  },

  async loadCartCount() {
    try {
      const result = await api.request('/api/cart/count')
      this.setData({ cartCount: result.count || 0 })
      getApp().globalData.cartCount = result.count || 0
    } catch (error) {
      this.setData({ cartCount: getApp().globalData.cartCount || 0 })
    }
  },

  onCategory(event) {
    this.setData({ category: event.currentTarget.dataset.value })
    this.loadProducts()
  },

  onSort(event) {
    this.setData({ sort: event.currentTarget.dataset.value })
    this.loadProducts()
  },

  onKeyword(event) {
    this.setData({ keyword: event.detail.value })
  },

  onSearch() {
    this.loadProducts()
  },

  async addToCart(event) {
    const productId = event.currentTarget.dataset.id
    const product = this.data.products.find(item => item.id === productId)
    if (!product || product.stock <= 0 || this.data.addingId) return
    this.setData({ addingId: productId })
    try {
      const result = await api.request('/api/cart', 'POST', {
        productId,
        quantity: 1,
        customText: ''
      })
      const count = result.count || 0
      this.setData({ cartCount: count })
      getApp().globalData.cartCount = count
      wx.showToast({ title: '已加入购物车', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ addingId: '' })
    }
  },

  onProduct(event) {
    wx.navigateTo({ url: `/pages/product/product?id=${event.currentTarget.dataset.id}` })
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' })
  },

  goCart() {
    wx.navigateTo({ url: '/pages/cart/cart' })
  },

  goFarm() {
    wx.switchTab({ url: '/pages/farm/farm' })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 红文创商城', path: '/pages/mall/mall' }
  }
})


