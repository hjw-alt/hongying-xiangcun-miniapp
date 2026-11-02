import fs from 'fs';
import path from 'path';
import { createRequire } from 'module';
import { fileURLToPath } from 'url';

const here = path.dirname(fileURLToPath(import.meta.url));
const nodeRuntime = 'C:/Users/12495/.cache/codex-runtimes/codex-primary-runtime/dependencies/node';
const require = createRequire(path.join(nodeRuntime, 'index.js'));
const sharp = require('sharp');

const outDir = path.join(here, '..', 'frontend', 'assets', 'tabbar-red');

const icons = {
  home: '<path d="M4 10.5 12 4l8 6.5V20a1 1 0 0 1-1 1h-4.5v-6h-5v6H5a1 1 0 0 1-1-1z"/>',
  mall: '<path d="M5 8h14l-1.2 12a1 1 0 0 1-1 .9H7.2a1 1 0 0 1-1-.9z"/><path d="M9 8V6a3 3 0 0 1 6 0v2"/>',
  guide: '<path d="M12 21s7-6.2 7-11a7 7 0 1 0-14 0c0 4.8 7 11 7 11z"/><circle cx="12" cy="10" r="2.6"/>',
  story: '<path d="M5 5.5A2.5 2.5 0 0 1 7.5 3H19v15.5H7.5A2.5 2.5 0 0 0 5 21z"/><path d="M9 8h6M9 12h4"/>',
  study: '<path d="M12 4 3 8l9 4 9-4z"/><path d="M7 10.5V16c0 1.4 2.2 2.6 5 2.6s5-1.2 5-2.6v-5.5"/>',
  farm: '<path d="M12 21c4.4 0 8-3.6 8-8 0-5.2-8-10-8-10S4 7.8 4 13c0 4.4 3.6 8 8 8z"/><path d="M12 17v-5M9.5 14h5"/>',
};

const colors = { normal: '#9A8578', active: '#B32A2A' };

function svg(name, color) {
  const body = icons[name];
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="81" height="81" fill="none" stroke="${color}" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">${body}</svg>`;
}

await fs.promises.mkdir(outDir, { recursive: true });
let count = 0;
for (const name of Object.keys(icons)) {
  await sharp(Buffer.from(svg(name, colors.normal))).png({ compressionLevel: 9 }).toFile(path.join(outDir, `${name}.png`));
  await sharp(Buffer.from(svg(name, colors.active))).png({ compressionLevel: 9 }).toFile(path.join(outDir, `${name}-active.png`));
  count += 2;
}
console.log(`generated ${count} tabbar icons`);
