import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  // Route any /api calls to fail, so we see the global error toast
  await page.route('**/api/config', route => {
    route.fulfill({
      status: 500,
      contentType: 'text/html',
      body: '<html><body>500 Internal Server Error</body></html>'
    });
  });

  // Navigate to the app running on localhost
  console.log("Navigating to http://localhost:5173/...");
  await page.goto('http://localhost:5173/');

  // Wait for the global toast container and its error message to appear
  console.log("Waiting for toast...");
  await page.waitForSelector('.toast', { state: 'visible', timeout: 5000 });

  // Wait a small moment for animations to complete
  await page.waitForTimeout(500);

  // Take a screenshot of the toast
  console.log("Taking screenshot...");
  await page.screenshot({ path: '/home/jules/verification/toast_error.png' });
  console.log("Screenshot saved to /home/jules/verification/toast_error.png");

  // Let's also check inline errors inside components.
  // The HomeView uses ErrorState when config/auth fails (if we set it up that way, wait actually HomeView loads Config in main.ts, so maybe we see the main ErrorState or Toast).
  // We'll just screenshot the page to see both the Toast and any inline error/loading states.

  await browser.close();
})();
