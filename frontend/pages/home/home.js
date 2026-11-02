const api = require('../../utils/request.js')
const fallbackBanners = [
  { id: 'banner-village', title: '红映湘村', subtitle: '湘潭红色村落数字文旅与文创助农', imageUrl: '/assets/red/banner-village.jpg', linkType: 'village', linkId: 'village-shaoshan' },
  { id: 'banner-market', title: '红文创商城', subtitle: '文创好物与助农特产，每一单都分给村里', imageUrl: '/assets/red/banner-market.jpg', linkType: 'tab', linkId: 'mall' },
  { id: 'banner-story', title: '红迹数字故事', subtitle: '扫码读故事，听湘潭红色村落的记忆', imageUrl: '/assets/red/banner-story.jpg', linkType: 'tab', linkId: 'story' }
]

const fallbackVillages = [
  { id: 'village-shaoshan', name: '韶山村', region: '韶山市韶山乡', coverUrl: '/assets/red/village-shaoshan.jpg', spotCount: 3, storyCount: 2 },
  { id: 'village-wushi', name: '乌石村', region: '湘潭县乌石镇', coverUrl: '/assets/red/village-wushi.jpg', spotCount: 3, storyCount: 1 },
  { id: 'village-quanhu', name: '泉湖村', region: '湘乡市龙洞镇', coverUrl: '/assets/red/village-quanhu.jpg', spotCount: 2, storyCount: 1 },
  { id: 'village-xingtang', name: '杏花村', region: '湘潭县白石镇', coverUrl: '/assets/red/village-xingtang.jpg', spotCount: 2, storyCount: 2 }
]

const fallbackRoutes = [
  { id: 'route-shaoshan', name: '伟人故里红色研学线', theme: '伟人故里', summary: '故居、纪念馆与乡村记忆串联的一日研学路线', coverUrl: '/assets/red/route-shaoshan.jpg', days: 1, distanceKm: 18 },
  { id: 'route-wushi', name: '横刀立马英雄线', theme: '英雄故事', summary: '走进乌石峰下，读懂一代元帅的家国担当', coverUrl: '/assets/red/route-wushi.jpg', days: 1, distanceKm: 12 },
  { id: 'route-art', name: '齐白石艺术寻踪线', theme: '艺术研学', summary: '从星斗塘出发，寻找齐白石笔下的湘潭乡土', coverUrl: '/assets/red/route-art.jpg', days: 1, distanceKm: 16 }
]

const fallbackStories = [
  { id: 'story-baishi', title: '从木匠到大师：齐白石的乡土笔墨', category: '艺术记忆', era: '近现代', summary: '在星斗塘与老屋之间，读懂齐白石一生不教一日闲过的坚持。', coverUrl: '/assets/red/story-baishi.jpg', views: 1280 },
  { id: 'story-hengdao', title: '谁敢横刀立马，唯我彭大将军', category: '英雄故事', era: '革命年代', summary: '乌石峰下的少年走向战场，留下属于湘潭的英雄回响。', coverUrl: '/assets/red/story-hengdao.jpg', views: 968 },
  { id: 'story-xianglian', title: '一朵湘莲的红色乡愁', category: '乡村记忆', era: '当代', summary: '从荷塘到礼盒，湘莲产业把村民的日子越过越甜。', coverUrl: '/assets/red/story-xianglian.jpg', views: 736 },
  { id: 'story-guanqu', title: '南方红旗渠：韶山灌区记忆', category: '建设记忆', era: '上世纪六十年代', summary: '一渠清水穿过湘中丘陵，也串起一代人的奋斗故事。', coverUrl: '/assets/red/story-guanqu.jpg', views: 615 }
]

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

const fallbackFarmers = [
  { id: 'farmer-xianglian', name: '杏花村湘莲合作社', project: '湘莲红色主题礼盒与通心白莲加工', avatarUrl: '/assets/red/farmer-lotus.jpg', supportTarget: 1200, supportCount: 386 },
  { id: 'farmer-bamboo', name: '乌石竹编工坊', project: '竹编包装与竹编手作课程', avatarUrl: '/assets/red/farmer-bamboo.jpg', supportTarget: 800, supportCount: 245 },
  { id: 'farmer-tea', name: '银田茶油作坊', project: '茶油红色文化伴手礼', avatarUrl: '/assets/red/farmer-tea.jpg', supportTarget: 1000, supportCount: 418 }
]

function decorateBanner(item) {
  return { ...item, imageSrc: item.imageSrc || api.absoluteUrl(item.imageUrl) }
}

function decorateProduct(item) {
  const decorated = api.decorateProduct(item)
  return { ...decorated, coverSrc: item.coverSrc || decorated.coverSrc }
}

