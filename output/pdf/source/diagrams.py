# -*- coding: utf-8 -*-
"""企划书中的技术流程图与结构图，使用内联 SVG 以保证印刷清晰度。"""

from html import escape as esc

RED = "#8E1D25"
RED2 = "#A03132"
GOLD = "#C39A4B"
TINT_RED = "#FBEEEE"
TINT_GOLD = "#FBF3E4"
TINT_GREY = "#F3F1EC"
LINE = "#E1D6CD"


def _open(width, height):
    return (
        f'<svg viewBox="0 0 {width} {height}" xmlns="http://www.w3.org/2000/svg" '
        'font-family="Microsoft YaHei, Noto Sans SC, sans-serif">'
        "<defs>"
        '<marker id="arrow" markerWidth="9" markerHeight="9" refX="7.5" refY="3" orient="auto">'
        f'<path d="M0,0 L7.5,3 L0,6 z" fill="{RED}"/>'
        "</marker>"
        '<linearGradient id="band" x1="0" y1="0" x2="1" y2="0">'
        f'<stop offset="0" stop-color="{GOLD}"/><stop offset="1" stop-color="{RED}"/>'
        "</linearGradient>"
        "</defs>"
    )


def _arrow(x1, y1, x2, y2, color=RED):
    return (
        f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" '
        'stroke-width="1.2" marker-end="url(#arrow)"/>'
    )


def _text(x, y, content, size=10, fill="#4A4038", weight="400", anchor="middle", spacing="0"):
    return (
        f'<text x="{x}" y="{y}" text-anchor="{anchor}" font-size="{size}" fill="{fill}" '
        f'font-weight="{weight}" letter-spacing="{spacing}">{esc(content)}</text>'
    )


def architecture():
    w, h = 660, 380
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(_text(w / 2, 22, "红映湘村小程序系统总体架构", 12.5, RED, "700", spacing="0.6"))
    layers = [
        (34, "用户层", ["游客与青年群体", "研学学生与教师", "村集体与经营者"], TINT_RED),
        (110, "应用层 · 微信小程序", ["红色地图导览", "数字展陈讲解", "研学任务打卡", "文创商城", "个人中心"], TINT_GOLD),
        (204, "服务层 · Go 后端服务", ["内容与路线", "用户与收藏", "订单与预约", "文件上传", "管理接口鉴权"], TINT_GREY),
        (298, "数据与运营层", ["MySQL 业务数据库", "图片与音频资源库", "Vue 3 内容管理后台"], TINT_RED),
    ]
    for y, name, items, fill in layers:
        out.append(f'<rect x="14" y="{y}" width="{w - 28}" height="62" rx="6" fill="{fill}" stroke="{LINE}"/>')
        out.append(_text(26, y + 26, name, 11.5, RED, "700", "start"))
        count = len(items)
        available = w - 200
        gap = 8
        box_w = (available - gap * (count - 1)) / count
        for index, item in enumerate(items):
            bx = 162 + index * (box_w + gap)
            out.append(f'<rect x="{bx:.1f}" y="{y + 13}" width="{box_w:.1f}" height="36" rx="4" fill="#FFFFFF" stroke="#DDD3CB"/>')
            out.append(_text(bx + box_w / 2, y + 36, item, 10, "#4A4038"))
        if y < 290:
            out.append(_arrow(w / 2, y + 62, w / 2, y + 74))
    out.append("</svg>")
    return "".join(out)


