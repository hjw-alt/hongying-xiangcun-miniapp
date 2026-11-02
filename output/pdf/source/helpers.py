# -*- coding: utf-8 -*-
"""HTML 片段构造工具。"""

from html import escape as esc


def chapter(num, title):
    return f'<h2 class="chapter-title">{esc(num)}、{esc(title)}</h2>'


def sub(title):
    return f"<h3>{esc(title)}</h3>"


def subsub(title):
    return f"<h4>{esc(title)}</h4>"


def para(text, cls=None):
    attr = f' class="{cls}"' if cls else ""
    return f"<p{attr}>{esc(text)}</p>"


def box(text):
    return f'<div class="box">{esc(text)}</div>'


def note(text):
    return f'<p class="note">{esc(text)}</p>'


def table(headers, rows, widths=None):
    colgroup = ""
    if widths:
        colgroup = "<colgroup>" + "".join(f'<col style="width:{w}">' for w in widths) + "</colgroup>"
    head = "".join(f"<th>{esc(item)}</th>" for item in headers)
    body = ""
    for row in rows:
        cells = ""
        for index, cell in enumerate(row):
            cls = ' class="key"' if index == 0 else ""
            cells += f"<td{cls}>{esc(cell)}</td>"
        body += f"<tr>{cells}</tr>"
    return f"<table>{colgroup}<thead><tr>{head}</tr></thead><tbody>{body}</tbody></table>"


def figure(svg, caption, number):
    return (
        '<div class="figure">'
        + svg
        + f'<div class="fig-cap">图 {number}　{esc(caption)}</div>'
        + "</div>"
    )


def cards(items):
    """items: list of (title, text)"""
    inner = "".join(
        f'<div class="card"><h5>{esc(title)}</h5><p>{esc(text)}</p></div>' for title, text in items
    )
    return f'<div class="grid2">{inner}</div>'
