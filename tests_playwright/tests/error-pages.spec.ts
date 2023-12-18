import { test, expect } from "@playwright/test";
import { serverUrl } from "./config";

//page, context,
//page.on("console", (msg) => console.log(msg.text()));

test("redirect to 404", async ({ page }) => {
  const response = await page.goto(serverUrl + "/urlthatdoesntexist");
  const previousResponse = await response
    ?.request()
    .redirectedFrom()
    ?.response();
  const previousResponseStatus = previousResponse?.status();
  expect(previousResponseStatus, "status should be 308").toBe(308);
  expect(page.url(), "url should have '404'").toContain("404");
});

test("devmode prevent cache headers", async ({ page }) => {
  const response = await page.goto(serverUrl + "/");
  const headers = await response?.allHeaders();

  expect(headers?.["cache-control"], "check header 'cache-control'").toBe(
    "no-cache, no-store, must-revalidate"
  );
  expect(headers?.["pragma"], "check header 'pragma'").toBe("no-cache");
  expect(headers?.["expires"], "check header 'expires'").toBe("0");
});
