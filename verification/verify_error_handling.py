from playwright.sync_api import sync_playwright, Page, expect

def test_error_handling(page: Page):
    # Mock the images API to return a 500 error
    page.route("**/images?repository=repo1", lambda route: route.fulfill(
        status=500,
        content_type="application/json",
        body='{"error": "Internal Server Error"}'
    ))

    # Navigate to the repository view
    # We use dummy IDs since we are mocking the response
    # Port changed to 5174 based on logs
    page.goto("http://localhost:5174/projects/p1/registries/r1/repositories/repo1")

    # Expect the error state to be visible
    # The error state div has class 'error-state' and contains "Failed to load images."
    error_state = page.locator(".error-state")
    expect(error_state).to_be_visible(timeout=10000)
    expect(error_state).to_contain_text("Failed to load images.")

    # Also verify the toast notification if possible, though it might auto-dismiss
    # The toast component usually renders with a specific class or role
    # Based on previous file reads, it's ToastNotification

    # Take screenshot
    page.screenshot(path="verification/error_state.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_error_handling(page)
        finally:
            browser.close()
