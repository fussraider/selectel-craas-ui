from playwright.sync_api import sync_playwright, Page, expect

def test_sidebar_retry(page: Page):
    # 1. Mock projects API to FAIL initially
    page.route("**/projects", lambda route: route.fulfill(
        status=500,
        content_type="application/json",
        body='{"error": "Failed to fetch projects"}'
    ))

    # Navigate to home
    page.goto("http://localhost:5174/")

    # Check for error state in Sidebar
    sidebar = page.locator("aside")
    error_state = sidebar.locator(".error-state")
    expect(error_state).to_be_visible(timeout=5000)
    expect(error_state).to_contain_text("Failed to load registries.") # Or whatever generic message it defaults to if projects fail
    # Note: RepositorySidebar checks store.registries. If projects fail, registries is empty. store.error is set.
    # Title in Sidebar is "Failed to load registries."

    # 2. Update mock to SUCCEED
    page.unroute("**/projects") # Remove failure mock
    page.route("**/projects", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='[{"id": "p1", "name": "project1"}]'
    ))
    page.route("**/projects/p1/registries", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='[{"id": "r1", "name": "registry1"}]'
    ))
    page.route("**/projects/p1/registries/r1/repositories", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='[]'
    ))

    # 3. Click Retry in Sidebar
    retry_btn = error_state.locator("button", has_text="Retry")
    retry_btn.click()

    # 4. Verify Error Gone and Registry Visible
    expect(error_state).not_to_be_visible(timeout=5000)
    expect(sidebar.locator("text=registry1")).to_be_visible()

    page.screenshot(path="verification/sidebar_retry_success.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_sidebar_retry(page)
        finally:
            browser.close()
