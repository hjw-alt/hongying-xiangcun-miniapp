Component({
  data: {
    selected: 0,
    hidden: false,
    list: [
      { pagePath: '/pages/home/home', text: '首页', icon: '/assets/tabbar-red/home.png', activeIcon: '/assets/tabbar-red/home-active.png' },
      { pagePath: '/pages/mall/mall', text: '红文创', icon: '/assets/tabbar-red/mall.png', activeIcon: '/assets/tabbar-red/mall-active.png' },
      { pagePath: '/pages/guide/guide', text: '导览', icon: '/assets/tabbar-red/guide.png', activeIcon: '/assets/tabbar-red/guide-active.png' },
      { pagePath: '/pages/story/story', text: '红迹', icon: '/assets/tabbar-red/story.png', activeIcon: '/assets/tabbar-red/story-active.png' },
      { pagePath: '/pages/study/study', text: '研学', icon: '/assets/tabbar-red/study.png', activeIcon: '/assets/tabbar-red/study-active.png' },
      { pagePath: '/pages/farm/farm', text: '助农', icon: '/assets/tabbar-red/farm.png', activeIcon: '/assets/tabbar-red/farm-active.png' }
    ]
  },
  methods: {
    onTap(event) {
      const index = Number(event.currentTarget.dataset.index)
      const target = this.data.list[index]
      if (!target || index === this.data.selected) return
      wx.switchTab({ url: target.pagePath })
    }
  }
})
