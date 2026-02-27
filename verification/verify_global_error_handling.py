from playwright.sync_api import sync_playwright, Page, expect

def test_home_error_handling(page: Page):
    # Mock projects API failure
    page.route("**/projects", lambda route: route.fulfill(
        status=500,
        content_type="application/json",
        body='{"error": "Failed to fetch projects"}'
    ))

    # Navigate to home
    page.goto("http://localhost:5174/")

    # Check for error container in HomeView
    # Refactored to use ErrorState component with class .error-state
    error_container = page.locator(".error-state")
    expect(error_container).to_be_visible(timeout=5000)
    expect(error_container).to_contain_text("Failed to load projects.")

    page.screenshot(path="verification/home_error.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_home_error_handling(page)
        finally:
            browser.close()
