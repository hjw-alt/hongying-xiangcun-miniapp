# -*- coding: utf-8 -*-
"""把原文档封面页与重新生成的目录、正文合并为最终 PDF。"""

import pathlib
import sys

from pypdf import PdfReader, PdfWriter

HERE = pathlib.Path(__file__).resolve().parent
COVER = HERE / "cover-original.pdf"
BODY = HERE / "proposal-body.pdf"
FINAL = HERE.parent / "红映湘村项目企划书-修订版.pdf"
FALLBACK = HERE.parent / "红映湘村项目企划书-修订版2.pdf"


def main():
    if not COVER.exists():
        raise SystemExit("缺少原封面页 cover-original.pdf")
    if not BODY.exists():
        raise SystemExit("缺少正文 proposal-body.pdf，请先运行 print.mjs")

    writer = PdfWriter()
    cover = PdfReader(str(COVER))
    writer.add_page(cover.pages[0])

    body = PdfReader(str(BODY))
    for page in body.pages:
        writer.add_page(page)

    target = FINAL
    try:
        with open(target, "wb") as stream:
            writer.write(stream)
    except PermissionError:
        # 目标文件可能正在被 PDF 阅读器打开，改为另存新版本
        target = FALLBACK
        with open(target, "wb") as stream:
            writer.write(stream)
    print(f"merged: {target} ({len(writer.pages)} pages)")


if __name__ == "__main__":
    sys.exit(main())
