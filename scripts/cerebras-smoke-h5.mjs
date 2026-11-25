import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const { origin, login, password, screenshot } = JSON.parse(readFileSync(0, 'utf8'));
assert.match(origin, /^http:\/\/127\.0\.0\.1:\d+$/);
const { default: puppeteer } = await import(process.env.CAMS_PUPPETEER_MODULE || 'puppeteer-core');
const prompt = 'Give three short numbered checklist items for a fictional campus workshop. Start each with its number.';
const browser = await puppeteer.launch({
    executablePath: process.env.CAMS_CHROME || '/usr/bin/google-chrome',
    headless: true,
    args: ['--no-sandbox', '--disable-background-networking'],
});

try {
    const page = await browser.newPage();
    page.setDefaultTimeout(90000);
    await page.setViewport({ width: 390, height: 844, isMobile: true, hasTouch: true });
    await page.setRequestInterception(true);
    page.on('request', request => {
        if (request.url().startsWith(origin + '/') || request.url().startsWith('data:')) void request.continue();
        else void request.abort();
    });

    const dashboard = encodeURIComponent('/modules/president/pages/president-activity-list');
    await page.goto(origin + '/#/modules/auth/pages/login?returnUrl=' + dashboard, { waitUntil: 'domcontentloaded' });
    const loginResponse = page.waitForResponse(response => response.url() === origin + '/api/auth/login');
    await page.locator('input.uni-input-input:not([type=password])').fill(login, { typingThreshold: 0 });
    await page.locator('input.uni-input-input[type=password]').fill(password, { typingThreshold: 0 });
    await page.locator('uni-button').filter(button => button.textContent.trim() === 'Sign in').click();
    assert.equal((await loginResponse).status(), 200, 'synthetic president signed in');
    await page.locator('.assistant-entry').click();
    await page.waitForFunction(() => document.querySelector('.new-chat') &&
        !document.body.innerText.includes('Updating conversations…'));
    await page.locator('.new-chat').click();
    await page.locator('.assistant-composer textarea').fill(prompt);

    const replyResponse = page.waitForResponse(response => response.request().method() === 'POST' &&
        response.url().startsWith(origin + '/api/') && response.url().endsWith('/reply'));
    await page.locator('.assistant-composer .send').click();
    const response = await replyResponse;
    assert.equal(response.status(), 200, 'CAMS returned a completed hosted reply');
    const result = await response.json();
    assert.equal(result.prompt, prompt);
    const lines = result.reply.trim().split(/\r?\n/).filter(line => line.trim());
    assert.equal(lines.length, 3, 'hosted model returned three nonempty lines');
    lines.forEach((line, index) => assert.match(line, new RegExp('^\\s*' + (index + 1) + '[.)]\\s+\\S')));

    await page.waitForFunction(expected => document.querySelector('.transcript .reply')?.textContent.trim() === expected,
        {}, result.reply.trim());
    await page.reload({ waitUntil: 'domcontentloaded' });
    await page.waitForFunction(expected => document.querySelector('.transcript .reply')?.textContent.trim() === expected,
        {}, result.reply.trim());
    assert.equal(await page.$$eval('.transcript .reply', replies => replies.length), 1);
    if (screenshot) await page.screenshot({ path: screenshot, fullPage: true });
    console.log('PASS: hosted reply reached H5 and reopened as one complete saved turn');
    console.log(result.reply);
} finally {
    await browser.close();
}
