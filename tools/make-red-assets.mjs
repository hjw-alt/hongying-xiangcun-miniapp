import fs from 'fs';
import path from 'path';
import { createRequire } from 'module';
import { fileURLToPath } from 'url';

const here = path.dirname(fileURLToPath(import.meta.url));
const nodeRuntime = 'C:/Users/12495/.cache/codex-runtimes/codex-primary-runtime/dependencies/node';
const require = createRequire(path.join(nodeRuntime, 'index.js'));
const sharp = require('sharp');

const outDir = path.join(here, '..', 'backend', 'static', 'red');

const RED = '#B32A2A';
const DEEP = '#7E1B1B';
const GOLD = '#C79A46';
const INK = '#33221C';
const CREAM = '#FBF3E6';
const JADE = '#4E6B55';
const SKY = '#F6E2C8';

const palettes = {
  red: [RED, DEEP, GOLD, CREAM],
  dawn: ['#C8452F', '#8E2119', '#E0A94F', '#FDF1DE'],
  jade: [JADE, '#2F4636', GOLD, '#F4EEDC'],
  gold: ['#D9A441', '#A9702A', RED, '#FFF6E3'],
  ink: ['#4A3A32', INK, GOLD, '#F7EFE2'],
  sky: ['#C97B4A', '#8E3B2A', '#E8C88A', '#FFF4E2'],
};

