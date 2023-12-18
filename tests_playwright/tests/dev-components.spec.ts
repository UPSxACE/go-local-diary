import { test, expect } from "@playwright/test";
import { serverUrl } from "./config";

test("redirect from dev to dev/components", async ({ page }) => {
  const response = await page.goto(serverUrl + "/dev");
  const previousResponse = await response
    ?.request()
    .redirectedFrom()
    ?.response();
  const previousResponseStatus = previousResponse?.status();
  expect(previousResponseStatus, "status should be 308").toBe(308);
  expect(page.url(), "url should have dev/components").toContain(
    "dev/components"
  );
});

// This test is more precise if there are more tabs
test("dev/components ux", async ({ page, request }) => {
  await page.goto(serverUrl + "/dev/components");
  const tabs = await page.locator(".tab a").all();
  const sidebarInfoButton = page.locator("#showcase-header-button-info");
  const homeLink = page.locator("#home-link");
  const refreshButton = page.locator("#refresh-data");

  const waitHtmxLoad = async (func: Function) => {
    const promise = page.locator("body").evaluate((element) => {
      return new Promise((resolve) => {
        element.addEventListener("htmx:load", (event) => {
          return resolve(event);
        });
      });
    });
    func();
    await promise;
  };

  const tabsActive = async () => {
    const tabsActive = await page.locator(".tab-active").count();
    return tabsActive;
  };

  const sidebarInfoVisible = async () => {
    const sidebar = page.locator("#sidebar-info-wrapper");

    const width = await sidebar.evaluate((el) => {
      return window.getComputedStyle(el).getPropertyValue("width");
    });

    return width !== "0px" && width !== "0";
  };
  expect(await sidebarInfoVisible(), "info sidebar should not be visible").toBe(
    false
  );
  expect(await tabsActive(), "should have 0 tabs active").toBe(0);
  if (tabs.length > 0) {
    // click first tab
    await waitHtmxLoad(() => tabs[0].click());
    expect(page.url(), "check url query params").toContain("&ct=0&cp=1&e=1");
    expect(await tabsActive(), "only one 1 tab should be active").toBe(1);
    expect(await sidebarInfoVisible(), "sidebar should be visible").toBe(true);
    // click info sidebar button
    await sidebarInfoButton.click();
    await page.waitForTimeout(500); // wait animation
    expect(await sidebarInfoVisible(), "sidebar should not be visible").toBe(
      false
    );
    if (tabs.length > 1) {
      // click second tab
      await waitHtmxLoad(() => tabs[1].click());
      expect(
        await sidebarInfoVisible(),
        "info sidebar should not be visible"
      ).toBe(false);
      // click info sidebar button
      await sidebarInfoButton.click();
      await page.waitForTimeout(500); // wait animation
      expect(await sidebarInfoVisible(), "info sidebar should be visible").toBe(
        true
      );
      // click other tab
      await waitHtmxLoad(() => tabs[0].click());
      expect(await sidebarInfoVisible(), "info sidebar should be visible").toBe(
        true
      );
      expect(await tabsActive(), "only 1 tab should be active").toBe(1);
    } else {
      // click info sidebar button
      await sidebarInfoButton.click();
      await page.waitForTimeout(500); // wait animation
      expect(await sidebarInfoVisible(), "sidebar should be visible").toBe(
        true
      );
    }
    // click home
    await waitHtmxLoad(() => homeLink.click());
    await page.waitForTimeout(500); // wait animation
    expect(await tabsActive(), "no tab should be active").toBe(0);
    expect(
      await sidebarInfoVisible(),
      "info sidebar should not be visible"
    ).toBe(false);
    // click first tab
    await waitHtmxLoad(() => tabs[0].click());
    expect(await tabsActive(), "only 1 tab should be active").toBe(1);
    expect(await sidebarInfoVisible(), "info sidebar should be visible").toBe(
      true
    );
  }
  const responsePromise = page.waitForResponse(
    serverUrl + "/dev/components/refresh"
  );
  const urlBeforeRefresh = page.url();
  await refreshButton.click();
  expect(
    (await responsePromise).status(),
    "response status should be 307"
  ).toBe(307);
  expect(
    page.url(),
    "url should be the same before and after redirection/request"
  ).toBe(urlBeforeRefresh);
});
