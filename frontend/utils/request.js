const config = require('../config.js')

function absoluteUrl(path) {
  if (!path) return ''
  if (/^https?:\/\//.test(path) || /^\/assets\//.test(path) || /^wxfile:\/\//.test(path)) return path
  return config.BASE_URL + path
}

function getUserKey() {
  let key = wx.getStorageSync('userKey')
  if (!key) {
    key = `u_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
    wx.setStorageSync('userKey', key)
  }
  return key
}

function formatDateTime(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const pad = number => String(number).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function decorateArtwork(item) {
  const themeClass = item.theme === '纸影人物' ? 'puppet' : item.theme === '湘潭印象' ? 'impress' : 'story'
  return {
    ...item,
    coverUrl: absoluteUrl(item.imageUrl),
    themeClass,
    createdAtText: formatDateTime(item.createdAt),
    updatedAtText: formatDateTime(item.updatedAt)
  }
}

function decorateClip(item) {
  return {
    ...item,
    videoSrc: absoluteUrl(item.videoUrl),
    posterUrl: absoluteUrl(item.coverUrl),
    createdAtText: formatDateTime(item.createdAt)
  }
}

function decorateActivity(item) {
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    priceText: item.priceCents === 0 ? '免费' : `¥${(item.priceCents / 100).toFixed(0)}`,
    startTimeText: formatDateTime(item.startTime),
    endTimeText: formatDateTime(item.endTime),
    bookingDisabled: item.remaining <= 0
  }
}

function decorateBooking(item) {
  return {
    ...item,
    startTimeText: formatDateTime(item.startTime),
    createdAtText: formatDateTime(item.createdAt),
    feeText: item.totalFeeCents === 0 ? '免费' : `¥${(item.totalFeeCents / 100).toFixed(2)}`
  }
}

function decorateVillage(item) {
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    honorText: (item.honors || []).join(' · ')
  }
}

function decorateSpot(item) {
  return {
    ...item,
    imageSrc: absoluteUrl(item.imageUrl),
    audioSrc: absoluteUrl(item.audioUrl)
  }
}

function decorateRoute(item) {
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    metaText: `${item.days} 天 · 约 ${item.distanceKm} 公里`
  }
}

function decorateStory(item) {
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    audioSrc: absoluteUrl(item.audioUrl),
    arSrc: absoluteUrl(item.arImage),
    createdAtText: formatDateTime(item.createdAt),
    viewText: item.views >= 1000 ? `${(item.views / 1000).toFixed(1)}k` : String(item.views)
  }
}

function decorateProduct(item = {}) {
  const images = (item.images || []).map(absoluteUrl)
  const pointsPrice = Number(item.pointsPrice || 0)
  const rewardPoints = Number(item.rewardPoints || 0)
  const rewardExp = Number(item.rewardExp || 0)
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    imageList: images.length ? images : [absoluteUrl(item.coverUrl)],
    priceText: `¥${(Number(item.priceCents || 0) / 100).toFixed(2)}`,
    originalText: Number(item.originalCents || 0) > Number(item.priceCents || 0) ? `¥${(Number(item.originalCents || 0) / 100).toFixed(2)}` : '',
    pointsPrice,
    pointsPriceText: pointsPrice > 0 ? `${pointsPrice} 积分` : '暂不支持积分兑换',
    rewardPoints,
    rewardPointsText: rewardPoints > 0 ? `购买 +${rewardPoints} 积分` : '',
    rewardExp,
    rewardExpText: rewardExp > 0 ? `成长 +${rewardExp} 经验` : '',
    salesText: Number(item.sales || 0) >= 1000 ? `${(Number(item.sales) / 1000).toFixed(1)}k` : String(Number(item.sales || 0))
  }
}

function decorateFarmer(item) {
  return {
    ...item,
    avatarSrc: absoluteUrl(item.avatarUrl),
    progress: item.supportTarget > 0 ? Math.min(100, Math.round((item.supportCount / item.supportTarget) * 100)) : 0
  }
}

function decorateTask(item = {}) {
  const points = Number(item.points || 0)
  const rewardExp = Number(item.rewardExp || points)
  return {
    ...item,
    coverSrc: absoluteUrl(item.coverUrl),
    pointText: `+${points} 积分`,
    expText: `+${rewardExp} 经验`,
    rewardExp
  }
}

function decorateOrder(item = {}) {
  const paymentType = item.paymentType || 'cash'
  const pointsSpent = Number(item.pointsSpent || 0)
  return {
    ...item,
    paymentType,
    paymentTypeText: paymentType === 'points' ? '积分兑换' : '现金购买',
    totalText: paymentType === 'points' ? `${pointsSpent} 积分` : `¥${(Number(item.totalCents || 0) / 100).toFixed(2)}`,
    rewardText: Number(item.rewardPoints || 0) > 0 || Number(item.rewardExp || 0) > 0
      ? `获得 ${Number(item.rewardPoints || 0)} 积分 · ${Number(item.rewardExp || 0)} 经验`
      : '',
    createdAtText: formatDateTime(item.createdAt),
    items: (item.items || []).map(row => ({
      ...row,
      priceText: paymentType === 'points' && Number(row.pointsPrice || 0) > 0
        ? `${Number(row.pointsPrice)} 积分`
        : `¥${(Number(row.priceCents || 0) / 100).toFixed(2)}`
    }))
  }
}


function decorateProfile(profile = {}) {
  const experience = Number(profile.experience || 0)
  const nextLevelExp = Number(profile.nextLevelExp || 0)
  const currentLevelExp = Number(profile.currentLevelExp || 0)
  const progress = Math.max(0, Math.min(100, Number(profile.progress || 0)))
  return {
    ...profile,
    levelText: `V${profile.level || 2}`,
    pointsText: `${Number(profile.points || 0)} 积分`,
    experienceText: `${experience} 经验`,
    progressStyle: `width:${progress}%`,
    progressText: nextLevelExp > currentLevelExp ? `${experience} / ${nextLevelExp}` : `${experience} 经验`,
    balanceText: `¥${(Number(profile.balanceCents || 0) / 100).toFixed(2)}`,
    nextLevelText: nextLevelExp > experience ? `距离下一等级还需 ${nextLevelExp - experience} 经验` : '已达到最高等级'
  }
}

function decorateLedger(item = {}) {
  const change = Number(item.changeAmount || 0)
  return {
    ...item,
    changeText: `${change > 0 ? '+' : ''}${change}`,
    changeClass: change >= 0 ? 'positive' : 'negative',
    createdAtText: formatDateTime(item.createdAt)
  }
}

function decorateMedal(item = {}) {
  const progress = Number(item.progress || 0)
  const target = Math.max(1, Number(item.target || item.conditionValue || 1))
  return {
    ...item,
    iconSrc: absoluteUrl(item.iconUrl),
    coverSrc: absoluteUrl(item.coverUrl),
    progress,
    target,
    progressText: `${progress}/${target}`,
    progressStyle: `width:${Math.max(0, Math.min(100, Math.round(progress * 100 / target)))}%`,
    statusText: item.unlocked ? '已点亮' : '待解锁',
    unlockedAtText: item.unlockedAt ? formatDateTime(item.unlockedAt) : ''
  }
}

function request(path, method = 'GET', data = {}) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: config.BASE_URL + path,
      method,
      data,
      header: { 'X-User-Key': getUserKey() },
      enableHttp2: true,
      success(response) {
        if (response.statusCode >= 200 && response.statusCode < 300) {
          resolve(response.data)
          return
        }
        const message = (response.data && response.data.error) || `请求失败（${response.statusCode}）`
        reject(new Error(message))
      },
      fail(error) {
        reject(new Error(error.errMsg || '网络请求失败'))
      }
    })
  })
}

function uploadFile(filePath, endpoint) {
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: config.BASE_URL + endpoint,
      filePath,
      name: 'file',
      header: { 'X-User-Key': getUserKey() },
      success(response) {
        let body
        try {
          body = JSON.parse(response.data)
        } catch (error) {
          reject(new Error('上传返回数据异常'))
          return
        }
        if (response.statusCode >= 200 && response.statusCode < 300 && body.url) {
          resolve(body.url)
          return
        }
        reject(new Error(body.error || '图片上传失败'))
      },
      fail(error) {
        reject(new Error(error.errMsg || '图片上传失败'))
      }
    })
  })
}

function uploadImage(filePath) {
  return uploadFile(filePath, '/api/upload')
}

function uploadVideo(filePath) {
  return uploadFile(filePath, '/api/upload/video')
}

module.exports = {
  BASE_URL: config.BASE_URL,
  absoluteUrl,
  getUserKey,
  decorateArtwork,
  decorateClip,
  decorateActivity,
  decorateBooking,
  decorateVillage,
  decorateSpot,
  decorateRoute,
  decorateStory,
  decorateProduct,
  decorateFarmer,
  decorateTask,
  decorateOrder,
  decorateProfile,
  decorateLedger,
  decorateMedal,
  formatDateTime,
  request,
  uploadImage,
  uploadVideo
}


