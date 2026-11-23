// Zero-dependency Chrome DevTools smoke test for M1 responsive foundation UI.
// Intended for CI on Node 22 with preinstalled Google Chrome/Chromium.
import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import { mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const site = process.env.UI_URL || "http://127.0.0.1:4173/";
const chrome = process.env.CHROME_BIN || "google-chrome";
const work = mkdtempSync(join(tmpdir(), "delta-m1-browser-"));
const screenshots = "browser-screenshots";
mkdirSync(screenshots, { recursive: true });
const browser = spawn(chrome, [
  "--headless=new", "--no-sandbox", "--disable-dev-shm-usage",
  "--disable-gpu", "--no-first-run", "--remote-allow-origins=*",
  "--remote-debugging-port=0", "--user-data-dir=" + work, "about:blank"
], { stdio: ["ignore", "ignore", "pipe"] });
let chromeDiagnostics = "";
browser.stderr?.on("data", (data) => {
  chromeDiagnostics = (chromeDiagnostics + data.toString()).slice(-4000);
});

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
let ws;
try {
  let target;
  for (let i = 0; i < 120; i++) {
    if (browser.exitCode !== null) throw new Error("Chrome exited during startup: " + chromeDiagnostics);
    try {
      const portText = readFileSync(join(work, "DevToolsActivePort"), "utf8").split("\n")[0];
      const port = Number(portText);
      if (!Number.isInteger(port) || port < 1) throw new Error("No debugger port yet");
      const response = await fetch("http://127.0.0.1:" + port + "/json", {
        signal: AbortSignal.timeout(1500)
      });
      const pages = await response.json();
      target = pages.find((page) => page.type === "page");
      if (target?.webSocketDebuggerUrl) break;
    } catch { /* Chrome is still starting */ }
    await sleep(200);
  }
  assert.ok(target?.webSocketDebuggerUrl, "Chrome debugging target unavailable: " + chromeDiagnostics);
  ws = new WebSocket(target.webSocketDebuggerUrl);
  await new Promise((resolve, reject) => {
    ws.addEventListener("open", resolve, { once: true });
    ws.addEventListener("error", reject, { once: true });
  });
  let nextID = 0;
  const pending = new Map();
  ws.addEventListener("message", ({ data }) => {
    const message = JSON.parse(data);
    const operation = pending.get(message.id);
    if (operation) {
      pending.delete(message.id);
      if (message.error) operation.reject(new Error(JSON.stringify(message.error)));
      else operation.resolve(message.result);
    }
  });
  function cdp(method, params = {}) {
    return new Promise((resolve, reject) => {
      const id = ++nextID;
      pending.set(id, { resolve, reject });
      ws.send(JSON.stringify({ id, method, params }));
    });
  }
  async function evaluate(expression) {
    const result = await cdp("Runtime.evaluate", {
      expression, returnByValue: true, awaitPromise: true
    });
    if (result.exceptionDetails) throw new Error(JSON.stringify(result.exceptionDetails));
    return result.result.value;
  }
  async function until(expression, expected, label) {
    for (let i = 0; i < 50; i++) {
      const value = await evaluate(expression);
      if (value === expected) return;
      await sleep(100);
    }
    throw new Error("Timed out waiting for " + label);
  }
  await cdp("Page.enable");
  await cdp("Runtime.enable");
  for (const width of [360, 768, 1280, 1920]) {
    await cdp("Emulation.setDeviceMetricsOverride", {
      width, height: 900, deviceScaleFactor: 1, mobile: width <= 768
    });
    await cdp("Page.navigate", { url: site });
    await until('document.querySelector("main h1")?.textContent', "Trading overview", "overview " + width);
    const state = await evaluate(String.raw`(() => ({
      warning: document.body.innerText.includes("Live trading is disabled."),
      disconnected: document.body.innerText.includes("No account"),
      overflow: document.documentElement.scrollWidth > document.documentElement.clientWidth + 2,
      skipLink: !!document.querySelector('a.skip-link[href="#main"]'),
      main: !!document.querySelector('main#main'),
      unlabeled: Array.from(document.querySelectorAll("button")).filter((button) => {
        const style = getComputedStyle(button);
        if (style.display === "none" || style.visibility === "hidden") return false;
        return !button.getAttribute("aria-label") && !button.textContent.trim() &&
          !button.getAttribute("aria-labelledby");
      }).length
    }))()`);
    assert.ok(state.warning && state.disconnected && state.skipLink && state.main,
      "Foundation safety/accessibility elements absent at " + width + ": " + JSON.stringify(state));
    assert.equal(state.unlabeled, 0, "Unlabeled visible buttons at " + width);
    assert.equal(state.overflow, false, "Horizontal document overflow at " + width);
    const shot = await cdp("Page.captureScreenshot", { format: "png", captureBeyondViewport: false });
    writeFileSync(join(screenshots, "overview-" + width + ".png"), Buffer.from(shot.data, "base64"));
    console.log("PASS overview at " + width + "px");
    if (width === 1280) {
      await evaluate(String.raw`Array.from(document.querySelectorAll("nav button")).find(
        (button) => button.textContent.trim() === "Strategies"
      )?.click()`);
      await until('document.querySelector("main h1")?.textContent', "Strategies", "strategy navigation");
      await until('document.body.innerText.includes("Configure Strategy")', true, "Configure Strategy action");
      const shot2 = await cdp("Page.captureScreenshot", { format: "png" });
      writeFileSync(join(screenshots, "strategies-1280.png"), Buffer.from(shot2.data, "base64"));
      console.log("PASS strategy workspace navigation");
    }
  }

  // A 320 CSS-pixel viewport approximates 200% zoom in a 640px-wide
  // physical viewport. This is a reflow check, not a screen-reader audit.
  await cdp("Emulation.setDeviceMetricsOverride", {
    width: 320, height: 900, deviceScaleFactor: 1, mobile: true
  });
  await cdp("Emulation.setEmulatedMedia", {
    features: [{ name: "prefers-reduced-motion", value: "reduce" }]
  });
  await cdp("Page.navigate", { url: site });
  await until('document.querySelector("main h1")?.textContent', "Trading overview", "narrow reflow");
  assert.equal(
    await evaluate('document.documentElement.scrollWidth > document.documentElement.clientWidth + 2'),
    false, "Horizontal overflow in 320 CSS-pixel narrow-reflow mode"
  );
  assert.equal(await evaluate('matchMedia("(prefers-reduced-motion: reduce)").matches'), true);
  const narrowShot = await cdp("Page.captureScreenshot", { format: "png" });
  writeFileSync(join(screenshots, "overview-320-reflow.png"), Buffer.from(narrowShot.data, "base64"));
  console.log("PASS 320px narrow-reflow and reduced-motion preferences");

  // CSS zoom is an automated 200% layout stress test, not a substitute for
  // browser-native zoom and assistive-technology acceptance on real devices.
  await cdp("Emulation.setDeviceMetricsOverride", {
    width: 640, height: 900, deviceScaleFactor: 1, mobile: false
  });
  await cdp("Page.navigate", { url: site });
  await until('document.querySelector("main h1")?.textContent', "Trading overview", "zoom stress");
  const zoomState = await evaluate(String.raw`(() => {
    // Chrome does not apply root-element zoom consistently in headless mode;
    // zoom the body and confirm the computed value instead.
    document.body.style.zoom = "200%";
    return {
      zoom: parseFloat(getComputedStyle(document.body).zoom),
      inlineZoom: document.body.style.zoom,
      overflow: document.documentElement.scrollWidth > document.documentElement.clientWidth + 2,
      safety: document.body.innerText.includes("Live trading is disabled."),
      skipLink: !!document.querySelector('a.skip-link[href="#main"]')
    };
  })()`);
  assert.ok(zoomState.zoom >= 2 && zoomState.safety && zoomState.skipLink,
    "200% CSS zoom safety check failed: " + JSON.stringify(zoomState));
  assert.equal(zoomState.overflow, false,
    "Horizontal document overflow during 200% CSS zoom stress");
  const zoomShot = await cdp("Page.captureScreenshot", { format: "png" });
  writeFileSync(join(screenshots, "overview-css-zoom-200.png"), Buffer.from(zoomShot.data, "base64"));
  console.log("PASS 200% CSS zoom layout stress (manual browser zoom still required)");

  // Validate the actual accessibility tree names, not merely visible text.
  await cdp("Accessibility.enable");
  const ax = await cdp("Accessibility.getFullAXTree");
  const unnamed = ax.nodes.filter(node =>
    node.role?.value === "button" && !node.ignored && !node.name?.value?.trim()
  );
  assert.equal(unnamed.length, 0, "Unnamed accessible buttons: " + unnamed.map(n => n.nodeId).join(", "));
  console.log("PASS accessible button names");

  // Exercise native keyboard events at a mobile viewport.
  await cdp("Emulation.setDeviceMetricsOverride", {
    width: 360, height: 900, deviceScaleFactor: 1, mobile: true
  });
  await cdp("Page.navigate", { url: site });
  await until('document.querySelector("main h1")?.textContent', "Trading overview", "mobile keyboard");
  await cdp("Input.dispatchKeyEvent", {
    type: "keyDown", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "keyUp", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9
  });
  await until('document.activeElement?.classList.contains("skip-link")', true, "keyboard skip link");
  await evaluate('document.querySelector(".hamburger")?.focus()');
  await cdp("Input.dispatchKeyEvent", {
    type: "rawKeyDown", key: "Enter", code: "Enter", windowsVirtualKeyCode: 13, nativeVirtualKeyCode: 13
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "char", key: "Enter", code: "Enter", text: "\r", unmodifiedText: "\r",
    windowsVirtualKeyCode: 13, nativeVirtualKeyCode: 13
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "keyUp", key: "Enter", code: "Enter", windowsVirtualKeyCode: 13, nativeVirtualKeyCode: 13
  });
  console.log("Keyboard activation state:", await evaluate('JSON.stringify({focused:document.activeElement?.outerHTML?.slice(0,180), expanded:document.querySelector(".hamburger")?.getAttribute("aria-expanded"), display:getComputedStyle(document.querySelector(".hamburger")).display, width:innerWidth})'));
  await until('document.querySelector(".hamburger")?.getAttribute("aria-expanded")', "true", "mobile navigation open");
  console.log("Focus placement state:", await evaluate('JSON.stringify({focused: document.activeElement?.outerHTML?.slice(0,180), close: document.querySelector(".close-nav")?.outerHTML?.slice(0,180), navVisible: getComputedStyle(document.querySelector(".sidebar")).visibility})'));
  await until('document.activeElement?.getAttribute("aria-label")', "Close menu", "menu focus placement");
  await cdp("Input.dispatchKeyEvent", {
    type: "keyDown", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9, modifiers: 8
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "keyUp", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9, modifiers: 8
  });
  await until('document.activeElement?.textContent.trim()', "Settings", "menu keyboard focus containment");
  assert.equal(await evaluate('document.querySelector(".mobile-scrim")?.tabIndex'), -1,
    "Backdrop must not be in the mobile keyboard tab order");
  assert.equal(await evaluate('document.querySelector(".mobile-scrim")?.getAttribute("aria-hidden")'), "true",
    "Backdrop must not compete with the drawer close button in the accessibility tree");

  // Even if focus is moved outside by browser automation or an extension,
  // the next Tab must return to the first control within the open drawer.
  await evaluate('document.querySelector(".hamburger")?.focus()');
  await cdp("Input.dispatchKeyEvent", {
    type: "keyDown", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "keyUp", key: "Tab", code: "Tab", windowsVirtualKeyCode: 9
  });
  await until('document.activeElement?.getAttribute("aria-label")', "Close menu",
    "focus trapped after external focus");
  await cdp("Input.dispatchKeyEvent", {
    type: "keyDown", key: "Escape", code: "Escape", windowsVirtualKeyCode: 27
  });
  await cdp("Input.dispatchKeyEvent", {
    type: "keyUp", key: "Escape", code: "Escape", windowsVirtualKeyCode: 27
  });
  await until('document.querySelector(".hamburger")?.getAttribute("aria-expanded")', "false", "Escape closes navigation");
  await until('document.activeElement?.getAttribute("aria-label")', "Open navigation", "Escape focus restoration");
  console.log("PASS keyboard skip link, mobile menu Enter/Shift+Tab/Escape and focus restoration");

  // Check the state transition in the opposite direction: an open drawer
  // must close when a window resize moves to the desktop breakpoint.
  await evaluate('document.querySelector(".hamburger")?.click()');
  await until('document.querySelector(".hamburger")?.getAttribute("aria-expanded")', "true",
    "drawer reopened before resize");
  await until('document.activeElement?.getAttribute("aria-label")', "Close menu",
    "drawer focus before desktop resize");
  await cdp("Emulation.setDeviceMetricsOverride", {
    width: 1280, height: 900, deviceScaleFactor: 1, mobile: false
  });
  await until('document.querySelector(".hamburger")?.getAttribute("aria-expanded")', "false",
    "desktop breakpoint resets mobile drawer");
  assert.equal(await evaluate('document.querySelector(".mobile-scrim") === null'), true,
    "Mobile backdrop remained mounted after desktop resize");
  await until('document.activeElement?.textContent.trim()', "Overview",
    "desktop sidebar navigation focus");
  await cdp("Emulation.setDeviceMetricsOverride", {
    width: 360, height: 900, deviceScaleFactor: 1, mobile: true
  });
  await until('document.activeElement?.getAttribute("aria-label")', "Open navigation",
    "mobile breakpoint restores focus to menu trigger");
  console.log("PASS bidirectional viewport transitions, drawer cleanup and focus recovery");

} finally {
  ws?.close();
  if (browser.exitCode === null) {
    const stopped = new Promise((resolve) => browser.once("exit", resolve));
    browser.kill("SIGTERM");
    await Promise.race([stopped, sleep(3000)]);
  }
  // Chrome may still flush cache after exit. Do not fail valid UI checks for
  // a temporary-profile cleanup race on an ephemeral CI runner.
  try {
    rmSync(work, { recursive: true, force: true, maxRetries: 10, retryDelay: 150 });
  } catch (error) {
    console.warn("Temporary browser profile cleanup deferred:", error.message);
  }
}
