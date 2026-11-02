import fs from 'fs';
import path from 'path';
import { createRequire } from 'module';
import { fileURLToPath } from 'url';

const here = path.dirname(fileURLToPath(import.meta.url));
const nodeRuntime = 'C:/Users/12495/.cache/codex-runtimes/codex-primary-runtime/dependencies/node';
const require = createRequire(path.join(nodeRuntime, 'index.js'));
const sharp = require('sharp');

const outDir = path.join(here, '..', 'frontend', 'assets', 'icons');

// 统一使用线性图标，与小程序界面的描边风格保持一致
const icons = {
  cart: `
    <circle cx="9" cy="20.5" r="1.4"/>
    <circle cx="18" cy="20.5" r="1.4"/>
    <path d="M2.5 3h2.2l2.3 11.2a1.8 1.8 0 0 0 1.8 1.4h8.6a1.8 1.8 0 0 0 1.8-1.4L21 7H6"/>`,
  location: `
    <path d="M20 10.5c0 5.5-8 11.5-8 11.5s-8-6-8-11.5a8 8 0 0 1 16 0z"/>
    <circle cx="12" cy="10.5" r="3"/>`,
  star: `
    <path d="M12 3.4l2.6 5.4 5.9.8-4.3 4.1 1 5.9-5.2-2.8-5.2 2.8 1-5.9L3.5 9.6l5.9-.8z"/>`,
  starFilled: `
    <path d="M12 3.4l2.6 5.4 5.9.8-4.3 4.1 1 5.9-5.2-2.8-5.2 2.8 1-5.9L3.5 9.6l5.9-.8z" fill="#B32A2A"/>`
};

const variants = [
  { suffix: '', color: '#6B5442' },
  { suffix: '-red', color: '#B32A2A' }
];

await fs.promises.mkdir(outDir, { recursive: true });
let count = 0;
for (const [name, body] of Object.entries(icons)) {
  for (const variant of variants) {
    if (name === 'starFilled' && variant.suffix) continue;
    const stroke = name === 'starFilled' ? 'none' : variant.color;
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="96" height="96" fill="none" stroke="${stroke}" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">${body}</svg>`;
    const file = `${name}${variant.suffix}.png`;
    await sharp(Buffer.from(svg)).png({ compressionLevel: 9 }).toFile(path.join(outDir, file));
    count += 1;
  }
}
console.log(`generated ${count} ui icons into ${outDir}`);