function frame(w, h, palette, motif) {
  const [, deep, accent, light] = palette;
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${w} ${h}" width="${w}" height="${h}">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="${light}"/>
      <stop offset="1" stop-color="${accent}" stop-opacity="0.55"/>
    </linearGradient>
    <radialGradient id="glow" cx="0.5" cy="0.28" r="0.6">
      <stop offset="0" stop-color="#FFFFFF" stop-opacity="0.85"/>
      <stop offset="1" stop-color="#FFFFFF" stop-opacity="0"/>
    </radialGradient>
  </defs>
  <rect width="${w}" height="${h}" fill="url(#bg)"/>
  <rect width="${w}" height="${h}" fill="url(#glow)"/>
  ${motif}
  <rect x="${w * 0.045}" y="${h * 0.045}" width="${w * 0.91}" height="${h * 0.91}" fill="none" stroke="${deep}" stroke-opacity="0.22" stroke-width="2"/>
</svg>`;
}

function hills(w, h, palette) {
  const [main, deep, accent] = palette;
  return `
  <path d="M0 ${h * 0.72} C ${w * 0.18} ${h * 0.56}, ${w * 0.3} ${h * 0.66}, ${w * 0.45} ${h * 0.6} C ${w * 0.62} ${h * 0.53}, ${w * 0.74} ${h * 0.66}, ${w} ${h * 0.56} L ${w} ${h} L 0 ${h} Z" fill="${deep}" opacity="0.24"/>
  <path d="M0 ${h * 0.8} C ${w * 0.22} ${h * 0.7}, ${w * 0.42} ${h * 0.82}, ${w * 0.62} ${h * 0.74} C ${w * 0.8} ${h * 0.67}, ${w * 0.9} ${h * 0.76}, ${w} ${h * 0.72} L ${w} ${h} L 0 ${h} Z" fill="${main}" opacity="0.5"/>
  <circle cx="${w * 0.78}" cy="${h * 0.26}" r="${w * 0.09}" fill="${accent}" opacity="0.9"/>`;
}

function house(w, h, palette, x, y, scale) {
  const [, deep, accent] = palette;
  const bw = 150 * scale;
  const bh = 96 * scale;
  return `
  <g transform="translate(${x}, ${y})">
    <path d="M${-bw * 0.62} 0 L0 ${-bh * 0.72} L${bw * 0.62} 0 Z" fill="${deep}" opacity="0.9"/>
    <rect x="${-bw * 0.5}" y="0" width="${bw}" height="${bh * 0.72}" fill="#FFFFFF" opacity="0.85" stroke="${deep}" stroke-width="3"/>
    <rect x="${-bw * 0.1}" y="${bh * 0.22}" width="${bw * 0.2}" height="${bh * 0.5}" fill="${deep}" opacity="0.7"/>
    <circle cx="${-bw * 0.28}" cy="${bh * 0.28}" r="${bw * 0.07}" fill="${accent}" opacity="0.85"/>
    <circle cx="${bw * 0.28}" cy="${bh * 0.28}" r="${bw * 0.07}" fill="${accent}" opacity="0.85"/>
  </g>`;
}

function flag(w, h, palette, x, y) {
  const [main] = palette;
  return `<g transform="translate(${x}, ${y})">
    <rect x="0" y="0" width="6" height="120" fill="${INK}" opacity="0.8"/>
    <path d="M6 8 L86 26 L6 52 Z" fill="${main}"/>
  </g>`;
}

function museumMotif(w, h, palette) {
  const [, deep, accent] = palette;
  const baseY = h * 0.74;
  const cols = [0, 1, 2, 3, 4].map(i => `<rect x="${w * (0.3 + i * 0.1)}" y="${baseY - 120}" width="${w * 0.045}" height="120" fill="#FFFFFF" opacity="0.88" stroke="${deep}" stroke-width="2"/>`).join('');
  return `${hills(w, h, palette)}
  <path d="M${w * 0.26} ${baseY - 128} L${w * 0.5} ${baseY - 186} L${w * 0.74} ${baseY - 128} Z" fill="${deep}" opacity="0.9"/>
  ${cols}
  <rect x="${w * 0.24}" y="${baseY}" width="${w * 0.52}" height="${h * 0.06}" fill="${deep}" opacity="0.85"/>
  <rect x="${w * 0.24}" y="${baseY + h * 0.06}" width="${w * 0.52}" height="${h * 0.04}" fill="${accent}" opacity="0.7"/>`;
}

function statueMotif(w, h, palette) {
  const [, deep, accent] = palette;
  return `${hills(w, h, palette)}
  <rect x="${w * 0.38}" y="${h * 0.66}" width="${w * 0.24}" height="${h * 0.14}" fill="${deep}" opacity="0.85"/>
  <path d="M${w * 0.44} ${h * 0.66} L${w * 0.47} ${h * 0.32} L${w * 0.53} ${h * 0.32} L${w * 0.56} ${h * 0.66} Z" fill="${deep}" opacity="0.92"/>
  <circle cx="${w * 0.5}" cy="${h * 0.27}" r="${w * 0.055}" fill="${deep}" opacity="0.92"/>
  <path d="M${w * 0.53} ${h * 0.4} L${w * 0.66} ${h * 0.34}" stroke="${deep}" stroke-width="10" stroke-linecap="round" opacity="0.92"/>
  <circle cx="${w * 0.78}" cy="${h * 0.22}" r="${w * 0.07}" fill="${accent}" opacity="0.85"/>`;
}

function lotusMotif(w, h, palette) {
  const [, deep, accent] = palette;
  const petals = [0, 1, 2, 3, 4].map(i => {
    const angle = -50 + i * 25;
    return `<ellipse cx="0" cy="-42" rx="15" ry="46" fill="#FFFFFF" opacity="0.92" stroke="${deep}" stroke-width="2" transform="rotate(${angle})"/>`;
  }).join('');
  return `${hills(w, h, palette)}
  <g transform="translate(${w * 0.5}, ${h * 0.55})">${petals}</g>
  <ellipse cx="${w * 0.28}" cy="${h * 0.72}" rx="${w * 0.14}" ry="${h * 0.06}" fill="${deep}" opacity="0.65"/>
  <ellipse cx="${w * 0.72}" cy="${h * 0.76}" rx="${w * 0.16}" ry="${h * 0.06}" fill="${deep}" opacity="0.55"/>
  <circle cx="${w * 0.5}" cy="${h * 0.55}" r="${w * 0.035}" fill="${accent}"/>`;
}

function bambooMotif(w, h, palette) {
  const [, deep, accent] = palette;
  const stalks = [0.3, 0.42, 0.6, 0.72].map((ratio, index) => {
    const x = w * ratio;
    const top = h * (0.16 + (index % 2) * 0.08);
    return `<rect x="${x}" y="${top}" width="${w * 0.028}" height="${h * 0.66}" fill="${deep}" opacity="0.78"/>
    ${[0, 1, 2, 3].map(i => `<ellipse cx="${x + (index % 2 ? -w * 0.06 : w * 0.06)}" cy="${top + 40 + i * 52}" rx="${w * 0.06}" ry="${h * 0.014}" fill="${accent}" opacity="0.7" transform="rotate(${index % 2 ? 18 : -18} ${x} ${top + 40 + i * 52})"/>`).join('')}`;
  }).join('');
  return `${hills(w, h, palette)}${stalks}`;
}

function shrimpMotif(w, h, palette) {
  const [, deep, accent] = palette;
  return `<rect x="${w * 0.12}" y="${h * 0.12}" width="${w * 0.76}" height="${h * 0.76}" fill="#FFFDF7" stroke="${deep}" stroke-opacity="0.3" stroke-width="2"/>
  <path d="M${w * 0.3} ${h * 0.58} C ${w * 0.38} ${h * 0.36}, ${w * 0.58} ${h * 0.32}, ${w * 0.62} ${h * 0.5} C ${w * 0.66} ${h * 0.66}, ${w * 0.44} ${h * 0.72}, ${w * 0.36} ${h * 0.6}" fill="none" stroke="${INK}" stroke-width="6" stroke-linecap="round"/>
  <path d="M${w * 0.62} ${h * 0.48} L${w * 0.78} ${h * 0.4} M${w * 0.62} ${h * 0.52} L${w * 0.8} ${h * 0.52} M${w * 0.6} ${h * 0.56} L${w * 0.76} ${h * 0.64}" stroke="${INK}" stroke-width="4" stroke-linecap="round"/>
  <circle cx="${w * 0.34}" cy="${h * 0.47}" r="4" fill="${INK}"/>
  <path d="M${w * 0.3} ${h * 0.58} L${w * 0.2} ${h * 0.68}" stroke="${INK}" stroke-width="3" stroke-linecap="round"/>
  <circle cx="${w * 0.82}" cy="${h * 0.24}" r="${w * 0.04}" fill="${accent}" opacity="0.8"/>`;
}

function fieldMotif(w, h, palette) {
  const [, deep, accent] = palette;
  const rows = [0, 1, 2, 3, 4].map(i => `<path d="M0 ${h * (0.58 + i * 0.08)} Q ${w * 0.5} ${h * (0.5 + i * 0.08)} ${w} ${h * (0.6 + i * 0.08)}" fill="none" stroke="${deep}" stroke-width="4" opacity="${0.5 - i * 0.07}"/>`).join('');
  return `${hills(w, h, palette)}${rows}
  <circle cx="${w * 0.2}" cy="${h * 0.2}" r="${w * 0.06}" fill="${accent}" opacity="0.9"/>`;
}

function mapMotif(w, h, palette) {
  const [, deep, accent] = palette;
  const pins = [[0.28, 0.34], [0.52, 0.24], [0.72, 0.44], [0.44, 0.62], [0.66, 0.72]].map(([px, py], index) =>
    `<g transform="translate(${w * px}, ${h * py})">
      <circle r="${index === 0 ? 14 : 11}" fill="${index === 0 ? accent : deep}" opacity="0.9"/>
      <circle r="5" fill="#FFFFFF"/>
    </g>`).join('');
  return `<rect x="${w * 0.15}" y="${h * 0.14}" width="${w * 0.7}" height="${h * 0.72}" rx="14" fill="#FFFDF7" stroke="${deep}" stroke-opacity="0.3" stroke-width="3"/>
  <path d="M${w * 0.2} ${h * 0.4} Q ${w * 0.42} ${h * 0.2} ${w * 0.56} ${h * 0.4} T ${w * 0.8} ${h * 0.62}" fill="none" stroke="${RED}" stroke-width="5" stroke-dasharray="12 10"/>
  ${pins}`;
}

function productMotif(w, h, palette, shape) {
  const [main, deep, accent] = palette;
  const cx = w * 0.5;
  const cy = h * 0.54;
  switch (shape) {
    case 'bag':
      return `<path d="M${cx - w * 0.2} ${cy - h * 0.16} H${cx + w * 0.2} L${cx + w * 0.24} ${cy + h * 0.26} H${cx - w * 0.24} Z" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <path d="M${cx - w * 0.1} ${cy - h * 0.16} C${cx - w * 0.1} ${cy - h * 0.3}, ${cx + w * 0.1} ${cy - h * 0.3}, ${cx + w * 0.1} ${cy - h * 0.16}" fill="none" stroke="${deep}" stroke-width="5"/>
      <rect x="${cx - w * 0.12}" y="${cy}" width="${w * 0.24}" height="${h * 0.14}" fill="${main}" opacity="0.85"/>
      <path d="M${cx - w * 0.1} ${cy + h * 0.07} H${cx + w * 0.1}" stroke="#FFFFFF" stroke-width="3"/>`;
    case 'badge':
      return `<circle cx="${cx}" cy="${cy}" r="${w * 0.22}" fill="${deep}"/>
      <circle cx="${cx}" cy="${cy}" r="${w * 0.17}" fill="${accent}" opacity="0.9"/>
      <path d="M${cx - w * 0.08} ${cy + h * 0.02} L${cx + w * 0.08} ${cy - h * 0.06}" stroke="#FFFDF7" stroke-width="6" stroke-linecap="round"/>
      <circle cx="${cx - w * 0.05}" cy="${cy - h * 0.05}" r="5" fill="#FFFDF7"/>`;
    case 'sticker':
      return [0, 1, 2, 3].map(i => {
        const x = cx - w * 0.2 + (i % 2) * w * 0.24;
        const y = cy - h * 0.16 + Math.floor(i / 2) * h * 0.2;
        return `<rect x="${x}" y="${y}" width="${w * 0.17}" height="${h * 0.16}" rx="8" fill="#FFFDF7" stroke="${deep}" stroke-width="3" transform="rotate(${i % 2 ? 6 : -6} ${x + w * 0.08} ${y + h * 0.08})"/>
        <path d="M${x + w * 0.03} ${y + h * 0.11} L${x + w * 0.08} ${y + h * 0.04} L${x + w * 0.14} ${y + h * 0.11}" fill="none" stroke="${main}" stroke-width="3"/>`;
      }).join('');
    case 'box':
      return `<rect x="${cx - w * 0.19}" y="${cy - h * 0.1}" width="${w * 0.38}" height="${h * 0.3}" rx="10" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <path d="M${cx - w * 0.19} ${cy - h * 0.1} L${cx} ${cy - h * 0.24} L${cx + w * 0.19} ${cy - h * 0.1} L${cx} ${cy + h * 0.04} Z" fill="${accent}" opacity="0.85"/>
      <circle cx="${cx + w * 0.06}" cy="${cy + h * 0.08}" r="${w * 0.05}" fill="${main}" opacity="0.85"/>`;
    case 'scarf':
      return `<rect x="${cx - w * 0.24}" y="${cy - h * 0.18}" width="${w * 0.48}" height="${h * 0.36}" rx="10" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <path d="M${cx - w * 0.16} ${cy + h * 0.06} C ${cx - w * 0.04} ${cy - h * 0.12}, ${cx + w * 0.06} ${cy + h * 0.1}, ${cx + w * 0.18} ${cy - h * 0.08}" fill="none" stroke="${main}" stroke-width="5" stroke-linecap="round"/>
      <path d="M${cx - w * 0.16} ${cy + h * 0.12} C ${cx - w * 0.02} ${cy - h * 0.02}, ${cx + w * 0.08} ${cy + h * 0.16}, ${cx + w * 0.18} ${cy - h * 0.02}" fill="none" stroke="${accent}" stroke-width="4" stroke-linecap="round"/>`;
    case 'seal':
      return `<rect x="${cx - w * 0.16}" y="${cy - h * 0.16}" width="${w * 0.32}" height="${h * 0.32}" rx="14" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <rect x="${cx - w * 0.1}" y="${cy - h * 0.1}" width="${w * 0.2}" height="${h * 0.2}" rx="8" fill="${main}"/>
      <path d="M${cx - w * 0.05} ${cy} H${cx + w * 0.05} M${cx} ${cy - h * 0.05} V${cy + h * 0.05}" stroke="#FFFDF7" stroke-width="4" stroke-linecap="round"/>`;
    case 'card':
      return `<rect x="${cx - w * 0.22}" y="${cy - h * 0.15}" width="${w * 0.44}" height="${h * 0.3}" rx="10" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <rect x="${cx - w * 0.17}" y="${cy - h * 0.09}" width="${w * 0.16}" height="${h * 0.18}" fill="${main}" opacity="0.85"/>
      <path d="M${cx - w * 0.14} ${cy + h * 0.06} H${cx + w * 0.16} M${cx - w * 0.14} ${cy + h * 0.11} H${cx + w * 0.08}" stroke="${deep}" stroke-width="3" stroke-linecap="round"/>`;
    case 'book':
      return `<rect x="${cx - w * 0.18}" y="${cy - h * 0.18}" width="${w * 0.36}" height="${h * 0.36}" rx="6" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <rect x="${cx - w * 0.18}" y="${cy - h * 0.18}" width="${w * 0.07}" height="${h * 0.36}" fill="${main}"/>
      <path d="M${cx - w * 0.06} ${cy - h * 0.06} H${cx + w * 0.12} M${cx - w * 0.06} ${cy + h * 0.02} H${cx + w * 0.12}" stroke="${accent}" stroke-width="4" stroke-linecap="round"/>`;
    case 'bottle':
      return `<rect x="${cx - w * 0.1}" y="${cy - h * 0.22}" width="${w * 0.2}" height="${h * 0.42}" rx="12" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <rect x="${cx - w * 0.05}" y="${cy - h * 0.28}" width="${w * 0.1}" height="${h * 0.08}" fill="${deep}"/>
      <rect x="${cx - w * 0.08}" y="${cy - h * 0.04}" width="${w * 0.16}" height="${h * 0.14}" fill="${main}" opacity="0.85"/>`;
    case 'jar':
      return `<path d="M${cx - w * 0.14} ${cy - h * 0.16} H${cx + w * 0.14} L${cx + w * 0.11} ${cy + h * 0.18} H${cx - w * 0.11} Z" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
      <rect x="${cx - w * 0.16}" y="${cy - h * 0.22}" width="${w * 0.32}" height="${h * 0.07}" rx="6" fill="${main}"/>
      <circle cx="${cx}" cy="${cy + h * 0.04}" r="${w * 0.05}" fill="${accent}" opacity="0.85"/>`;
    case 'digital':
      return `<rect x="${cx - w * 0.2}" y="${cy - h * 0.14}" width="${w * 0.4}" height="${h * 0.28}" rx="12" fill="${deep}"/>
      <rect x="${cx - w * 0.14}" y="${cy - h * 0.07}" width="${w * 0.12}" height="${h * 0.14}" fill="#FFFDF7" opacity="0.92"/>
      <path d="M${cx + w * 0.02} ${cy - h * 0.05} H${cx + w * 0.14} M${cx + w * 0.02} ${cy} H${cx + w * 0.14} M${cx + w * 0.02} ${cy + h * 0.05} H${cx + w * 0.1}" stroke="${accent}" stroke-width="4" stroke-linecap="round"/>`;
    default:
      return `<circle cx="${cx}" cy="${cy}" r="${w * 0.2}" fill="${main}" opacity="0.85"/>`;
  }
}

function taskMotif(w, h, palette) {
  const [, deep, accent] = palette;
  return `${hills(w, h, palette)}
  <rect x="${w * 0.26}" y="${h * 0.3}" width="${w * 0.48}" height="${h * 0.44}" rx="12" fill="#FFFDF7" stroke="${deep}" stroke-width="4"/>
  <rect x="${w * 0.3}" y="${h * 0.24}" width="${w * 0.16}" height="${h * 0.09}" rx="6" fill="${deep}"/>
  <path d="M${w * 0.36} ${h * 0.5} L${w * 0.44} ${h * 0.6} L${w * 0.62} ${h * 0.42}" fill="none" stroke="${accent}" stroke-width="8" stroke-linecap="round" stroke-linejoin="round"/>`;
}

function arMotif(w, h, palette) {
  const [main, deep, accent] = palette;
  return `<rect x="${w * 0.14}" y="${h * 0.16}" width="${w * 0.72}" height="${h * 0.68}" fill="#F2E3CB" stroke="${deep}" stroke-width="4"/>
  ${house(w, h, palette, w * 0.5, h * 0.68, 1.05).replace('opacity="0.85"', 'opacity="0.55"')}
  <rect x="${w * 0.14}" y="${h * 0.16}" width="${w * 0.72}" height="${h * 0.68}" fill="${main}" opacity="0.12"/>
  ${[[0.18, 0.2], [0.82, 0.2], [0.18, 0.8], [0.82, 0.8]].map(([px, py]) =>
    `<path d="M${w * px - 18} ${h * py} h18 v18 M${w * px} ${h * py - 18} v18 h-18" stroke="${accent}" stroke-width="4" fill="none"/>`).join('')}`;
}

const assets = [
  ['banner-village.png', 1200, 720, 'dawn', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.36, h * 0.66, 1.5)}${flag(w, h, p, w * 0.66, h * 0.24)}`],
  ['banner-market.png', 1200, 720, 'gold', (w, h, p) => `${productMotif(w, h, p, 'bag')}${productMotif(w * 0.62, h, p, 'jar').replace(/translate/g, 'translate')}${flag(w, h, p, w * 0.82, h * 0.2)}`],
  ['banner-story.png', 1200, 720, 'red', (w, h, p) => `${mapMotif(w, h, p)}`],
  ['village-shaoshan.png', 900, 640, 'dawn', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.42, h * 0.68, 1.4)}${flag(w, h, p, w * 0.7, h * 0.26)}`],
  ['village-wushi.png', 900, 640, 'ink', (w, h, p) => `${hills(w, h, p)}${statueMotif(w, h, p)}`],
  ['village-quanhu.png', 900, 640, 'jade', (w, h, p) => `${fieldMotif(w, h, p)}${house(w, h, p, w * 0.72, h * 0.6, 0.9)}`],
  ['village-xingtang.png', 900, 640, 'sky', (w, h, p) => `${hills(w, h, p)}${lotusMotif(w, h, p)}`],
  ['village-yinshan.png', 900, 640, 'gold', (w, h, p) => `${fieldMotif(w, h, p)}`],
  ['village-dongshan.png', 900, 640, 'red', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['spot-maoguju.png', 800, 600, 'dawn', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.5, h * 0.7, 1.5)}`],
  ['spot-tongxiang.png', 800, 600, 'gold', (w, h, p) => `${statueMotif(w, h, p)}`],
  ['spot-museum.png', 800, 600, 'red', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['spot-pengjinianguan.png', 800, 600, 'red', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['spot-pengguju.png', 800, 600, 'ink', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.5, h * 0.7, 1.4)}`],
  ['spot-defengmu.png', 800, 600, 'jade', (w, h, p) => `${hills(w, h, p)}${statueMotif(w, h, p)}`],
  ['spot-chengengguju.png', 800, 600, 'ink', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.44, h * 0.7, 1.3)}`],
  ['spot-quanhu-tian.png', 800, 600, 'jade', (w, h, p) => `${fieldMotif(w, h, p)}`],
  ['spot-xingdoutang.png', 800, 600, 'sky', (w, h, p) => `${hills(w, h, p)}${lotusMotif(w, h, p)}`],
  ['spot-jieshan.png', 800, 600, 'sky', (w, h, p) => `${shrimpMotif(w, h, p)}`],
  ['spot-yintianlaojie.png', 800, 600, 'gold', (w, h, p) => `${hills(w, h, p)}${house(w, h, p, w * 0.34, h * 0.68, 1.2)}${house(w, h, p, w * 0.68, h * 0.7, 1.0)}`],
  ['spot-guanqu.png', 800, 600, 'jade', (w, h, p) => `${fieldMotif(w, h, p)}`],
  ['spot-dongshantang.png', 800, 600, 'red', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['spot-dongshan-museum.png', 800, 600, 'gold', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['route-shaoshan.png', 900, 640, 'dawn', (w, h, p) => `${mapMotif(w, h, p)}`],
  ['route-wushi.png', 900, 640, 'ink', (w, h, p) => `${mapMotif(w, h, p)}`],
  ['route-art.png', 900, 640, 'sky', (w, h, p) => `${mapMotif(w, h, p)}`],
  ['story-qiushi.png', 900, 640, 'red', (w, h, p) => `${museumMotif(w, h, p)}`],
  ['story-nongyun.png', 900, 640, 'dawn', (w, h, p) => `${house(w, h, p, w * 0.5, h * 0.68, 1.4)}`],
  ['story-hengdao.png', 900, 640, 'ink', (w, h, p) => `${statueMotif(w, h, p)}`],
  ['story-baishi.png', 900, 640, 'sky', (w, h, p) => `${shrimpMotif(w, h, p)}`],
  ['story-guanqu.png', 900, 640, 'jade', (w, h, p) => `${fieldMotif(w, h, p)}`],
  ['story-xianglian.png', 900, 640, 'gold', (w, h, p) => `${lotusMotif(w, h, p)}`],
  ['ar-dongshan.png', 900, 640, 'red', (w, h, p) => `${arMotif(w, h, p)}`],
  ['ar-shaoshan.png', 900, 640, 'dawn', (w, h, p) => `${arMotif(w, h, p)}`],
  ['ar-wushi.png', 900, 640, 'ink', (w, h, p) => `${arMotif(w, h, p)}`],
  ['ar-xingtang.png', 900, 640, 'sky', (w, h, p) => `${arMotif(w, h, p)}`],
  ['ar-yinshan.png', 900, 640, 'jade', (w, h, p) => `${arMotif(w, h, p)}`],
  ['ar-xianglian.png', 900, 640, 'gold', (w, h, p) => `${arMotif(w, h, p)}`],
  ['farmer-lotus.png', 800, 800, 'gold', (w, h, p) => `${lotusMotif(w, h, p)}`],
  ['farmer-bamboo.png', 800, 800, 'jade', (w, h, p) => `${bambooMotif(w, h, p)}`],
  ['farmer-tea.png', 800, 800, 'jade', (w, h, p) => `${fieldMotif(w, h, p)}`],
  ['farmer-pepper.png', 800, 800, 'red', (w, h, p) => `${productMotif(w, h, p, 'jar')}`],
  ['task-maoguju.png', 800, 600, 'dawn', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['task-nongyun.png', 800, 600, 'red', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['task-hengdao.png', 800, 600, 'ink', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['task-baishi.png', 800, 600, 'sky', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['task-guanqu.png', 800, 600, 'jade', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['task-chengeng.png', 800, 600, 'gold', (w, h, p) => `${taskMotif(w, h, p)}`],
  ['product-bag.png', 800, 800, 'red', (w, h, p) => `${productMotif(w, h, p, 'bag')}`],
  ['product-badge.png', 800, 800, 'gold', (w, h, p) => `${productMotif(w, h, p, 'badge')}`],
  ['product-sticker.png', 800, 800, 'dawn', (w, h, p) => `${productMotif(w, h, p, 'sticker')}`],
  ['product-blindbox.png', 800, 800, 'red', (w, h, p) => `${productMotif(w, h, p, 'box')}`],
  ['product-scarf.png', 800, 800, 'sky', (w, h, p) => `${productMotif(w, h, p, 'scarf')}`],
  ['product-seal.png', 800, 800, 'ink', (w, h, p) => `${productMotif(w, h, p, 'seal')}`],
  ['product-postcard.png', 800, 800, 'gold', (w, h, p) => `${productMotif(w, h, p, 'card')}`],
  ['product-book.png', 800, 800, 'dawn', (w, h, p) => `${productMotif(w, h, p, 'book')}`],
  ['product-lotus.png', 800, 800, 'gold', (w, h, p) => `${productMotif(w, h, p, 'jar')}`],
  ['product-camellia.png', 800, 800, 'jade', (w, h, p) => `${productMotif(w, h, p, 'bottle')}`],
  ['product-bamboo.png', 800, 800, 'jade', (w, h, p) => `${productMotif(w, h, p, 'box')}`],
  ['product-chili.png', 800, 800, 'red', (w, h, p) => `${productMotif(w, h, p, 'jar')}`],
  ['product-digital.png', 800, 800, 'ink', (w, h, p) => `${productMotif(w, h, p, 'digital')}`],
];

await fs.promises.mkdir(outDir, { recursive: true });
let count = 0;
for (const [name, w, h, paletteName, build] of assets) {
  const palette = palettes[paletteName];
  const svg = frame(w, h, palette, build(w, h, palette).replace(/translate\(0, 0\)/g, ''));
  await sharp(Buffer.from(svg)).png({ compressionLevel: 9 }).toFile(path.join(outDir, name));
  count += 1;
}
console.log(`generated ${count} assets into ${outDir}`);
