# -*- coding: utf-8 -*-
"""成品校验：页码一致性、文字越界、空白页与图形渲染。"""

import json
import pathlib
import re

import pdfplumber

ROOT = pathlib.Path(__file__).resolve().parents[3]
PDF = ROOT / "output" / "pdf" / "红映湘村项目企划书-修订版.pdf"


def main():
    pdf = pdfplumber.open(PDF)
    pages = [page.extract_text() or "" for page in pdf.pages]
    report = {"pages": len(pdf.pages)}

    heads = {}
    for index, text in enumerate(pages, 1):
        if index == 2:
            continue
        for line in text.splitlines():
            match = re.match(r"^([一二三四五六七八九十]+)、(\S.*)$", line.strip())
            if match and match.group(1) not in heads:
                heads[match.group(1)] = (index, match.group(2))
    report["chapter_pages"] = {key: value[0] for key, value in heads.items()}

    toc = {}
    for line in pages[1].splitlines():
        match = re.match(r"^([一二三四五六七八九十]+)、(\S.*?)\s+(\d+)$", line.strip())
        if match:
            toc[match.group(1)] = int(match.group(3))
    report["toc_pages"] = toc
    report["mismatch"] = {
        key: {"toc": toc.get(key), "actual": heads.get(key, (None,))[0]}
        for key in set(toc) | set(heads)
        if toc.get(key) != heads.get(key, (None,))[0]
    }

    issues = []
    for index, page in enumerate(pdf.pages, 1):
        for char in page.chars:
            if char["x0"] < 38 or char["x1"] > 558 or char["top"] < 30 or char["bottom"] > 812:
                issues.append(
                    {
                        "page": index,
                        "x": round(char["x0"]),
                        "y": round(char["top"]),
                        "text": char["text"],
                    }
                )
    report["out_of_bounds_count"] = len(issues)
    report["out_of_bounds_sample"] = issues[:10]

    report["vector_pages"] = [
        index
        for index, page in enumerate(pdf.pages, 1)
        if len(page.curves) + len(page.rects) > 40
    ]
    report["empty_pages"] = [index for index, text in enumerate(pages, 1) if len(text.strip()) < 10]
    report["new_chapter_present"] = {
        "ch6": "红映湘村小程序系统" in "".join(pages),
        "ch7": "红映湘村文创设计体系" in "".join(pages),
        "architecture": "系统总体架构" in "".join(pages),
        "data_flow": "内容与数据运营闭环" in "".join(pages),
        "product_path": "文创产品开发与转化路径" in "".join(pages),
    }
    report["figure_captions"] = re.findall(r"图 \d+[^\n]{0,30}", "".join(pages))

    out = pathlib.Path(__file__).parent / "verify.json"
    out.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(report, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
