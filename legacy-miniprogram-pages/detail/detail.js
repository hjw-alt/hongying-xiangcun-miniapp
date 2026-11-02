const api = require('../../utils/request.js')

Page({
  data: {
    id: '',
    artwork: null,
    comments: [],
    favorited: false,
    nickname: '',
    content: '',
    loading: true,
    submittingComment: false
  },

  onLoad(options) {
    this.setData({ id: options.id || '' })
  },

  onShow() {
    if (this.data.id) {
      this.loadDetail()
    }
  },

  async loadDetail() {
    try {
      const response = await api.request(`/api/artworks/${this.data.id}`)
      const favorite = await api.request(`/api/favorites/artwork/${this.data.id}`)
      this.setData({
        artwork: api.decorateArtwork(response.artwork),
        favorited: favorite.favorited,
        comments: response.comments.map(item => ({
          ...item,
          createdAtText: api.formatDateTime(item.createdAt)
        })),
        loading: false
      })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async onLike() {
    try {
      const item = await api.request(`/api/artworks/${this.data.id}/like`, 'POST')
      this.setData({ artwork: api.decorateArtwork({ ...this.data.artwork, ...item }) })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async onFavorite() {
    try {
      const response = await api.request(`/api/favorites/artwork/${this.data.id}/toggle`, 'POST')
      this.setData({ favorited: response.favorited })
      wx.showToast({ title: response.favorited ? '已收藏' : '已取消收藏', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onNicknameInput(event) {
    this.setData({ nickname: event.detail.value })
  },

  onContentInput(event) {
    this.setData({ content: event.detail.value })
  },

  async submitComment() {
    if (this.data.submittingComment) return
    const content = this.data.content.trim()
    if (!content) {
      wx.showToast({ title: '请填写评论', icon: 'none' })
      return
    }
    this.setData({ submittingComment: true })
    try {
      await api.request(`/api/artworks/${this.data.id}/comments`, 'POST', {
        nickname: this.data.nickname,
        content
      })
      this.setData({ content: '' })
      await this.loadDetail()
      wx.showToast({ title: '评论成功', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submittingComment: false })
    }
  },

  onDeleteComment(event) {
    const { id } = event.currentTarget.dataset
    wx.showModal({
      title: '删除评论',
      content: '确定要删除这条评论吗？',
      confirmColor: '#A93628',
      success: async result => {
        if (!result.confirm) return
        try {
          await api.request(`/api/artworks/${this.data.id}/comments/${id}`, 'DELETE')
          this.loadDetail()
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  },

  onEdit() {
    getApp().globalData.editingId = this.data.id
    wx.switchTab({ url: '/pages/create/create' })
  },

  onDelete() {
    wx.showModal({
      title: '删除作品',
      content: `确定要删除《${this.data.artwork.title}》吗？`,
      confirmColor: '#A93628',
      success: async result => {
        if (!result.confirm) return
        try {
          await api.request(`/api/artworks/${this.data.id}`, 'DELETE')
          wx.navigateBack()
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  },

  onShareAppMessage() {
    const artwork = this.data.artwork
    return {
      title: artwork ? `影绘湘潭 · ${artwork.title}` : '影绘湘潭',
      path: `/pages/detail/detail?id=${this.data.id}`
    }
  }
})
