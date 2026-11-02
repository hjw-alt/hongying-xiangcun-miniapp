const api = require('../../utils/request.js')

const emptyForm = {
  title: '',
  author: '',
  theme: '纸影人物',
  description: '',
  tagsText: ''
}

Page({
  data: {
    heroImage: api.absoluteUrl('/static/seed/shadow-workshop.png'),
    form: { ...emptyForm },
    themes: ['纸影人物', '湘潭印象', '戏台故事'],
    tempImage: '',
    storedImage: '',
    coverUrl: '',
    editingId: '',
    submitting: false
  },

  onShow() {
    const editingId = getApp().globalData.editingId || ''
    if (editingId && editingId !== this.data.editingId) {
      this.loadForEdit(editingId)
    }
  },

  onHide() {
    if (this.data.editingId && !getApp().globalData.editingId) {
      this.resetForm()
    }
  },

  async loadForEdit(id) {
    try {
      const response = await api.request(`/api/artworks/${id}`)
      const item = api.decorateArtwork(response.artwork)
      this.setData({
        editingId: id,
        form: {
          title: item.title,
          author: item.author,
          theme: item.theme,
          description: item.description,
          tagsText: item.tags.join(' ')
        },
        storedImage: item.imageUrl,
        coverUrl: item.coverUrl,
        tempImage: ''
      })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onFieldInput(event) {
    const key = event.currentTarget.dataset.key
    this.setData({ [`form.${key}`]: event.detail.value })
  },

  onThemeSelect(event) {
    this.setData({ 'form.theme': event.currentTarget.dataset.theme })
  },

  chooseImage() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      sizeType: ['compressed'],
      success: response => {
        const file = response.tempFiles[0]
        if (file.size > 8 * 1024 * 1024) {
          wx.showToast({ title: '图片不能超过 8MB', icon: 'none' })
          return
        }
        this.setData({ tempImage: file.tempFilePath, coverUrl: file.tempFilePath })
      }
    })
  },

  removeEditState() {
    getApp().globalData.editingId = ''
    this.resetForm()
  },

  resetForm() {
    this.setData({
      form: { ...emptyForm },
      tempImage: '',
      storedImage: '',
      coverUrl: '',
      editingId: ''
    })
  },

  validateBeforeSave(imageUrl) {
    const form = this.data.form
    if (!form.title.trim()) return '请填写作品名称'
    if (!form.author.trim()) return '请填写作者'
    if (!imageUrl) return '请上传作品图片'
    return ''
  },

  async save() {
    if (this.data.submitting) return
    let imageUrl = this.data.tempImage ? this.data.tempImage : this.data.storedImage
    const message = this.validateBeforeSave(imageUrl)
    if (message) {
      wx.showToast({ title: message, icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      if (this.data.tempImage) {
        imageUrl = await api.uploadImage(this.data.tempImage)
      }
      const tags = this.data.form.tagsText.split(/[,，、\s]+/).filter(Boolean)
      const payload = { ...this.data.form, tags, imageUrl }
      const saved = this.data.editingId
        ? await api.request(`/api/artworks/${this.data.editingId}`, 'PUT', payload)
        : await api.request('/api/artworks', 'POST', payload)
      getApp().globalData.editingId = ''
      this.resetForm()
      wx.showToast({ title: saved.title + ' 已保存', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  }
})