def modules():
    w, h = 660, 336
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(f'<rect x="238" y="10" width="184" height="36" rx="18" fill="{RED}"/>')
    out.append(_text(330, 33, "红映湘村小程序", 12.5, "#FFFFFF", "700"))
    items = [
        ("红色地图导览", ["主题路线推荐", "景点地图标记", "步行时长提示", "离线文字导览"]),
        ("数字展陈讲解", ["语音讲解点播", "图文故事展陈", "老照片对比", "扫码听讲解"]),
        ("研学任务打卡", ["任务卡与答题", "定位打卡集章", "成果海报生成", "小组排行"]),
        ("文创与农产品", ["文创商品详情", "农产品专区", "联名款预售", "订单与物流"]),
        ("个人中心", ["我的收藏", "我的打卡", "我的订单", "我的作品"]),
    ]
    box_w, gap, start_x = 122, 10.5, 12
    for index, (title, children) in enumerate(items):
        x = start_x + index * (box_w + gap)
        out.append(f'<rect x="{x}" y="76" width="{box_w}" height="238" rx="6" fill="{TINT_GOLD}" stroke="{LINE}"/>')
        out.append(f'<rect x="{x}" y="76" width="{box_w}" height="34" rx="6" fill="{RED2}"/>')
        out.append(_text(x + box_w / 2, 99, title, 10.6, "#FFFFFF", "700"))
        for child_index, child in enumerate(children):
            cy = 124 + child_index * 44
            out.append(f'<rect x="{x + 9}" y="{cy}" width="{box_w - 18}" height="32" rx="4" fill="#FFFFFF" stroke="#DDD3CB"/>')
            out.append(_text(x + box_w / 2, cy + 21, child, 10, "#4A4038"))
        out.append(f'<line x1="{x + box_w / 2}" y1="46" x2="{x + box_w / 2}" y2="76" stroke="{GOLD}" stroke-width="1.2"/>')
    out.append(f'<line x1="{start_x + box_w / 2}" y1="60" x2="{start_x + 4 * (box_w + gap) + box_w / 2}" y2="60" stroke="{GOLD}" stroke-width="1.2"/>')
    out.append("</svg>")
    return "".join(out)


def journey():
    w, h = 660, 258
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(_text(w / 2, 24, "游客使用流程：从“看到”到“看懂”再到“带走”", 12, RED, "700", spacing="0.4"))
    steps = [
        ("① 扫码进入", ["村口与展板二维码", "微信搜索进入"], TINT_RED),
        ("② 定位选线", ["获取当前位置", "推荐主题路线"], TINT_GOLD),
        ("③ 边走边听", ["语音讲解", "图文与老照片"], TINT_GREY),
        ("④ 任务打卡", ["完成答题任务", "集章生成海报"], TINT_RED),
        ("⑤ 转化下单", ["文创与农产品", "研学活动报名"], TINT_GOLD),
        ("⑥ 分享复访", ["分享成果海报", "收藏与复购"], TINT_GREY),
    ]
    box_w, gap, start_x = 94, 9.6, 12
    for index, (title, lines, fill) in enumerate(steps):
        x = start_x + index * (box_w + gap)
        out.append(f'<rect x="{x}" y="54" width="{box_w}" height="122" rx="6" fill="{fill}" stroke="{LINE}"/>')
        out.append(_text(x + box_w / 2, 80, title, 10.8, RED, "700"))
        for line_index, line in enumerate(lines):
            out.append(_text(x + box_w / 2, 104 + line_index * 17, line, 9.4))
        if index < len(steps) - 1:
            out.append(_arrow(x + box_w, 115, x + box_w + gap - 1, 115))
    out.append(f'<rect x="12" y="196" width="636" height="40" rx="4" fill="{TINT_GOLD}" stroke="{LINE}"/>')
    out.append(_text(330, 213, "线下触点：二维码展板 · 语音讲解牌 · 打卡印章点 · 文创展销台 · 研学集合点", 9.6, "#6E5447"))
    out.append(_text(330, 229, "线上触点：小程序首页 · 路线详情页 · 任务页 · 商品详情页 · 个人中心", 9.6, "#6E5447"))
    out.append("</svg>")
    return "".join(out)


def dataflow():
    w, h = 660, 272
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(_text(w / 2, 26, "内容与数据运营闭环", 12, RED, "700", spacing="0.4"))
    nodes = [
        (18, "游客行为", ["浏览与打卡", "收藏与下单"]),
        (186, "Go 后端服务", ["业务接口", "数据校验"]),
        (354, "MySQL 数据库", ["内容与订单", "用户与行为"]),
        (522, "管理后台", ["内容维护", "数据看板"]),
    ]
    for x, title, lines in nodes:
        out.append(f'<rect x="{x}" y="52" width="120" height="80" rx="6" fill="#FCFAF7" stroke="{LINE}"/>')
        out.append(_text(x + 60, 76, title, 10.8, RED, "700"))
        for index, line in enumerate(lines):
            out.append(_text(x + 60, 98 + index * 17, line, 9.6))
    out.append(_arrow(138, 92, 184, 92))
    out.append(_arrow(306, 92, 352, 92))
    out.append(_arrow(474, 92, 520, 92))
    out.append(f'<path d="M582 132 L582 176 L78 176 L78 134" fill="none" stroke="{RED}" stroke-width="1.2" marker-end="url(#arrow)"/>')
    out.append(_text(330, 170, "内容更新 · 路线优化 · 商品调整 · 新活动上线", 9.8, "#6E5447"))
    out.append(f'<rect x="150" y="200" width="360" height="52" rx="6" fill="{TINT_RED}" stroke="{LINE}"/>')
    out.append(_text(330, 221, "运营决策依据", 10.5, RED, "700"))
    out.append(_text(330, 239, "到访人数 · 打卡完成率 · 停留时长 · 转化率 · 复购与满意度", 9.4))
    out.append("</svg>")
    return "".join(out)


