import * as htmx from "htmx.org";
import * as feather from "feather-icons/dist/feather.min";
import { $load, $htmxLoad, $realLoad } from "../../lib/easy-dom";
window.htmx = htmx;


// Will save some persistent data so it's not lost on htmx-boost page changes (related to UI mostly)
window.$state = {};
// Force reload when using the go back page button (history API)
window.addEventListener("popstate", () => location.reload());

// feather.js
$load(() => {
  feather.replace();
  $htmxLoad(() => feather.replace());
});
