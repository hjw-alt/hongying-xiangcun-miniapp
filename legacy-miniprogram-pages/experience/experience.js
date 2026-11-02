const api = require('../../utils/request.js')

Page({
  data: {
    tabs: [
      { key: 'online', title: '在线影戏' },
      { key: 'ar', title: '影偶生成' },
      { key: 'diy', title: 'DIY 教程' },
      { key: 'events', title: '研学预约' }
    ],
    tabMeta: {
      online: {
        title: '影戏剧场',
        sub: '灯亮、锣响，戏就开了',
        image: api.absoluteUrl('/static/seed/shadow-stage.png')
      },
      ar: {
        title: '影偶生成',
        sub: '拍下自己，生成一段纸影戏',
        image: api.absoluteUrl('/static/seed/art-puppet.png')
      },
      diy: {
        title: 'DIY 教程',
        sub: '从画稿到装杆，做一只自己的影偶',
        image: api.absoluteUrl('/static/seed/shadow-workshop.png')
      },
      events: {
        title: '研学预约',
        sub: '走进戏台、课堂与展演现场',
        image: api.absoluteUrl('/static/seed/activity-night.png')
      }
    },
    headerTitle: '影戏剧场',
    headerSub: '灯亮、锣响，戏就开了',
    headerImage: api.absoluteUrl('/static/seed/shadow-stage.png'),
    activeTab: 'online',
    clips: [],
    activeClip: null,
    favState: {},
    cameraError: false,
    photoTemp: '',
    photoReady: false,
    arAuthor: '',
    playTitle: '',
    presets: ['亮相', '舞剑', '赶路'],
    selectedPreset: '亮相',
    modes: ['彩影', '墨影'],
    selectedMode: '彩影',
    generatingPlay: false,
    currentPlay: null,
    myPlays: [],
    diySteps: [
      { title: '定造型', text: '选穆桂英、孙悟空或地方人物侧影，用铅笔确定七分脸与动态线。' },
      { title: '拆纹样', text: '把头冠、肩甲、衣袖分区，用镂空花纹区分前后层次。' },
      { title: '上颜色', text: '红、黄、绿、黑为主色，先平涂大块，再用细笔描边。' },
      { title: '剪刻装杆', text: '剪下纸板后刻空纹样，装上头签和手签，灯光下试操纵。' }
    ],
    diyImage: '',
    diyAuthor: '',
    savingDIY: false,
    activities: [],
    bookingVisible: false,
    selectedActivity: null,
    bookingForm: { name: '', phone: '', participants: 1 },
    submittingBooking: false
  },

  onShow() {
    const tab = getApp().globalData.experienceTab || 'online'
    this.applyTab(tab)
    this.loadTab(tab)
  },

  onTabTap(event) {
    const tab = event.currentTarget.dataset.tab
    this.applyTab(tab)
    this.loadTab(tab)
  },

  applyTab(tab) {
    const meta = this.data.tabMeta[tab] || this.data.tabMeta.online
    this.setData({
      activeTab: tab,
      headerTitle: meta.title,
      headerSub: meta.sub,
      headerImage: meta.image
    })
  },

  noop() {},

  loadTab(tab) {
    if (tab === 'online') this.loadClips()
    if (tab === 'ar') this.loadMyPlays()
    if (tab === 'events') this.loadActivities()
  },

  async loadClips() {
    try {
      const clips = await api.request('/api/clips')
      const decorated = clips.map(api.decorateClip)
      const favState = {}
      for (const item of decorated) {
        const response = await api.request(`/api/favorites/clip/${item.id}`)
        favState[item.id] = response.favorited
      }
      this.setData({
        clips: decorated,
        favState,
        activeClip: this.data.activeClip || decorated[0] || null
      })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onPlayClip(event) {
    const item = this.data.clips.find(row => row.id === event.currentTarget.dataset.id)
    this.setData({ activeClip: item })
  },

  async likeClip(event) {
    const id = event.currentTarget.dataset.id
    try {
      const item = await api.request(`/api/clips/${id}/like`, 'POST')
      const clips = this.data.clips.map(row => row.id === id ? api.decorateClip({ ...row, ...item }) : row)
      this.setData({
        clips,
        activeClip: this.data.activeClip && this.data.activeClip.id === id ? api.decorateClip({ ...this.data.activeClip, ...item }) : this.data.activeClip
      })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async favoriteClip(event) {
    const id = event.currentTarget.dataset.id
    try {
      const response = await api.request(`/api/favorites/clip/${id}/toggle`, 'POST')
      this.setData({ [`favState.${id}`]: response.favorited })
      wx.showToast({ title: response.favorited ? '已收藏' : '已取消', icon: 'none' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  onCameraError() {
    this.setData({ cameraError: true })
  },

  onARAuthor(event) {
    this.setData({ arAuthor: event.detail.value })
  },

  takePhoto() {
    if (this.data.photoReady) {
      this.setData({ photoTemp: '', photoReady: false })
      if (this.data.cameraError) {
        this.choosePhoto()
      } else {
        wx.showToast({ title: '已重置，可以重拍', icon: 'none' })
      }
      return
    }
    if (this.data.cameraError || !wx.createCameraContext) {
      this.choosePhoto()
      return
    }
    let context
    try {
      context = wx.createCameraContext(this)
    } catch (error) {
      context = wx.createCameraContext()
    }
    context.takePhoto({
      quality: 'high',
      success: result => {
        if (!result || !result.tempImagePath) {
          this.setData({ cameraError: true })
          wx.showToast({ title: '没有拍到画面，改为选择照片', icon: 'none' })
          this.choosePhoto()
          return
        }
        this.setData({ photoTemp: result.tempImagePath, photoReady: true })
        wx.showToast({ title: '已拍下，可以保存作品', icon: 'none' })
      },
      fail: () => {
        this.setData({ cameraError: true })
        wx.showToast({ title: '相机不可用，改为选择照片', icon: 'none' })
        this.choosePhoto()
      }
    })
  },

  choosePhoto() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: result => {
        const file = result.tempFiles && result.tempFiles[0]
        if (!file || !file.tempFilePath) return
        this.setData({ photoTemp: file.tempFilePath, photoReady: true })
      }
    })
  },

  onPlayTitle(event) {
    this.setData({ playTitle: event.detail.value })
  },

  onPresetSelect(event) {
    this.setData({ selectedPreset: event.currentTarget.dataset.preset })
  },

  onModeSelect(event) {
    this.setData({ selectedMode: event.currentTarget.dataset.mode })
  },

  async loadMyPlays() {
    try {
      const items = await api.request('/api/shadow-plays')
      this.setData({
        myPlays: items.map(item => ({
          ...item,
          gifSrc: api.absoluteUrl(item.gifUrl),
          posterSrc: api.absoluteUrl(item.posterUrl),
          createdAtText: api.formatDateTime(item.createdAt)
        }))
      })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  async generatePlay() {
    if (this.data.generatingPlay) return
    if (!this.data.photoTemp) {
      wx.showToast({ title: '请先拍摄或选择照片', icon: 'none' })
      return
    }
    this.setData({ generatingPlay: true })
    wx.showLoading({ title: '正在生成纸影戏', mask: true })
    try {
      const sourceUrl = await api.uploadImage(this.data.photoTemp)
      const play = await api.request('/api/shadow-plays', 'POST', {
        title: this.data.playTitle.trim() || '我的纸影戏',
        author: this.data.arAuthor.trim() || '匿名同学',
        preset: this.data.selectedPreset,
        mode: this.data.selectedMode,
        sourceUrl
      })
      this.setData({
        currentPlay: {
          ...play,
          gifSrc: api.absoluteUrl(play.gifUrl),
          posterSrc: api.absoluteUrl(play.posterUrl),
          createdAtText: api.formatDateTime(play.createdAt)
        }
      })
      await this.loadMyPlays()
      wx.showToast({ title: '纸影戏生成完成', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message || '生成失败', icon: 'none' })
    } finally {
      wx.hideLoading()
      this.setData({ generatingPlay: false })
    }
  },

  selectPlay(event) {
    const item = this.data.myPlays.find(row => row.id === event.currentTarget.dataset.id)
    if (!item) return
    this.setData({ currentPlay: item })
  },

  deletePlay(event) {
    const id = event.currentTarget.dataset.id
    wx.showModal({
      title: '删除影戏',
      content: '确定要删除这段生成的纸影戏吗？',
      confirmColor: '#A93628',
      success: async result => {
        if (!result.confirm) return
        try {
          await api.request(`/api/shadow-plays/${id}`, 'DELETE')
          if (this.data.currentPlay && this.data.currentPlay.id === id) {
            this.setData({ currentPlay: null })
          }
          await this.loadMyPlays()
          wx.showToast({ title: '已删除', icon: 'none' })
        } catch (error) {
          wx.showToast({ title: error.message, icon: 'none' })
        }
      }
    })
  },

  onDIYAuthor(event) {
    this.setData({ diyAuthor: event.detail.value })
  },

  chooseDIYImage() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: result => this.setData({ diyImage: result.tempFiles[0].tempFilePath })
    })
  },

  async saveDIYWork() {
    if (this.data.savingDIY) return
    if (!this.data.diyImage) {
      wx.showToast({ title: '请先上传作品图', icon: 'none' })
      return
    }
    this.setData({ savingDIY: true })
    try {
      const imageUrl = await api.uploadImage(this.data.diyImage)
      await api.request('/api/artworks', 'POST', {
        title: '简易纸影人DIY',
        author: this.data.diyAuthor || '匿名同学',
        theme: '纸影人物',
        description: '按分步教程完成的纸影人练习，记录造型与镂空纹样尝试。',
        tags: ['DIY', '简易纸影人'],
        imageUrl
      })
      this.setData({ diyImage: '' })
      wx.showToast({ title: 'DIY作品已保存', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ savingDIY: false })
    }
  },

  async loadActivities() {
    try {
      const activities = await api.request('/api/activities')
      const decorated = activities.map(api.decorateActivity)
      this.setData({ activities: decorated })
      const pendingBooking = getApp().globalData.bookingActivityId
      if (pendingBooking) {
        getApp().globalData.bookingActivityId = ''
        const activity = decorated.find(row => row.id === pendingBooking)
        if (activity && !activity.bookingDisabled) {
          this.setData({
            bookingVisible: true,
            selectedActivity: activity,
            bookingForm: { name: '', phone: '', participants: 1 }
          })
        }
      }
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    }
  },

  openActivityDetail(event) {
    wx.navigateTo({ url: `/pages/activity/activity?id=${event.currentTarget.dataset.id}` })
  },

  openBooking(event) {
    const activity = this.data.activities.find(row => row.id === event.currentTarget.dataset.id)
    if (!activity || activity.bookingDisabled) return
    this.setData({
      bookingVisible: true,
      selectedActivity: activity,
      bookingForm: { name: '', phone: '', participants: 1 }
    })
  },

  closeBooking() {
    this.setData({ bookingVisible: false, selectedActivity: null })
  },

  onBookingField(event) {
    const key = event.currentTarget.dataset.key
    this.setData({ [`bookingForm.${key}`]: event.detail.value })
  },

  changeParticipants(event) {
    const type = event.currentTarget.dataset.type
    let value = this.data.bookingForm.participants + (type === 'add' ? 1 : -1)
    value = Math.min(10, Math.max(1, value))
    this.setData({ 'bookingForm.participants': value })
  },

  async submitBooking() {
    if (this.data.submittingBooking || !this.data.selectedActivity) return
    this.setData({ submittingBooking: true })
    try {
      await api.request('/api/bookings', 'POST', {
        activityId: this.data.selectedActivity.id,
        ...this.data.bookingForm,
        participants: Number(this.data.bookingForm.participants)
      })
      this.closeBooking()
      await this.loadActivities()
      wx.showToast({ title: '预约成功', icon: 'success' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submittingBooking: false })
    }
  }
})