def cultural_products():
    w, h = 660, 312
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(_text(w / 2, 26, "文创产品开发与转化路径", 12, RED, "700", spacing="0.4"))
    stages = [
        ("文化资源采集", ["村史与人物", "建筑与器物", "口述与影像"]),
        ("设计元素提取", ["纹样与图形", "色彩与材质", "字体与版式"]),
        ("IP 与视觉系统", ["IP 形象设计", "标志与辅助形", "视觉规范手册"]),
        ("产品线设计", ["纪念文创", "生活文创", "农产品包装"]),
        ("打样与实测", ["小批量打样", "用户试用", "成本核算"]),
        ("量产与渠道", ["本地工坊生产", "小程序与村内", "研学与展会展销"]),
    ]
    box_w, gap, start_x = 94, 9.6, 12
    for index, (title, lines) in enumerate(stages):
        x = start_x + index * (box_w + gap)
        out.append(f'<rect x="{x}" y="52" width="{box_w}" height="132" rx="6" fill="{TINT_GOLD}" stroke="{LINE}"/>')
        out.append(f'<rect x="{x}" y="52" width="{box_w}" height="34" rx="6" fill="{GOLD}"/>')
        out.append(_text(x + box_w / 2, 74, title, 10.2, "#FFFFFF", "700"))
        for line_index, line in enumerate(lines):
            out.append(_text(x + box_w / 2, 104 + line_index * 18, line, 9.4))
        if index < len(stages) - 1:
            out.append(_arrow(x + box_w, 118, x + box_w + gap - 1, 118))
    out.append(f'<rect x="12" y="204" width="636" height="88" rx="6" fill="{TINT_RED}" stroke="{LINE}"/>')
    out.append(_text(28, 226, "收益共享机制", 11, RED, "700", "start"))
    out.append(_text(28, 246, "设计费、销售分成与研学服务收入按约定比例回流村集体和农户，用于红色资源维护与乡村美育。", 9.6, "#4A4038", "400", "start"))
    out.append(_text(28, 266, "所有产品先小批量测试，再依据真实销量决定量产规模，避免库存积压与一次性创作。", 9.6, "#4A4038", "400", "start"))
    out.append(_text(28, 284, "包装与视觉统一使用同一套设计规范，保证线上线下形象一致、便于复制到其他村落。", 9.6, "#4A4038", "400", "start"))
    out.append("</svg>")
    return "".join(out)


def roadmap():
    w, h = 660, 250
    out = [_open(w, h), f'<rect width="{w}" height="{h}" fill="#FFFFFF"/>']
    out.append(_text(w / 2, 26, "小程序与内容迭代路线", 12, RED, "700", spacing="0.4"))
    out.append(f'<line x1="70" y1="118" x2="590" y2="118" stroke="#DDD3CB" stroke-width="3"/>')
    stages = [
        (70, "第 1-3 月", "原型期", ["需求调研与内容梳理", "数据表与接口设计", "小程序框架搭建"]),
        (330, "第 4-9 月", "试点期", ["示范村内容上线", "任务与打卡测试", "文创商品试销"]),
        (590, "第 10-12 月", "复制期", ["标准化内容模板", "多村数据接入", "运营数据复盘"]),
    ]
    for x, period, name, lines in stages:
        anchor = "middle"
        out.append(f'<circle cx="{x}" cy="118" r="9" fill="{RED}"/>')
        out.append(_text(x, 96, period, 11, RED, "700", anchor))
        out.append(_text(x, 148, name, 11, RED2, "700", anchor))
        for index, line in enumerate(lines):
            out.append(_text(x, 170 + index * 16, line, 9.4, "#4A4038", "400", anchor))
    out.append(f'<rect x="12" y="212" width="636" height="28" rx="4" fill="{TINT_GOLD}" stroke="{LINE}"/>')
    out.append(_text(330, 231, "每阶段设置可验证成果：内容上线记录 · 用户测试反馈 · 成本收入台账", 9.6, "#6E5447"))
    out.append("</svg>")
    return "".join(out)
