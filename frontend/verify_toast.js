import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  // Route any /api calls to immediately fail so we simulate an error
  await page.route('**/api/config', route => {
    route.fulfill({
      status: 500,
      contentType: 'text/html',
      body: '<html><body>500 Internal Server Error</body></html>'
    });
  });

  // Load the page
  console.log("Loading page...");
  await page.goto('http://localhost:5173/');

  try {
    console.log("Waiting for toast container...");
    await page.waitForSelector('.toast-container', { timeout: 2000 });
    console.log("Toast container found. Waiting for a toast...");
    await page.waitForSelector('.toast', { timeout: 5000 });
    console.log("Toast found successfully!");

    // Check toast content
    const toastText = await page.locator('.toast').innerText();
    console.log("Toast content:", toastText);
  } catch (e) {
    console.error("Failed:", e.message);
    const html = await page.locator('.toast-container').innerHTML();
    console.log("Toast container HTML:", html);
  }

  await browser.close();
})();
