# -*- coding: utf-8 -*-
"""装配《红映湘村项目企划书》HTML。"""

import pathlib

import content_a
import content_b
import content_c
import content_d
from style import CSS
from helpers import esc

OUT_HTML = pathlib.Path(__file__).resolve().parent / "proposal.html"

FOOT_LEFT = "红映湘村项目企划书"


def chapters():
    return [
        ("一", "项目摘要", ["1.1 项目基本信息", "1.2 项目概述", "1.3 项目愿景与使命", "1.4 红旅特色"], content_a.ch1()),
        ("二", "项目背景与政策依据", ["2.1 国家层面的发展方向", "2.2 湖南省与湘潭市的现实基础", "2.3 学校与专业基础"], content_a.ch2()),
        ("三", "现实问题与市场机会", ["3.1 红色村落面临的主要问题", "3.2 项目机会", "3.3 竞争格局"], content_a.ch3()),
        ("四", "项目定位与目标", ["4.1 项目定位", "4.2 核心价值主张", "4.3 阶段性目标", "4.4 首年建议目标"], content_a.ch4()),
        ("五", "产品与服务系统", ["5.1 红色文化品牌服务", "5.2 数字导览与互动内容", "5.3 文创产品与农产品包装", "5.4 研学课程与活动运营", "5.5 数字资产与年度维护", "5.6 产品服务矩阵"], content_a.ch5()),
        ("六", "红映湘村小程序系统", ["6.1 建设定位与用户价值", "6.2 系统总体架构", "6.3 功能模块设计", "6.4 核心功能与展示内容", "6.5 关键技术实现", "6.6 游客使用流程", "6.7 内容与数据运营闭环", "6.8 技术可行性与已有基础", "6.9 迭代路线"], content_b.ch6_pages()),
        ("七", "红映湘村文创设计体系", ["7.1 设计定位与原则", "7.2 IP 形象与视觉识别系统", "7.3 纹样与视觉元素提取", "7.4 产品线与包装设计", "7.5 文创与农产品联动", "7.6 线上线下体验闭环", "7.7 打样与量产路径"], content_c.ch7_pages()),
        ("八", "设计创新与文化真实性保障", ["8.1 设计创新", "8.2 文化真实性原则", "8.3 知识产权与素材合规", "8.4 可复制方法"], content_d.ch8()),
        ("九", "目标用户与市场需求", ["9.1 用户画像", "9.2 市场进入策略", "9.3 需求验证方法"], content_d.ch9()),
        ("十", "商业模式与盈利路径", ["10.1 收入来源", "10.2 成本结构", "10.3 定价策略", "10.4 收益分配与公益反哺"], content_d.ch10()),
        ("十一", "市场推广与运营方案", ["11.1 品牌传播", "11.2 渠道建设", "11.3 用户运营", "11.4 服务质量控制"], content_d.ch11()),
        ("十二", "实施计划与阶段目标", ["12.1 第一年实施安排", "12.2 三年发展路径", "12.3 关键里程碑"], content_d.ch12()),
        ("十三", "团队构成与组织管理", ["13.1 团队角色", "13.2 决策机制", "13.3 能力建设"], content_d.ch13()),
        ("十四", "财务规划与资金用途", ["14.1 测算说明", "14.2 启动资金建议", "14.3 三年收入规划"], content_d.ch14() + content_d.ch14b()),
        ("十五", "社会价值与育人价值", ["15.1 文化价值", "15.2 经济价值", "15.3 社会与育人价值", "15.4 规划社会效益指标"], content_d.ch15()),
        ("十六", "风险分析与应对措施", ["风险类别与应对措施"], content_d.ch16()),
        ("十七", "评价指标与可持续发展", ["17.1 四维评价体系", "17.2 数据台账", "17.3 可持续机制"], content_d.ch17()),
    ]


def cover_page():
    body = (
        '<div class="cover-mark"></div>'
        '<h1 class="cover-title">红映湘村</h1>'
        '<div class="cover-sub">湘潭红色村落数字文旅与文创助农计划</div>'
        '<div class="cover-line"></div>'
        '<div class="cover-note">项目企划书</div>'
        '<div class="cover-meta">'
        "<div>参赛赛道：青年红色筑梦之旅赛道（红旅赛道）</div>"
        "<div>项目所在地：湖南省湘潭市</div>"
        "<div>申报单位：湘潭大学艺术学院学生团队</div>"
        "<div>项目负责人：袁纳一　　指导教师：陈娟</div>"
        "<div>日期：2026 年 9 月</div>"
        "</div>"
        '<div class="cover-disclaimer">本企划书中的经营数据均为规划测算，不构成真实经营业绩。</div>'
    )
    return f'<section class="page cover">{body}</section>'


def toc_page(entries, first_content_page):
    rows = []
    page = first_content_page
    for number, title, subsections, pages in entries:
        highlight = " toc-new" if number in ("六", "七") else ""
        rows.append(
            f'<div class="toc-item{highlight}">'
            f'<span class="toc-name">{esc(number)}、{esc(title)}</span>'
            '<span class="toc-dots"></span>'
            f'<span class="toc-page">{page}</span>'
            "</div>"
        )
        page += len(pages)
    body = (
        '<h2 class="chapter-title">目录</h2>'
        + "".join(rows)
        + '<div class="toc-group">本版修订重点：新增第六章“红映湘村小程序系统”，包含系统总体架构、功能模块设计、核心技术实现、游客使用流程与数据运营闭环；新增第七章“红映湘村文创设计体系”，包含 IP 与视觉识别、纹样提取、产品线规划与打样量产路径。</div>'
        + '<div class="page-foot"><span>红映湘村项目企划书</span><span>第 2 页</span></div>'
    )
    return f'<section class="page">{body}</section>'


def build():
    entries = chapters()
    first_content_page = 3
    # 封面沿用用户原文档第 1 页（含 logo），由 merge.py 在合并阶段插回，这里只生成目录与正文。
    pages_html = [toc_page(entries, first_content_page)]
    page_number = first_content_page
    for number, title, subsections, chapter_pages in entries:
        for index, body in enumerate(chapter_pages):
            foot = (
                '<div class="page-foot">'
                f"<span>{esc(number)}、{esc(title)}</span>"
                f"<span>第 {page_number} 页</span>"
                "</div>"
            )
            pages_html.append(f'<section class="page"><div class="page-body">{body}</div>{foot}</section>')
            page_number += 1
    document = (
        "<!DOCTYPE html><html lang=\"zh-CN\"><head><meta charset=\"utf-8\">"
        "<title>红映湘村项目企划书</title>"
        f"<style>{CSS}</style></head><body>{''.join(pages_html)}</body></html>"
    )
    OUT_HTML.write_text(document, encoding="utf-8")
    return page_number - 1


if __name__ == "__main__":
    total = build()
    print(f"pages: {total}")
