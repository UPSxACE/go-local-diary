// Will save some persistent data so it's not lost on htmx-boost page changes (related to UI mostly)
window.$state = {};
// Force reload when using the go back page button (history API)
window.addEventListener("popstate", () => location.reload());
