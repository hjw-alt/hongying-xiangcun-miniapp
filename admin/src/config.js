// 后台可维护内容定义：新增模块时只需在这里补充一项
export const productCategories = ['红色经典IP', '齐白石艺术', '湘村风物助农', '数字文创']
export const taskTypes = ['现场探访', '史料研读', '主题写作', '写生实践', '数据调查', '小组协作']
export const storyCategories = ['人物故事', '革命事件', '求学立志', '艺术人生', '建设故事', '风物故事']

export const resources = {
  banners: {
    label: '首页轮播',
    desc: '小程序首页顶部轮播图与跳转目标',
    listColumns: [
      { key: 'image_url', label: '图片', type: 'image' },
      { key: 'title', label: '标题' },
      { key: 'subtitle', label: '副标题' },
      { key: 'link_type', label: '跳转' },
      { key: 'sort', label: '排序', width: '80px' }
    ],
    fields: [
      { key: 'title', label: '标题', type: 'text' },
      { key: 'subtitle', label: '副标题', type: 'text' },
      { key: 'image_url', label: '图片', type: 'image' },
      { key: 'link_type', label: '跳转类型', type: 'select', options: ['village', 'tab'] },
      { key: 'link_id', label: '跳转目标', type: 'text', hint: '村落填 village-xxx，标签页填 mall / story / guide / study' },
      { key: 'sort', label: '排序', type: 'number' },
      { key: 'active', label: '是否启用', type: 'bool' }
    ]
  },
  villages: {
    label: '红色村落',
    desc: '村落概况、历史与展示信息',
    listColumns: [
      { key: 'cover_url', label: '封面', type: 'image' },
      { key: 'name', label: '村落' },
      { key: 'region', label: '所在地区' },
      { key: 'spot_count', label: '景点', width: '70px' },
      { key: 'story_count', label: '故事', width: '70px' },
      { key: 'product_count', label: '产品', width: '70px' }
    ],
    fields: [
      { key: 'name', label: '村落名称', type: 'text' },
      { key: 'region', label: '所属地区', type: 'text' },
      { key: 'cover_url', label: '封面图', type: 'image' },
      { key: 'summary', label: '村落简介', type: 'textarea' },
      { key: 'history', label: '历史脉络', type: 'textarea', big: true },
      { key: 'honors', label: '荣誉称号', type: 'list' },
      { key: 'themes', label: '主题标签', type: 'list' },
      { key: 'lat', label: '纬度', type: 'number' },
      { key: 'lng', label: '经度', type: 'number' },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  },
  spots: {
    label: '村落景点',
    desc: '纪念馆、故居等红色点位',
    listColumns: [
      { key: 'image_url', label: '图片', type: 'image' },
      { key: 'name', label: '景点' },
      { key: 'village_id', label: '所属村落', type: 'ref', ref: 'villages' },
      { key: 'category', label: '类型' },
      { key: 'sort', label: '排序', width: '70px' }
    ],
    fields: [
      { key: 'name', label: '景点名称', type: 'text' },
      { key: 'village_id', label: '所属村落', type: 'ref', ref: 'villages' },
      { key: 'category', label: '点位类型', type: 'text' },
      { key: 'summary', label: '点位介绍', type: 'textarea' },
      { key: 'image_url', label: '图片', type: 'image' },
      { key: 'audio_url', label: '音频讲解地址', type: 'text' },
      { key: 'lat', label: '纬度', type: 'number' },
      { key: 'lng', label: '经度', type: 'number' },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  },
  routes: {
    label: '红色线路',
    desc: '主题游线与点位组合',
    listColumns: [
      { key: 'cover_url', label: '封面', type: 'image' },
      { key: 'name', label: '线路' },
      { key: 'theme', label: '主题' },
      { key: 'days', label: '天数', width: '70px' },
      { key: 'distance_km', label: '里程', width: '80px' }
    ],
    fields: [
      { key: 'name', label: '线路名称', type: 'text' },
      { key: 'theme', label: '主题', type: 'text' },
      { key: 'cover_url', label: '封面图', type: 'image' },
      { key: 'summary', label: '线路介绍', type: 'textarea' },
      { key: 'days', label: '天数', type: 'number' },
      { key: 'distance_km', label: '里程（公里）', type: 'number' },
      { key: 'spot_ids', label: '点位 ID 列表', type: 'list' },
      { key: 'spot_names', label: '点位名称列表', type: 'list' },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  },
  stories: {
    label: '红迹故事',
    desc: '数字故事内容与 AR 素材',
    listColumns: [
      { key: 'cover_url', label: '封面', type: 'image' },
      { key: 'title', label: '标题' },
      { key: 'category', label: '分类' },
      { key: 'era', label: '年代' },
      { key: 'views', label: '阅读', width: '80px' }
    ],
    fields: [
      { key: 'title', label: '故事标题', type: 'text' },
      { key: 'village_id', label: '所属村落', type: 'ref', ref: 'villages' },
      { key: 'category', label: '分类', type: 'select', options: storyCategories },
      { key: 'era', label: '年代', type: 'text' },
      { key: 'cover_url', label: '封面图', type: 'image' },
      { key: 'summary', label: '摘要', type: 'textarea' },
      { key: 'content', label: '正文', type: 'textarea', big: true },
      { key: 'audio_url', label: '音频地址', type: 'text' },
      { key: 'ar_image', label: 'AR 画面', type: 'image' },
      { key: 'views', label: '阅读量', type: 'number' }
    ]
  },
  products: {
    label: '红文创商品',
    desc: '文创与助农商品、DIY 定制与故事卡',
    listColumns: [
      { key: 'cover_url', label: '封面', type: 'image' },
      { key: 'name', label: '商品' },
      { key: 'category', label: '分类' },
      { key: 'price_cents', label: '价格(分)', width: '90px' },
      { key: 'stock', label: '库存', width: '70px' },
      { key: 'sales', label: '销量', width: '70px' },
      { key: 'status', label: '状态', width: '70px' }
    ],
    fields: [
      { key: 'name', label: '商品名称', type: 'text' },
      { key: 'category', label: '商品分类', type: 'select', options: productCategories },
      { key: 'cover_url', label: '封面图', type: 'image' },
      { key: 'images', label: '图集', type: 'list' },
      { key: 'price_cents', label: '售价（分）', type: 'number' },
      { key: 'points_price', label: '积分兑换价', type: 'number', hint: '填 0 表示该商品暂不支持积分兑换' },
      { key: 'reward_points', label: '现金购买奖励积分', type: 'number' },
      { key: 'reward_exp', label: '购买奖励经验', type: 'number' },
      { key: 'original_cents', label: '原价（分）', type: 'number' },
      { key: 'stock', label: '库存', type: 'number' },
      { key: 'sales', label: '销量', type: 'number' },
      { key: 'summary', label: '一句话卖点', type: 'textarea' },
      { key: 'detail', label: '商品详情', type: 'textarea', big: true },
      { key: 'story_card', label: '文创故事卡', type: 'textarea' },
      { key: 'tags', label: '标签', type: 'list' },
      { key: 'customizable', label: '支持 DIY 定制', type: 'bool' },
      { key: 'farmer_id', label: '承接农户', type: 'ref', ref: 'farmers' },
      { key: 'farmer_name', label: '农户名称', type: 'text' },
      { key: 'village_id', label: '关联村落', type: 'ref', ref: 'villages' },
      { key: 'status', label: '上架状态', type: 'select', options: ['上架', '下架'] },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  },
  farmers: {
    label: '合作农户',
    desc: '农户、合作社与帮扶进度',
    listColumns: [
      { key: 'avatar_url', label: '头像', type: 'image' },
      { key: 'name', label: '农户 / 工坊' },
      { key: 'village_name', label: '所在村' },
      { key: 'project', label: '帮扶项目' },
      { key: 'support_count', label: '订单', width: '80px' }
    ],
    fields: [
      { key: 'name', label: '农户 / 工坊名称', type: 'text' },
      { key: 'village_id', label: '所属村落', type: 'ref', ref: 'villages' },
      { key: 'village_name', label: '村落名称', type: 'text' },
      { key: 'avatar_url', label: '头像 / 图片', type: 'image' },
      { key: 'intro', label: '农户介绍', type: 'textarea' },
      { key: 'project', label: '帮扶项目', type: 'text' },
      { key: 'support_target', label: '目标订单数', type: 'number' },
      { key: 'support_count', label: '当前支持数', type: 'number' }
    ]
  },
  tasks: {
    label: '研学任务',
    desc: '研学打卡任务与积分设置',
    listColumns: [
      { key: 'cover_url', label: '封面', type: 'image' },
      { key: 'title', label: '任务' },
      { key: 'type', label: '类型' },
      { key: 'village_id', label: '村落', type: 'ref', ref: 'villages' },
      { key: 'points', label: '积分', width: '70px' }
    ],
    fields: [
      { key: 'title', label: '任务标题', type: 'text' },
      { key: 'village_id', label: '所属村落', type: 'ref', ref: 'villages' },
      { key: 'spot_id', label: '关联景点 ID', type: 'text' },
      { key: 'type', label: '任务类型', type: 'select', options: taskTypes },
      { key: 'summary', label: '任务说明', type: 'textarea' },
      { key: 'points', label: '积分', type: 'number' },
      { key: 'reward_exp', label: '奖励经验', type: 'number' },
      { key: 'cover_url', label: '封面图', type: 'image' },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  },
  medals: {
    label: '徽章管理',
    desc: '成长徽章、点亮条件与展示素材',
    listColumns: [
      { key: 'icon_url', label: '图标', type: 'image' },
      { key: 'name', label: '徽章名称' },
      { key: 'subtitle', label: '解锁说明' },
      { key: 'condition_type', label: '条件类型' },
      { key: 'condition_value', label: '目标', width: '70px' },
      { key: 'active', label: '启用', width: '70px' }
    ],
    fields: [
      { key: 'id', label: '徽章 ID', type: 'text', hint: '建议使用英文或短横线，创建后不建议修改' },
      { key: 'name', label: '徽章名称', type: 'text' },
      { key: 'subtitle', label: '解锁说明', type: 'text' },
      { key: 'description', label: '徽章详情', type: 'textarea' },
      { key: 'icon_url', label: '徽章图标', type: 'image' },
      { key: 'cover_url', label: '详情封面', type: 'image' },
      { key: 'condition_type', label: '条件类型', type: 'select', options: ['first_checkin', 'checkin_count', 'points_total', 'first_order', 'order_count', 'experience'] },
      { key: 'condition_value', label: '目标值', type: 'number' },
      { key: 'sort', label: '排序', type: 'number' },
      { key: 'active', label: '是否启用', type: 'bool' }
    ]
  },
  point_rules: {
    label: '积分/成长规则',
    desc: '统一维护积分奖励与成长经验规则',
    listColumns: [
      { key: 'name', label: '规则名称' },
      { key: 'event_type', label: '事件类型' },
      { key: 'points', label: '积分', width: '70px' },
      { key: 'experience', label: '经验', width: '70px' },
      { key: 'active', label: '启用', width: '70px' }
    ],
    fields: [
      { key: 'id', label: '规则 ID', type: 'text', hint: '建议使用英文或短横线' },
      { key: 'name', label: '规则名称', type: 'text' },
      { key: 'event_type', label: '事件类型', type: 'select', options: ['task_checkin', 'cash_order', 'points_order'] },
      { key: 'points', label: '积分变化', type: 'number' },
      { key: 'experience', label: '奖励经验', type: 'number' },
      { key: 'description', label: '规则说明', type: 'textarea' },
      { key: 'active', label: '是否启用', type: 'bool' },
      { key: 'sort', label: '排序', type: 'number' }
    ]
  }
}
