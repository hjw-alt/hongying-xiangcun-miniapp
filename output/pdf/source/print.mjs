import { createRequire } from 'module';
import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';

const runtime = 'C:/Users/12495/.cache/codex-runtimes/codex-primary-runtime/dependencies/node';
const require = createRequire(path.join(runtime, 'index.js'));
const { chromium } = require('playwright');

const here = path.dirname(fileURLToPath(import.meta.url));
const htmlPath = path.join(here, 'proposal.html');
const outDir = path.resolve(here, '..');
const outFile = path.join(here, 'proposal-body.pdf');

fs.mkdirSync(outDir, { recursive: true });

const browser = await chromium.launch({ executablePath: 'C:/Program Files/Google/Chrome/Application/chrome.exe' });
const page = await browser.newPage();
await page.emulateMedia({ media: 'print' });
await page.goto('file:///' + htmlPath.replace(/\\/g, '/'), { waitUntil: 'load' });
await page.waitForTimeout(600);

const report = await page.evaluate(() => {
  const pages = [...document.querySelectorAll('.page')];
  const overflow = [];
  pages.forEach((section, index) => {
    const body = section.querySelector('.page-body') || section;
    const over = body.scrollHeight - body.clientHeight;
    if (over > 2) {
      overflow.push({ page: index + 1, over, scrollHeight: body.scrollHeight, clientHeight: body.clientHeight });
    }
  });
  return { total: pages.length, overflow };
});

console.log(JSON.stringify(report, null, 2));

await page.pdf({
  path: outFile,
  format: 'A4',
  printBackground: true,
  preferCSSPageSize: true,
  margin: { top: '15mm', right: '15mm', bottom: '14mm', left: '15mm' }
});

await browser.close();
console.log('pdf: ' + outFile);