Page({
  data: {
    banners: [],
    villages: [],
    routes: [],
    stories: [],
    products: [],
    farmers: [],
    farmSummary: { farmers: 0, products: 0, supportCount: 0 },
    profile: api.decorateProfile({ nickname: '湘村体验者', level: 2, levelName: '初见', points: 0, experience: 0, nextLevelExp: 30, progress: 0, coupons: 0, balanceCents: 0 }),
    medals: [],
    unlockedCount: 0,
    cartCount: 0,
    serviceItems: [
      { key: 'guide', title: '红色导览', sub: '六村线路', icon: '导' },
      { key: 'mall', title: '文创商城', sub: '湘村好物', icon: '创' },
      { key: 'story', title: '红迹故事', sub: '扫码读史', icon: '迹' },
      { key: 'study', title: '研学打卡', sub: '任务集章', icon: '学' }
    ],
    benefits: [
      { key: 'farm', title: '助农好物', sub: '每单分给村里', icon: '农' },
      { key: 'guide', title: '村落导览', sub: '精选主题线路', icon: '路' },
      { key: 'story', title: '数字故事', sub: '听见湘村记忆', icon: '听' },
      { key: 'study', title: '研学活动', sub: '打卡兑徽章', icon: '章' }
    ],
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 0, hidden: false })
    }
    this.loadData()
    this.loadCartCount()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      const [banners, villageResult, routes, storyResult, productResult, farmers, farmSummary, profileResult, medalResult] = await Promise.all([
        api.request('/api/banners'),
        api.request('/api/villages'),
        api.request('/api/routes'),
        api.request('/api/stories'),
        api.request('/api/products?sort=sales'),
        api.request('/api/farmers'),
        api.request('/api/farm-summary'),
        api.request('/api/profile'),
        api.request('/api/medals')
      ])
      this.setData({
        banners: banners.map(decorateBanner),
        villages: (villageResult.items || []).slice(0, 4).map(api.decorateVillage),
        routes: routes.slice(0, 3).map(api.decorateRoute),
        stories: (storyResult.items || []).slice(0, 4).map(api.decorateStory),
        products: (productResult.items || []).slice(0, 6).map(decorateProduct),
        farmers: farmers.slice(0, 3).map(api.decorateFarmer),
        farmSummary: farmSummary || { farmers: 0, products: 0, supportCount: 0 },
        profile: api.decorateProfile(profileResult),
        medals: (medalResult.items || []).map(api.decorateMedal).slice(0, 4),
        unlockedCount: (medalResult.items || []).filter(item => item.unlocked).length,
        loading: false
      })
    } catch (error) {
      this.setData({
        banners: fallbackBanners.map(decorateBanner),
        villages: fallbackVillages.slice(0, 4).map(api.decorateVillage),
        routes: fallbackRoutes.map(api.decorateRoute),
        stories: fallbackStories.map(api.decorateStory),
        products: fallbackProducts.slice(0, 6).map(decorateProduct),
        farmers: fallbackFarmers.map(api.decorateFarmer),
        farmSummary: { farmers: 32, products: 8, supportCount: 386 },
        profile: api.decorateProfile({ nickname: '湘村体验者', level: 2, levelName: '初见', points: 0, experience: 0, nextLevelExp: 30, progress: 0, coupons: 0, balanceCents: 0 }),
        medals: [],
        unlockedCount: 0,
        loading: false
      })
      console.warn('首页接口暂不可用，已使用本地示例数据', error)
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

  onQuickTap(event) {
    const key = event.currentTarget.dataset.key
    const map = {
      guide: '/pages/guide/guide',
      mall: '/pages/mall/mall',
      story: '/pages/story/story',
      study: '/pages/study/study',
      farm: '/pages/farm/farm'
    }
    if (map[key]) wx.switchTab({ url: map[key] })
  },

  onBenefitTap(event) {
    const key = event.currentTarget.dataset.key
    if (key === 'mall') {
      this.goMall()
      return
    }
    this.onQuickTap(event)
  },

  onBannerTap(event) {
    const banner = this.data.banners.find(item => item.id === event.currentTarget.dataset.id)
    if (!banner) return
    if (banner.linkType === 'village') {
      wx.navigateTo({ url: `/pages/village/village?id=${banner.linkId}` })
      return
    }
    const map = { mall: '/pages/mall/mall', story: '/pages/story/story', guide: '/pages/guide/guide', study: '/pages/study/study' }
    if (map[banner.linkId]) wx.switchTab({ url: map[banner.linkId] })
  },

  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' })
  },

  goMedals() {
    wx.navigateTo({ url: '/pages/medals/medals' })
  },

  onMedalTap(event) {
    const id = event.currentTarget.dataset.id
    if (id) wx.navigateTo({ url: `/pages/medal-detail/medal-detail?id=${id}` })
  },

  goCart() {
    wx.navigateTo({ url: '/pages/cart/cart' })
  },

  goGuide() {
    wx.switchTab({ url: '/pages/guide/guide' })
  },

  goVillages() {
    wx.switchTab({ url: '/pages/guide/guide' })
  },

  goStories() {
    wx.switchTab({ url: '/pages/story/story' })
  },

  goMall() {
    wx.switchTab({ url: '/pages/mall/mall' })
  },

  goFarm() {
    wx.switchTab({ url: '/pages/farm/farm' })
  },

  goStudy() {
    wx.switchTab({ url: '/pages/study/study' })
  },

  onVillage(event) {
    wx.navigateTo({ url: `/pages/village/village?id=${event.currentTarget.dataset.id}` })
  },

  onStory(event) {
    wx.navigateTo({ url: `/pages/story-detail/story-detail?id=${event.currentTarget.dataset.id}` })
  },

  onProduct(event) {
    wx.navigateTo({ url: `/pages/product/product?id=${event.currentTarget.dataset.id}` })
  },

  onRoute(event) {
    wx.switchTab({ url: '/pages/guide/guide' })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 湘潭红色村落数字文旅', path: '/pages/home/home' }
  }
})


