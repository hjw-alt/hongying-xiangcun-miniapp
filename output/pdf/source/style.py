# -*- coding: utf-8 -*-
"""企划书排版的样式系统，沿用原文档的红色标题与 A4 版心。"""

CSS = """
@page { size: A4; margin: 15mm 15mm 14mm 15mm; }
* { box-sizing: border-box; }
html, body { margin: 0; padding: 0; }
body {
  /* 使用微软雅黑：可变字体 Noto Sans SC 会让 PDF 文本层写入部首码位，影响复制与检索 */
  font-family: "Microsoft YaHei", "Noto Sans SC", sans-serif;
  color: #2F2A27;
  font-size: 10.4pt;
  line-height: 1.76;
  -webkit-print-color-adjust: exact;
  print-color-adjust: exact;
}
.page {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 180mm;
  height: 265mm;
  overflow: hidden;
  break-after: page;
  page-break-after: always;
}
.page:last-child { break-after: auto; page-break-after: auto; }
.page-body { flex: 1 1 auto; min-height: 0; }
.page-foot {
  flex: 0 0 auto;
  display: flex;
  justify-content: space-between;
  padding-top: 2.6mm;
  border-top: 0.6pt solid #E4DAD2;
  color: #9A8C82;
  font-size: 8.2pt;
  letter-spacing: 0.4pt;
}
.cover { justify-content: center; align-items: center; text-align: center; }
.cover-mark { width: 22mm; height: 1.6mm; margin-bottom: 11mm; background: linear-gradient(90deg, #C39A4B 0%, #8E1D25 100%); }
.cover-title { margin: 0; color: #9D2028; font-size: 40pt; font-weight: 700; letter-spacing: 13pt; line-height: 1.32; }
.cover-sub { margin-top: 7mm; color: #6E5447; font-size: 15pt; letter-spacing: 1.4pt; }
.cover-line { width: 96mm; height: 0.8pt; margin: 8mm 0; background: #DDD3CB; }
.cover-note { color: #8E1D25; font-size: 13pt; letter-spacing: 6pt; }
.cover-meta { margin-top: 15mm; color: #4A4038; font-size: 10.5pt; line-height: 2.05; }
.cover-meta b { color: #8E1D25; font-weight: 600; }
.cover-disclaimer { margin-top: 12mm; padding-top: 4mm; border-top: 0.6pt solid #E4DAD2; color: #9A8C82; font-size: 8.5pt; }
.chapter-title {
  margin: 0 0 3mm;
  padding-bottom: 2.2mm;
  border-bottom: 1.4pt solid #8E1D25;
  color: #8E1D25;
  font-size: 19pt;
  font-weight: 700;
  letter-spacing: 1pt;
}
h3 { margin: 4.2mm 0 1.7mm; color: #A03132; font-size: 11.5pt; font-weight: 700; letter-spacing: 0.5pt; }
h4 { margin: 3.4mm 0 1.2mm; color: #8E1D25; font-size: 10.4pt; font-weight: 700; }
p { margin: 0 0 2.2mm; text-align: justify; }
.note { color: #8A7C72; font-size: 9.2pt; line-height: 1.68; }
.box { margin: 3mm 0; padding: 3.2mm 4mm; border-left: 2.4pt solid #C39A4B; background: #FBF3E4; }
.box strong { color: #8E1D25; }
.grid2 { display: grid; grid-template-columns: 1fr 1fr; gap: 3.2mm; }
.grid3 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 3mm; }
.card { padding: 3mm 3.4mm; border: 0.6pt solid #E1D6CD; border-radius: 1.4mm; background: #FCFAF7; }
.card h5 { margin: 0 0 1.2mm; color: #8E1D25; font-size: 10pt; font-weight: 700; }
.card p { margin: 0; font-size: 9.3pt; line-height: 1.64; color: #4A4038; }
.card ul { margin: 0; padding-left: 4.4mm; font-size: 9.3pt; line-height: 1.64; color: #4A4038; }
table { width: 100%; border-collapse: collapse; margin: 2.4mm 0 2.8mm; font-size: 9.2pt; }
th, td { border: 0.6pt solid #E1D6CD; padding: 1.6mm 2.1mm; vertical-align: top; line-height: 1.58; }
th { background: #FBEEEE; color: #8E1D25; font-weight: 700; text-align: left; }
tbody tr:nth-child(even) td { background: #FCFAF7; }
td.key { color: #8E1D25; font-weight: 600; }
.figure { margin: 3.2mm 0 3.4mm; text-align: center; }
.figure svg { width: 100%; height: auto; }
.fig-cap { margin-top: 1.4mm; color: #8A7C72; font-size: 8.6pt; letter-spacing: 0.2pt; }
.toc-item { display: flex; align-items: baseline; gap: 2mm; padding: 1.4mm 0; border-bottom: 0.5pt dotted #E1D6CD; font-size: 10pt; color: #8E1D25; }
.toc-item .toc-name { flex: 0 0 auto; }
.toc-dots { flex: 1 1 auto; height: 0.5pt; background: repeating-linear-gradient(90deg,#CFC4BB 0 1pt,transparent 1pt 3pt); }
.toc-page { flex: 0 0 auto; color: #8E1D25; font-weight: 600; }
.toc-group { margin: 2.4mm 0 0; color: #8E1D25; font-size: 9pt; letter-spacing: 1pt; }
"""
