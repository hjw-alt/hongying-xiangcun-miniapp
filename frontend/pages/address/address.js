const api = require('../../utils/request.js')

Page({
  data: {
    addresses: [],
    selectMode: false,
    selectedId: '',
    formVisible: false,
    editingId: '',
    form: { receiver: '', phone: '', region: '', detail: '', isDefault: false },
    submitting: false,
    loading: true
  },

  onLoad(options) {
    this.setData({
      selectMode: options.select === '1',
      selectedId: options.selectedId || ''
    })
  },

  onShow() {
    this.loadAddresses()
  },

  async loadAddresses() {
    this.setData({ loading: true })
    try {
      const items = await api.request('/api/addresses')
      this.setData({ addresses: items, loading: false })
    } catch (error) {
      this.setData({ loading: false })
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onPick(event) {
    if (!this.data.selectMode) return
    const id = event.currentTarget.dataset.id
    const address = this.data.addresses.find(item => item.id === id)
    if (!address) return
    getApp().globalData.selectedAddress = address
    wx.navigateBack()
  },

  openCreate() {
    this.setData({
      formVisible: true,
      editingId: '',
      form: { receiver: '', phone: '', region: '', detail: '', isDefault: this.data.addresses.length === 0 }
    })
  },

  openEdit(event) {
    const address = this.data.addresses.find(item => item.id === event.currentTarget.dataset.id)
    if (!address) return
    this.setData({
      formVisible: true,
      editingId: address.id,
      form: {
        receiver: address.receiver,
        phone: address.phone,
        region: address.region,
        detail: address.detail,
        isDefault: address.isDefault
      }
    })
  },

  closeForm() {
    this.setData({ formVisible: false })
  },

  noop() {},

  onField(event) {
    const key = event.currentTarget.dataset.key
    this.setData({ [`form.${key}`]: event.detail.value })
  },

  toggleDefault() {
    this.setData({ 'form.isDefault': !this.data.form.isDefault })
  },

  async saveAddress() {
    if (this.data.submitting) return
    const form = this.data.form
    if (!form.receiver.trim()) {
      wx.showToast({ title: '请填写收货人', icon: 'none' })
      return
    }
    if (!/^1[3-9]\d{9}$/.test(form.phone.trim())) {
      wx.showToast({ title: '请填写正确的手机号', icon: 'none' })
      return
    }
    if (form.detail.trim().length < 4) {
      wx.showToast({ title: '请填写详细地址', icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      const payload = {
        receiver: form.receiver,
        phone: form.phone,
        region: form.region,
        detail: form.detail,
        isDefault: form.isDefault
      }
      if (this.data.editingId) {
        await api.request(`/api/addresses/${this.data.editingId}`, 'PUT', payload)
      } else {
        await api.request('/api/addresses', 'POST', payload)
      }
      this.setData({ formVisible: false })
      wx.showToast({ title: '已保存', icon: 'success' })
      await this.loadAddresses()
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },

  async setDefault(event) {
    try {
      await api.request(`/api/addresses/${event.currentTarget.dataset.id}/default`, 'POST')
      await this.loadAddresses()
      wx.showToast({ title: '已设为默认', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  removeAddress(event) {
    const id = event.currentTarget.dataset.id
    wx.showModal({
      title: '删除地址',
      content: '确定要删除这条收货地址吗？',
      confirmColor: '#B32A2A',
      success: async result => {
        if (!result.confirm) return
        try {
          await api.request(`/api/addresses/${id}`, 'DELETE')
          await this.loadAddresses()
          wx.showToast({ title: '已删除', icon: 'none' })
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  }
})
