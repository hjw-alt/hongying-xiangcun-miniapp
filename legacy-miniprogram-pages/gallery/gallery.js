const api = require('../../utils/request.js')

function splitColumns(items) {
  const leftColumn = []
  const rightColumn = []
  let leftWeight = 0
  let rightWeight = 0
  items.forEach(item => {
    const weight = item.tall ? 2 : 1
    if (leftWeight <= rightWeight) {
      leftColumn.push(item)
      leftWeight += weight
    } else {
      rightColumn.push(item)
      rightWeight += weight
    }
  })
  return { leftColumn, rightColumn }
}

Page({
  data: {
    heroImage: api.absoluteUrl('/static/seed/shadow-stage.png'),
    keyword: '',
    theme: '全部',
    sort: 'latest',
    themes: ['全部', '纸影人物', '湘潭印象', '戏台故事'],
    sorts: [
      { key: 'latest', label: '最新' },
      { key: 'likes', label: '最热' }
    ],
    stats: { artworkTotal: 0, likeTotal: 0, puppetTotal: 0, impressTotal: 0, storyTotal: 0 },
    artworks: [],
    leftColumn: [],
    rightColumn: [],
    loading: true,
    errorMessage: ''
  },

  onShow() {
    this.loadPage()
  },

  async loadPage() {
    this.setData({ loading: true, errorMessage: '' })
    try {
      const [stats, response] = await Promise.all([
        api.request('/api/stats'),
        this.loadArtworks()
      ])
      const artworks = response.items.map(api.decorateArtwork).map((item, index) => ({
        ...item,
        tall: index % 3 === 0
      }))
      const { leftColumn, rightColumn } = splitColumns(artworks)
      this.setData({
        stats,
        artworks,
        leftColumn,
        rightColumn,
        loading: false
      })
    } catch (error) {
      this.setData({
        loading: false,
        errorMessage: error.message || '服务暂时不可用'
      })
    }
  },

  async loadArtworks() {
    const query = [
      `keyword=${encodeURIComponent(this.data.keyword)}`,
      `theme=${encodeURIComponent(this.data.theme)}`,
      `sort=${this.data.sort}`
    ].join('&')
    return api.request(`/api/artworks?${query}`)
  },

  onKeywordInput(event) {
    this.setData({ keyword: event.detail.value })
  },

  onSearch() {
    this.loadPage()
  },

  onClearSearch() {
    this.setData({ keyword: '' })
    this.loadPage()
  },

  onThemeSelect(event) {
    this.setData({ theme: event.currentTarget.dataset.theme })
    this.loadPage()
  },

  onSortSelect(event) {
    this.setData({ sort: event.currentTarget.dataset.sort })
    this.loadPage()
  },

  onRefresh() {
    this.loadPage()
  },

  onArtworkTap(event) {
    const { id } = event.currentTarget.dataset
    wx.navigateTo({ url: `/pages/detail/detail?id=${id}` })
  },

  async onLikeTap(event) {
    const id = event.currentTarget.dataset.id
    try {
      const item = await api.request(`/api/artworks/${id}/like`, 'POST')
      const artworks = this.data.artworks.map(row => {
        if (row.id !== id) return row
        return api.decorateArtwork({ ...row, ...item })
      })
      const { leftColumn, rightColumn } = splitColumns(artworks)
      this.setData({ artworks, leftColumn, rightColumn })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  goCreate() {
    wx.switchTab({ url: '/pages/create/create' })
  },

  goManage() {
    wx.navigateTo({ url: '/pages/manage/manage' })
  }
})
