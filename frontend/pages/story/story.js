const api = require('../../utils/request.js')

Page({
  data: {
    categories: ['全部', '人物故事', '革命事件', '求学立志', '艺术人生', '建设故事', '风物故事'],
    category: '全部',
    keyword: '',
    stories: [],
    loading: true
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 3, hidden: false })
    }
    this.loadStories()
  },

  async loadStories() {
    this.setData({ loading: true })
    try {
      const query = `category=${encodeURIComponent(this.data.category)}&keyword=${encodeURIComponent(this.data.keyword)}`
      const result = await api.request(`/api/stories?${query}`)
      this.setData({
        stories: (result.items || []).map(api.decorateStory),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onCategory(event) {
    this.setData({ category: event.currentTarget.dataset.value })
    this.loadStories()
  },

  onKeyword(event) {
    this.setData({ keyword: event.detail.value })
  },

  onSearch() {
    this.loadStories()
  },

  onStory(event) {
    wx.navigateTo({ url: `/pages/story-detail/story-detail?id=${event.currentTarget.dataset.id}` })
  },

  onShareAppMessage() {
    return { title: '红映湘村 · 红迹数字故事', path: '/pages/story/story' }
  }
})
