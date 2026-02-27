from playwright.sync_api import sync_playwright, Page, expect

def test_sidebar_error_handling(page: Page):
    # Mock projects API to succeed (so sidebar loads)
    page.route("**/projects", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='[{"id": "p1", "name": "project1"}]'
    ))

    # Mock registries API failure (this is what the sidebar fetches for the project)
    page.route("**/projects/p1/registries", lambda route: route.fulfill(
        status=500,
        content_type="application/json",
        body='{"error": "Failed to fetch registries"}'
    ))

    # Navigate to home
    page.goto("http://localhost:5174/")

    # Check for error state in Sidebar
    # The sidebar is an aside element
    sidebar = page.locator("aside")
    error_state = sidebar.locator(".error-state")
    expect(error_state).to_be_visible(timeout=5000)
    expect(error_state).to_contain_text("Failed to load registries.")

    page.screenshot(path="verification/sidebar_error.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_sidebar_error_handling(page)
        finally:
            browser.close()
